package main

import (
	"bekasiberbagi/auth"
	"bekasiberbagi/campaign"
	"bekasiberbagi/handler"
	"bekasiberbagi/response"
	"bekasiberbagi/user"
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/bekasiberbagi?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	userRepository := user.NewRepository(db)
	campaignRepository := campaign.NewRepository(db)

	userService := user.NewService(userRepository)
	authService := auth.NewService()
	campaignService := campaign.NewService(campaignRepository)

	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)

	router := gin.Default()

	router.Use(static.Serve("/uploads", static.LocalFile("./uploads", true)))

	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/check-email-availability", userHandler.IsEmailAvailability)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)

	api.GET("/campaigns", campaignHandler.GetCampaigns)

	router.Run()
}

func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := response.APIResponseFailed("Unautorized", http.StatusUnauthorized)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		tokenString := ""

		arrayToken := strings.Split(authHeader, " ")

		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		token, err := authService.ValidateToken(tokenString)

		if err != nil {
			response := response.APIResponseFailed("Unautorized", http.StatusUnauthorized)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			response := response.APIResponseFailed("Unautorized", http.StatusUnauthorized)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userId := int(claim["user_id"].(float64))

		user, err := userService.GetUserById(userId)

		if err != nil {
			response := response.APIResponseFailed("Unautorized", http.StatusUnauthorized)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", user)
	}
}
