package store

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/mgironi/operation-fire-quasar/model"
)

// Satelites global variable
var satelites map[int]model.SateliteInfo

// Initialices in memory store
func Initialize() {

	// the satelites info
	InitializeSatelitesInfo()
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
