package internal

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Router struct {
	service Service
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) RateMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		IPv4 := context.Request.Header.Get("X-Forwarded-For")
		if len(IPv4) == 0 {
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}

		status, err := r.service.HaveAccess(context, IPv4)
		if err != nil {
			context.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if !status {
			context.AbortWithStatus(http.StatusTooManyRequests)
			return
		}

		context.Next()
	}
}

func (r *Router) GetData(context *gin.Context) {
	context.JSON(http.StatusOK, nil)
}

func (r *Router) ClearRate(context *gin.Context) {
	context.JSON(http.StatusOK, nil)
}

func (r *Router) InitRouter() *gin.Engine {
	router := gin.New()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))

	api := router.Group("/api")
	{
		api.GET("/send", r.RateMiddleware(), r.GetData)
		api.POST("/clear", r.ClearRate)
	}

	return router
}
