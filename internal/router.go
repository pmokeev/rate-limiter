package internal

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Router struct {
	service *Service
}

func NewRouter(service *Service) *Router {
	return &Router{service: service}
}

func (r *Router) RateMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		IPv4 := context.Request.Header.Get("X-Forwarded-For")
		if len(IPv4) == 0 {
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}

		access, err := r.service.CheckAccess(context, IPv4)
		if err != nil {
			context.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if !access.isAccess {
			context.Writer.Header().Set("X-Ratelimit-Retry-After", access.xRetryAfter.String())
			context.AbortWithStatus(http.StatusTooManyRequests)
			return
		}

		context.Writer.Header().Set("X-Ratelimit-Remaining", strconv.Itoa(int(access.xRemaining)))
		context.Writer.Header().Set("X-Ratelimit-Limit", strconv.Itoa(int(access.xLimit)))

		context.Next()
	}
}

func (r *Router) GetData(context *gin.Context) {
	data, err := r.service.GetData()
	if err != nil {
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	context.JSON(http.StatusOK, map[string]string{
		"message": data,
	})
}

func (r *Router) ClearRate(context *gin.Context) {
	IPv4 := context.Request.Header.Get("X-Forwarded-For")
	err := r.service.ClearRate(context, IPv4)
	if err != nil {
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	context.JSON(http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func (r *Router) InitRouter() *gin.Engine {
	router := gin.New()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := router.Group("/api")
	{
		api.GET("/send", r.RateMiddleware(), r.GetData)
		api.POST("/clear", r.ClearRate)
	}

	return router
}
