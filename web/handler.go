package web

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mgironi/operation-fire-quasar/location"
	"github.com/mgironi/operation-fire-quasar/message"
	"github.com/mgironi/operation-fire-quasar/model"
	"github.com/mgironi/operation-fire-quasar/store"

	"net/http"
)

// @BasePath /

// Returns Echo message as ping response
// @Summary ping
// @Schemes
// @Description response ping
// @Tags example
// @Produce text/plain
// @Success 200 {string} echo
// @Router /ping/ [get]
func PingHandler(c *gin.Context) {
	c.String(http.StatusOK, "echo")
}

// @BasePath /
// @Summary Obtiene la ubicacion de la nave y el mensaje que emite.
// @Description Basado en las distancias y mensajes que se reciben de cada satelite, se obtienen la posicion y el mensaje emitido.
// @Param Body body model.TopSecretRequest true "Las distancias y mensajes recibidos por los satelites"
// @Accept json
// @Produce json
// @Failure 404 {object} model.ErrorResponse
// @Success 200 {object} model.TopSecretRequest
// @Router /topsecret/ [POST]
func TopSecretHandler(c *gin.Context) {
	var requestData model.TopSecretRequest

	// parse json to struct
	err := c.ShouldBindJSON(&requestData)
	if err != nil {
		log.Printf("Error binding json. Trace: %s", err.Error())
		c.IndentedJSON(http.StatusBadRequest, model.ErrorResponse{Message: "malformed json."})
		return
	}

	// performs calculations, checks and response data
	DoCalculationsAndResponse("TopSecretHandler", requestData.Satellites, c)
}

func DoCalculationsAndResponse(handlerName string, satellitesData []model.SatelliteInfoRequest, c *gin.Context) {
	// treat request data to lists calculation form
	distances, messages, treatErr := TreatSatellitesData(satellitesData)
	if treatErr != nil {
		c.IndentedJSON(http.StatusNotFound, model.ErrorResponse{Message: treatErr.Error()})
		return
	}

	// calculates location
	x, y, locErr := location.CalculateLocation(distances)
	if locErr != nil {
		log.Printf("%s error with calculate location. Trace: %s", handlerName, locErr.Error())
		c.IndentedJSON(http.StatusNotFound, model.ErrorResponse{Message: "Can't calculate location. Please check distances."})
		return
	}

	message, msgsErr := message.ConsolidateMessage(messages)
	if msgsErr != nil {
		log.Printf("%s error with consolidate message. Trace: %s", handlerName, msgsErr.Error())
		c.IndentedJSON(http.StatusNotFound, model.ErrorResponse{Message: msgsErr.Error()})
		return
	}

	rspData := model.TopSecretResponse{
		Position: model.CoordinatesResponse{X: x, Y: y},
		Message:  message,
	}
	c.IndentedJSON(http.StatusOK, rspData)
}

func TreatSatellitesData(satellitesData []model.SatelliteInfoRequest) (distances []float32, messages [][]string, err error) {
	satellitesCount := store.GetSatellitesInfoCount()
	if len(satellitesData) < satellitesCount {
		log.Printf("Insufficient request data. Satelites distances: %d, need at least %d", len(satellitesData), satellitesCount)
		return distances, messages, errors.New("insufficient request data")
	}
	distances = make([]float32, satellitesCount)
	messages = make([][]string, satellitesCount)
	for _, rqSatelliteInfo := range satellitesData {
		// gets index synchronized satellite info
		satIdx := store.GetSatelliteInfoIndex(rqSatelliteInfo.Name)

		// sets distance to the satellite via index idem like stored satellite info
		distances[satIdx] = rqSatelliteInfo.Distance

		// sets message to the satellite via index idem like stored satellite info
		messages[satIdx] = rqSatelliteInfo.Message
	}
	return distances, messages, nil
}

// @BasePath /
// @Summary Recibe la distancia de la nave y el mensaje que recibido por un satelite.
// @Description Recibe la distancia y mensaje que recibe un satelite y devuelve el token de operacion para posterior tratamiento.
// @Param operation path string false "El token de operacion"
// @Param Body body model.TopSecretSplitRequest true "La distancia y el mensaje recibido por un satelite"
// @Accept json
// @Produce json
// @Failure 404 {object} model.ErrorResponse
// @Success 200 {object} model.TopSecretSplitPOSTResponse
// @Router /topsecret_split/{operation} [POST]
func TopSecretSplitPOSTHandler(c *gin.Context) {
	// get operation token
	operation := strings.TrimSpace(c.Param("operation"))

	var requestData model.SatelliteInfoRequest

	// parse json to struct
	err := c.ShouldBindJSON(&requestData)
	if err != nil {
		log.Printf("Error binding json. Trace: %s", err.Error())

		// if data not ok, send response 404
		c.IndentedJSON(http.StatusBadRequest, model.ErrorResponse{Message: "malformed json."})
		return
	}

	isValid, errMsgStr := validateSatelliteInfoRequestData(requestData)
	if !isValid {
		log.Printf("Error data is incomplete. Data: %v. Trace: %s", requestData, errMsgStr)

		// if data not ok, send response 404
		c.IndentedJSON(http.StatusNotFound, model.ErrorResponse{Message: errMsgStr})
		return
	}

	savedDataset := store.GetDataset(operation, requestData.Message)
	if savedDataset.Key == "" {
		// get operation token
		operation = store.GetNewOperationUUID()

		// initialize dataset
		saved := store.SaveNewDataset(operation, requestData)
		if saved {
			response := model.TopSecretSplitPOSTResponse{Operation: operation}
			c.IndentedJSON(http.StatusOK, response)
		} else {
			log.Printf("Error in save new dataset. unsaved operacion: %s, request data:%v", operation, requestData)
			c.IndentedJSON(http.StatusInternalServerError, model.ErrorResponse{Message: "Can't save data."})
			return
		}
		return
	}

	countData := len(savedDataset.Satellites) + 1

	consolidatedMessage := ""

	// checks if dataset is completed with new data
	if countData < store.GetSatellitesInfoCount() {
		// dataset is incomplete, consolidate partial message with data that have
		messages := [][]string{}
		for _, satData := range savedDataset.Satellites {
			messages = append(messages, satData.Message)
		}
		messages = append(messages, requestData.Message)

		// gets consolidated message
		var consErr error
		consolidatedMessage, consErr = message.ConsolidateMessage(messages)
		if consErr != nil {
			c.IndentedJSON(http.StatusNotFound, model.ErrorResponse{Message: "Can't consolidate message."})
			return
		}
	}

	// update dataset
	updated := store.UpdateDataset(operation, consolidatedMessage, savedDataset.Key, requestData)
	if !updated {
		log.Printf("Error in update dataset. operacion: %s, message: %s, previous key: %s, request data:%v", operation, consolidatedMessage, savedDataset.Key, requestData)
		c.IndentedJSON(http.StatusInternalServerError, model.ErrorResponse{Message: "Can't update data."})
		return
	}
	if operation == "" {
		operation = savedDataset.Operation
	}
	response := model.TopSecretSplitPOSTResponse{Operation: operation}
	c.IndentedJSON(http.StatusOK, response)
}

func validateSatelliteInfoRequestData(requestData model.SatelliteInfoRequest) (isValid bool, validationErrors string) {
	isValid = false
	validationErrors = ""
	if store.GetSatelliteInfoIndex(requestData.Name) == -1 {
		validationErrors += fmt.Sprintf("Not exists satellite reference data for '%s'", requestData.Name)
	}
	isValid = true
	return
}

// @BasePath /
// @Summary Recibe el token correspondiente a la operacion de colleccion de datos de distancias y mensajes enviados a /TopSecretSplit/ previamente.
// @Description Recibe el token de operacion y con el set de datos previamente recolectado basado en las distancias y mensajes que se reciben de cada satelite, se obtienen la posicion y el mensaje emitido.
// @Param operation path string true "El token de operacion"
// @Produce json
// @Failure 404 {object} model.ErrorResponse
// @Success 200 {object} model.TopSecretResponse
// @Router /topsecret_split/{operation} [GET]
func TopSecretSplitGETHandler(c *gin.Context) {
	// get operation token
	operation := c.Param("operation")

	if operation == "" {
		resErr := model.ErrorResponse{Message: "operation token required"}
		c.IndentedJSON(http.StatusNotFound, resErr)
	}

	// get dataset directly by operation
	dataset := store.GetDatasetByKey(operation)
	if dataset.Key == "" || dataset.Key != operation {
		resErr := model.ErrorResponse{Message: "Insufficient information"}
		c.IndentedJSON(http.StatusNotFound, resErr)
		return
	}

	// performs calculations, checks and response data
	DoCalculationsAndResponse("TopSecretSplitGETHandler", dataset.Satellites, c)
}
