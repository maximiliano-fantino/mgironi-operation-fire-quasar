package web

import (
	"github.com/gin-gonic/gin"
	"github.com/mgironi/operation-fire-quasar/support"

	"log"
	//	"operation-fire-quasar/docs"
	"github.com/mgironi/operation-fire-quasar/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitializeServer() {
	docs.SwaggerInfo.BasePath = "/"

	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.GET("/ping", PingHandler)
	router.POST("/topsecret/", TopSecretHandler)
	router.POST("/topsecret_split/:operation", TopSecretSplitPOSTHandler)
	router.GET("/topsecret_split/:operation", TopSecretSplitGETHandler)

	// swagger index
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Determine port for HTTP service.
	port := support.WebServerPort()

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	err := router.Run("0.0.0.0:" + port)
	if err != nil {
		log.Fatal(err)
	}

}
