package web

import (
	"errors"
	"log"

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
// @Summary Obtiene la ubicación de la nave y el mensaje que emite.
// @Description Basado en las distancias y mensajes que se reciben de cada satelite, se obtienen la posicion y el mensaje emitido.
// @Param Body body model.TopSecretRequest true "The disantes and messages recieved from each satellite"
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
		c.IndentedJSON(http.StatusNotFound, model.ErrorResponse{Error: "Error binding JSON."})
		return
	}

	// treat request data to lists calculation form
	distances, messages, treatErr := TreatRequestData(requestData)
	if treatErr != nil {
		c.IndentedJSON(http.StatusNotFound, model.ErrorResponse{Error: treatErr.Error()})
		return
	}

	// calculates location
	x, y, locErr := location.CalculateLocation(distances)
	if locErr != nil {
		log.Printf("TopSecretHandler error with calculate location. Trace: %s", locErr.Error())
		c.IndentedJSON(http.StatusNotFound, model.ErrorResponse{Error: "Can't calculate location. Please check distances."})
		return
	}

	message, msgsErr := message.ConsolidateMessage(messages)
	if msgsErr != nil {
		log.Printf("TopSecretHandler error with consolidate message. Trace: %s", msgsErr.Error())
		c.IndentedJSON(http.StatusNotFound, model.ErrorResponse{Error: msgsErr.Error()})
		return
	}

	rspData := model.TopSecretResponse{
		Position: model.CoordinatesResponse{X: x, Y: y},
		Message:  message,
	}
	c.IndentedJSON(http.StatusOK, rspData)
}

func TreatRequestData(requestData model.TopSecretRequest) (distances []float32, messages [][]string, err error) {
	satellitesCount := store.GetSatellitesInfoCount()
	if len(requestData.Satellites) < satellitesCount {
		log.Printf("Insufficient request data. Satelites distances: %d, need at least %d", len(requestData.Satellites), satellitesCount)
		return distances, messages, errors.New("insufficient request data")
	}
	distances = make([]float32, satellitesCount)
	messages = make([][]string, satellitesCount)
	for _, rqSatelliteInfo := range requestData.Satellites {
		// gets index synchronized satellite info
		satIdx := store.GetSatelliteInfoIndex(rqSatelliteInfo.Name)

		// sets distance to the satellite via index idem like stored satellite info
		distances[satIdx] = rqSatelliteInfo.Distance

		// sets message to the satellite via index idem like stored satellite info
		messages[satIdx] = rqSatelliteInfo.Message
	}
	return distances, messages, nil
}

func TopSecretSplitHandler(c *gin.Context) {
	c.String(http.StatusNotImplemented, "")
}
