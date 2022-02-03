package store

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/mgironi/operation-fire-quasar/model"
	"github.com/mgironi/operation-fire-quasar/support"

	"github.com/google/uuid"

	"github.com/gomodule/redigo/redis"

	fuzzy "github.com/paul-mannino/go-fuzzywuzzy"
)

// satelites variable
var satelites map[int]model.SateliteInfo

// redis connection pool
var redisPool *redis.Pool

var GetRedisConnection = func() redis.Conn {
	if redisPool == nil {
		log.Print("Redis pool connection not initialized")
		return nil
	}
	return redisPool.Get()
}

func InitializeCmd() {

	// the satelites info
	InitializeSatelitesInfo()
}

// Initialices in memory store
func Initialize() {

	// the satelites info
	InitializeSatelitesInfo()

	// loads memory cache connection (Redis)
	InitializeMemorycacheConnection()
}

// HELP message for passing stalites info by environment variables
const HELP_PASSING_SATELITES_INFO_ENV = "To load satelites information plese use format '<name>_<xcoord>,<ycoord>'. Example: kenobi_300.25,-340.78"

// Env key for Kenobi satelite info
const SATELITE_KENOBI_ENV string = "OFQ_KENOBI"

// Env Key for Skywalker satelite info
const SATELITE_SKYWALKER_ENV string = "OFQ_SKYWALKER"

// Env key for Sato satelite info
const SATELITE_SATO_ENV string = "OFQ_SATO"

// Initialices satelites info
func InitializeSatelitesInfo() {
	// defines environment key to get satelite info
	satelitesEnvsKeys := []string{SATELITE_KENOBI_ENV, SATELITE_SKYWALKER_ENV, SATELITE_SATO_ENV}

	// initialices satelites map
	var hasErrors bool
	satelites, hasErrors = ParseSatelitesInfoFromEnvs(satelitesEnvsKeys)

	// checks for parsing errors
	if hasErrors {
		// loads default
		LoadsDefaultSatelitesInfo()
	}
}

func InitializeMemorycacheConnection() {
	redisHost := support.RedisHost()
	redisPort := support.RedisPort()
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)
	redisConnProtocol := "tcp"
	const maxConnections = 10
	log.Printf("connecting to redis memorystore on %s...", redisAddr)
	redisPool = &redis.Pool{
		MaxIdle: maxConnections,
		Dial: func() (redis.Conn, error) {
			conn, cnErr := redis.Dial(redisConnProtocol, redisAddr)
			if cnErr != nil {
				log.Printf("Error dialing to redis. Protocol: '%s', address: '%s'", redisConnProtocol, redisAddr)
			}
			return conn, cnErr
		},
	}
}

// Loads default satelite info
func LoadsDefaultSatelitesInfo() {
	log.Printf("\nContinue loading default Satelites information ...")

	// loads default satelites info
	satelites = map[int]model.SateliteInfo{
		0: {Name: "kenobi", Location: model.Point{X: -500, Y: -200}},
		1: {Name: "skywalker", Location: model.Point{X: 100, Y: -100}},
		2: {Name: "sato", Location: model.Point{X: 500, Y: 100}},
	}
}

func GetSatellitesInfo() (satelliteList []model.SateliteInfo) {
	satelliteList = make([]model.SateliteInfo, len(satelites))
	satelliteCount := len(satelites)
	for i := 0; i < satelliteCount; i++ {
		satelliteList[i] = satelites[i]
	}
	return satelliteList
}

// Routput: the kwnown reference coordinates.
func GetKnownReferenceCoordinates() (points []model.Point) {
	checksAndInitialicesSatellitesInfo()
	points = make([]model.Point, len(satelites))
	for i, satelite := range satelites {
		points[i] = satelite.Location
	}
	return points
}

func GetSatellitesInfoCount() int {
	checksAndInitialicesSatellitesInfo()
	return len(satelites)
}

func checksAndInitialicesSatellitesInfo() {
	if len(satelites) == 0 {
		InitializeSatelitesInfo()
	}
}

// input: satellite ;name'
// output: the satellite info index in store. Returns -1 if 'name' not present
func GetSatelliteInfoIndex(name string) (index int) {
	checksAndInitialicesSatellitesInfo()

	for i, satInfo := range satelites {
		if satInfo.Name == name {
			index = i
			return i
		}
	}
	return -1
}

// Parses satelites info from environment variables
// input: environment variables keys to parse
func ParseSatelitesInfoFromEnvs(envKeys []string) (satelitesInfo map[int]model.SateliteInfo, hasErrors bool) {
	satelitesInfo = make(map[int]model.SateliteInfo, len(envKeys))

	// intitialices convertion errore list
	convertionErrors := make([]error, len(envKeys))

	hasErrors = false

	// applies convertion for each satelite info in system envirnoment
	for i, envKey := range envKeys {
		envValue, envPresent := os.LookupEnv(envKey)
		if envPresent {
			satelitesInfo[i], convertionErrors[i] = ConvertSateliteInfo(envValue)
		} else {
			log.Printf("WARN: env variable %s not present.", envKeys[i])
			hasErrors = true
		}
	}

	for i, convError := range convertionErrors {
		if convError != nil {
			hasErrors = true
			log.Printf("WARN: Can't parse '%s' env variable value, %s", envKeys[i], convError)
		}
	}
	return satelitesInfo, hasErrors
}

// Converts satelite info string to SateliteInfo
// input: satelite info string (format '<name>_<xcoord>,<ycoord>')
// output: the SateliteInfo. See also SateliteInfo struct
func ConvertSateliteInfo(infoStr string) (info model.SateliteInfo, err error) {
	const GENERIC_ERROR_MSG string = "Insuficient data. To load satelites information plese use format '<name>_<xcoord>,<ycoord>'. Example: kenobi_300.25,-340.78"

	infoList := strings.Split(infoStr, "_")
	if len(infoList) < 2 {
		return info, fmt.Errorf("no posible get satelite info '%s'. %s", infoStr, GENERIC_ERROR_MSG)
	}
	locationCoords := strings.Split(infoList[1], ",")
	if len(locationCoords) < 2 {
		return info, fmt.Errorf("no posible get satelite info location coordinates '%s'. %s", infoStr, GENERIC_ERROR_MSG)
	}

	parsedCoords := make([]float64, len(locationCoords))
	for i, coordStr := range locationCoords {
		var parseError error
		parsedCoords[i], parseError = strconv.ParseFloat(coordStr, 64)
		if parseError != nil {
			return info, fmt.Errorf("can't get satelite info location coordinate '%s'. %s. %s", infoStr, parseError.Error(), GENERIC_ERROR_MSG)
		}
	}
	info = model.SateliteInfo{Name: infoList[0], Location: model.Point{X: parsedCoords[0], Y: parsedCoords[1]}}
	return info, nil
}

var GetNewOperationUUID = func() (id string) {
	id = uuid.New().String()
	return
}

const DATASET_KEY_FORMAT_PATTERN = "%s:%s"

const REDIS_MATCH_PATTERN_WILDCARD = "*"

const MESSAGE_KEY_SEPARATOR = " "

func SaveNewDataset(operation string, dataValue model.SatelliteInfoRequest) (saved bool) {
	// build key with <operation>:<string_message>
	stringMsg := strings.Join(dataValue.Message, MESSAGE_KEY_SEPARATOR)
	dataSetKey := fmt.Sprintf(DATASET_KEY_FORMAT_PATTERN, operation, stringMsg)
	dataset := model.Dataset{
		Key:        dataSetKey,
		Operation:  operation,
		Satellites: []model.SatelliteInfoRequest{dataValue},
	}
	saved = SetKeyValuePair(dataSetKey, dataset)
	return
}

func UpdateDataset(operation string, consolidatedMessage string, previousKey string, dataValue model.SatelliteInfoRequest) (saved bool) {
	saved = false
	dataset := GetDatasetByKey(previousKey)

	// checks and updates operation in dataset
	if operation != "" {
		// set operation value if not defined
		if dataset.Operation == "" {
			dataset.Operation = operation
		} else if dataset.Operation != operation {
			// if the operation is not the same, then we have a problem
			log.Printf("WARN operation in dataset mismatch. Given: '%s'. Existent dataset: '%v'", operation, dataset)
		}
	} else {
		operation = dataset.Operation
	}

	// adds new satellite data
	dataset.Satellites = append(dataset.Satellites, dataValue)

	var newDataSetKey string
	if consolidatedMessage != "" {
		// build key with <operation>:<string_message>
		//stringMsg := strings.Join(dataValue.Message, MESSAGE_KEY_SEPARATOR)
		//newDataSetKey = fmt.Sprintf(DATASET_KEY_FORMAT_PATTERN, operation, stringMsg)
		newDataSetKey = fmt.Sprintf(DATASET_KEY_FORMAT_PATTERN, operation, consolidatedMessage)
	} else if operation != "" {
		// use key just with operation value
		newDataSetKey = operation
	} else {
		// no message and no operation, impossible to build key
		return
	}

	// updates with the new key
	dataset.Key = newDataSetKey

	// keep safe old dataset
	oldDataset := GetByKey(previousKey)

	// remove old key,value
	DeleteKey(previousKey)

	// save new key, value
	saved = SetKeyValuePair(newDataSetKey, dataset)
	if !saved {
		// Restores previous key
		SetKeyValuePair(newDataSetKey, oldDataset)
	}
	return
}

func DeleteKey(key string) (success bool) {
	success = false

	cnn := GetRedisConnection()
	if cnn == nil {
		return
	}

	// apply 'DEL'
	delCount, setErr := redis.Int(cnn.Do("DEL", key))
	if setErr != nil {
		log.Printf("Error in DEL to redis. Key: %s. Trace: %s", key, setErr.Error())
		return
	}
	/*if reply == nil {
		log.Printf("Error in DEL to redis. Key: %s. Reply is nil", key)
		return
	}*/

	// gets deleted count
	//delCount := redis.Int(reply)
	success = delCount >= 1
	if delCount > 1 {
		// more than one is a problem
		log.Printf("WARN in DEL to redis. Key: %s. Deleted more than 1 key. total deleted: %d", key, delCount)
	}
	return
}

func SetKeyValuePair(key string, value interface{}) (success bool) {
	success = false
	cnn := GetRedisConnection()
	if cnn == nil {
		return
	}

	serialized, srlErr := json.Marshal(value)
	if srlErr != nil {
		log.Printf("Error serializing data. Key: %s, value: %v. Trace: %s", key, value, srlErr.Error())
		return
	}
	// apply 'SET' only if not exists key, using 'NX' arg
	reply, setErr := cnn.Do("SET", key, serialized, "NX")
	if reply == nil || setErr != nil {
		log.Printf("Error in SET to redis. Key: %s, value: %v. Trace: %s", key, value, setErr.Error())
		return
	}
	success = (reply.(string) == "OK")
	return
}

func GetDataset(operation string, message []string) (dataset model.Dataset) {

	// if operataion is not empty then get key by operation
	if operation != "" {
		dataset = GetDatasetByOperation(operation)
	} else if len(message) > 0 {
		// if operation is empty should search by matching using partial message (assuming that exits a more complete message)
		dataset = GetDatasetByMessage(message)
	}

	return
}

func buildMessageMatchPattern(message []string) (flattedMessage string) {
	msgToWild := append([]string(nil), message...)
	for i := 0; i < len(msgToWild); i++ {
		if strings.TrimSpace(msgToWild[i]) == "" {
			msgToWild[i] = REDIS_MATCH_PATTERN_WILDCARD
		}
	}
	flattedMessage = strings.Join(msgToWild, " ")
	return
}

func GetDatasetByMessage(message []string) (dataset model.Dataset) {
	filledWildcardMsg := buildMessageMatchPattern(message)
	matchFilterOperation := fmt.Sprintf(DATASET_KEY_FORMAT_PATTERN, REDIS_MATCH_PATTERN_WILDCARD, filledWildcardMsg)
	keys := ScanKeys(matchFilterOperation)
	countKeys := len(keys)
	if countKeys == 0 {
		// search by full scan with fuzzywuzzy
		log.Printf("WARN Key not found, trying filter by fuzzy match process, message: %s", message)
		matchKey := ScanWithFuzzy(message)
		if matchKey != "" {
			dataset = GetDatasetByKey(matchKey)
		}
		return dataset
	}
	if countKeys > 1 {
		log.Printf("WARN found more than one key, scanning by message: %s. Using first key:%s", filledWildcardMsg, keys[0])
	}
	dataset = GetDatasetByKey(keys[0])
	return dataset
}

func GetDatasetByOperation(operation string) (dataset model.Dataset) {
	matchFilterOperation := operation + ":*"
	keys := ScanKeys(matchFilterOperation)
	countKeys := len(keys)
	if countKeys == 0 {
		log.Printf("Key not found, after scanning by operation: %s", operation)
		return dataset
	}
	if countKeys > 1 {
		log.Printf("WARN found more than one key, scanning by operation: %s. Using first key:%s", operation, keys[0])

	}
	dataset = GetDatasetByKey(keys[0])
	return dataset
}

func GetDatasetByKey(key string) (dataset model.Dataset) {
	serializedDataset := GetByKey(key)
	if len(serializedDataset) == 0 {
		return
	}
	umErr := json.Unmarshal([]byte(serializedDataset), &dataset)
	if umErr != nil {
		log.Printf("Error trying to unmarshal value '%s' to '[]model.SatelliteInfoRequest'", serializedDataset)
	}
	return dataset
}

func partialScan(matchFilter string, scanCursor *string, results *[]string) {
	cnn := GetRedisConnection()
	if cnn == nil {
		return
	}
	// run scan with match filter
	currentCursor := *scanCursor
	reply, scanErr := redis.Values(cnn.Do("SCAN", currentCursor, "MATCH", matchFilter))
	if len(reply) != 2 || scanErr != nil {
		log.Printf("Error in SCAN to redis. match filter: %s, current cursor: %s, reply: %v. Trace: %s", matchFilter, *scanCursor, reply, scanErr.Error())
		return
	}

	/*if _, ok := reply[0].([]uint8); !ok {
		log.Printf("Error in SCAN reply of redis, reply[0] is not uint8. type: %s.", reflect.TypeOf(reply[0]).String())
		return
	}*/

	/*if _, ok := reply[1].([]interface{}); !ok {
		log.Printf("Error in SCAN reply of redis, reply[1] is not []interface{}. type: %s.", reflect.TypeOf(reply[1]).String())
		return
	}*/

	// updates scan cursor
	var errRInt error
	*scanCursor, errRInt = redis.String(reply[0], nil)
	if errRInt != nil {
		log.Printf("Error in SCAN reply of redis, reply[0] parse error. Trace: %s.", errRInt.Error())
		return
	}

	// updates results
	var strErr error
	*results, strErr = redis.Strings(reply[1], nil)
	if strErr != nil {
		log.Printf("Error in SCAN reply of redis, reply[1] parse error. Trace: %s.", strErr.Error())
		return
	}

}

const FUZZY_PROCESS_MIN_SCORING_ACCEPTED int = 60

func ScanWithFuzzy(message []string) (matchKey string) {
	scanCursor := "0"
	// scan for all keys
	matchFilter := "*"
	extractionPhrase := strings.Join(message, " ")
	minScoring := FUZZY_PROCESS_MIN_SCORING_ACCEPTED
	keys := []string{}
	matchChoices := []matchChoiceKey{}

	// do partial scan
	partialScan(matchFilter, &scanCursor, &keys)

	// apply filter by fuzzy match scoring
	resultFilter := filterKeysByFuzzyProcessScore(extractionPhrase, minScoring, keys)

	// collect result filter
	matchChoices = append(matchChoices, resultFilter...)

	for scanCursor != "0" {
		// do partial scan
		partialScan(matchFilter, &scanCursor, &keys)

		// apply filter by fuzzy match scoring
		resultFilter = filterKeysByFuzzyProcessScore(extractionPhrase, minScoring, keys)

		// collect result filter
		matchChoices = append(matchChoices, resultFilter...)
	}

	// select result with the best scoring
	sort.Slice(matchChoices, func(i, j int) bool {
		return matchChoices[i].Score > matchChoices[j].Score
	})

	if len(matchChoices) > 0 {
		// selects the first elem of sorted results
		matchKey = matchChoices[0].key
	}
	return

}

type matchChoiceKey struct {
	Match string
	Score int
	key   string
}

func filterKeysByFuzzyProcessScore(extractionPhrase string, processMinScoring int, keys []string) (matchs []matchChoiceKey) {
	// clean key from operation segment
	choices := make([]string, len(keys))
	linkedKeys := make(map[string]string, len(keys))
	for _, key := range keys {
		sepIdx := strings.IndexAny(key, ":")
		// select as choices those keys with '_:_' pattern
		if sepIdx != -1 {
			// clean prefix of the key (operation) to get a clean choice
			choice := key[sepIdx+1:]
			choices = append(choices, choice)
			// maintain keys linked (the full key) to recover it after scoring filter
			linkedKeys[choice] = key
		}
	}

	// apply extract of the 3 choices with more match scoring
	matchPairs, matchErr := fuzzy.Extract(extractionPhrase, choices, 3)
	if matchErr != nil {
		log.Printf("Error during fuzzy extract process. extraction phrase: '%s' choices list size: '%d'. Trace: %s", extractionPhrase, len(choices), matchErr.Error())
		return
	}

	// build result match list with recovering full key
	matchs = make([]matchChoiceKey, len(matchPairs))
	for _, matchPair := range matchPairs {
		// only accept those who have a higher scoring
		if matchPair.Score >= processMinScoring {
			mtcChoiceKey := matchChoiceKey{
				key:   linkedKeys[matchPair.Match],
				Match: matchPair.Match,
				Score: matchPair.Score,
			}
			matchs = append(matchs, mtcChoiceKey)
		}
	}
	return
}

func ScanKeys(matchFilter string) (keys []string) {
	scanCursor := "0"
	partialScan(matchFilter, &scanCursor, &keys)
	for scanCursor != "0" {
		partialScan(matchFilter, &scanCursor, &keys)
	}
	return
}

func GetByKey(key string) (value string) {
	cnn := GetRedisConnection()
	if cnn == nil {
		return
	}
	var getErr error
	value, getErr = redis.String(cnn.Do("GET", key))
	if getErr != nil {
		log.Printf("Error in GET to redis. Key: %s. Trace: %s", key, getErr.Error())
		return
	}
	return
}
