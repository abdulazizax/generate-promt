package api

import (
	"context"
	"generate-promt-v1/api/handlers"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"log"

	"generate-promt-v1/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"

	_ "generate-promt-v1/api/docs"

	"github.com/gin-contrib/cors"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func New(openaiClient *openai.Client, logger logger.Logger) *gin.Engine {
	router := gin.Default()
	router.Use(gin.Recovery())

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AllowHeaders = append(config.AllowHeaders, "*")

	router.Use(cors.New(config))
	sheetService, err := sheets.NewService(context.Background(), option.WithCredentialsFile("service_account.json"))
	if err != nil {
		log.Fatalf("Unable to create Sheets service: %v", err)
	}
	handler := handlers.New(openaiClient, logger, sheetService)

	router.POST("/prompt/execute", handler.ExecutePrompt)

	url := ginSwagger.URL("swagger/doc.json") // The url pointing to API definition
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	return router
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, HEAD, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
