package server

import (
	"time"

	"golang.org/x/crypto/bcrypt"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lon9/discord-generalized-sound-bot/backend/config"
	"github.com/lon9/discord-generalized-sound-bot/backend/controllers"
)

// NewRouter returns gin router
func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	cfg := config.GetConfig()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.GetStringSlice("server.cors"),
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	authMiddleware := &jwt.GinJWTMiddleware{
		Realm:      "Admin zone",
		Key:        []byte(config.GetConfig().GetString("auth.secret")),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,
		Authenticator: func(username, password string, c *gin.Context) (interface{}, bool) {
			isSuccess := func(password string) bool {
				if err := bcrypt.CompareHashAndPassword([]byte(cfg.GetString("auth.password")), []byte(password)); err != nil {
					return false
				}
				return true
			}
			if username == cfg.GetString("auth.username") && isSuccess(password) {
				return nil, true
			}
			return nil, false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"status":  code,
				"message": message,
			})
		},
		TokenLookup:   "header:Authorization",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	}

	health := new(controllers.HealthController)
	categories := new(controllers.CategoriesController)
	sounds := new(controllers.SoundsController)

	router.GET("/health", health.Status)

	categoriesRoutes := router.Group("/categories")
	{
		categoriesRoutes.GET("/", categories.Index)
		categoriesRoutes.GET("/:id", categories.Show)
	}

	soundRoutes := router.Group("/sounds")
	{
		soundRoutes.GET("/", sounds.Index)
	}

	router.POST("/login", authMiddleware.LoginHandler)
	admin := router.Group("/admin")
	admin.Use(authMiddleware.MiddlewareFunc())
	{
		admin.POST("/sounds", sounds.Create)
		admin.GET("refresh_token", authMiddleware.RefreshHandler)
	}

	return router
}
