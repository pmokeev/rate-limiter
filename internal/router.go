package internal

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Router struct {
	service Service
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) RateMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		log.Println("Middleware")
		context.Next()
	}
}

func (r *Router) GetData(context *gin.Context) {
	context.JSON(http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func (r *Router) ClearRate(context *gin.Context) {
	context.JSON(http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func (r *Router) InitRouter() *gin.Engine {
	router := gin.New()

	api := router.Group("/api")
	{
		api.GET("/send", r.RateMiddleware(), r.GetData)
		api.POST("/clear", r.ClearRate)
	}

	return router
}
