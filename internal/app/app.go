package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"avito/internal/config"
	delivery "avito/internal/delivery/http"
	"avito/internal/gen"
	"avito/internal/log"
	"avito/internal/postgres"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Start() {
	cfg := config.InitConfig()

	log.Log.Info("Config Initialized")

	db := postgres.MustInitPg(cfg)
	defer db.Close()

	log.Log.Info("PG Initialized")

	g := gin.New()

	handlers := delivery.InitServer(db)

	g.Use(gin.Logger(), gin.Recovery())

	g.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           time.Hour,
	}))

	g.GET("/openapi.json", func(c *gin.Context) {
		swagger, _ := gen.GetSwagger()

		jsonData, err := json.Marshal(swagger)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to marshal OpenAPI spec",
			})

			return
		}

		c.Data(http.StatusOK, "application/json", jsonData)
	})

	g.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.URL("/openapi.json"),
	))

	gen.RegisterHandlers(g, handlers)

	log.Log.Info("Start Server")

	err := g.Run(fmt.Sprintf("%v:%v", cfg.ServiceHost, cfg.ServicePort))
	if err != nil {
		panic(fmt.Sprintf("error running client: %v", err.Error()))
	}
}
