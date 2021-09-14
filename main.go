package main

import (
	"bekasiberbagi/auth"
	"bekasiberbagi/campaign"
	"bekasiberbagi/handler"
	"bekasiberbagi/payment"
	"bekasiberbagi/response"
	"bekasiberbagi/transaction"
	"bekasiberbagi/user"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	webHandler "bekasiberbagi/web/handler"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	DB_HOST := os.Getenv("DB_HOST")
	DB_PORT := os.Getenv("DB_PORT")
	DB_NAME := os.Getenv("DB_NAME")
	DB_USER := os.Getenv("DB_USER")

	dsn := DB_USER + ":@tcp(" + DB_HOST + ":" + DB_PORT + ")/" + DB_NAME + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	userRepository := user.NewRepository(db)
	campaignRepository := campaign.NewRepository(db)
	transactionRepository := transaction.NewRepository(db)

	userService := user.NewService(userRepository)
	authService := auth.NewService()
	campaignService := campaign.NewService(campaignRepository)
	paymentService := payment.NewService()
	transactionService := transaction.NewService(transactionRepository, campaignRepository, paymentService)

	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	userWebHandler := webHandler.NewUserHandler(userService)
	campaignWebHandler := webHandler.NewCampaignHandler(campaignService, userService)
	transactionWebHandler := webHandler.NewTransactionHandler(transactionService)
	webAuthHandler := webHandler.NewWebAuthHandler(userService)

	router := gin.Default()
	router.Use(cors.Default())

	cookieStore := cookie.NewStore([]byte(auth.SECRET_KEY))
	router.Use(sessions.Sessions("bekasiberbagi", cookieStore))

	router.HTMLRender = loadTemplates("./web/templates")

	router.Use(static.Serve("/uploads", static.LocalFile("./uploads", true)))
	router.Use(static.Serve("/css", static.LocalFile("./web/assets/css", true)))
	router.Use(static.Serve("/js", static.LocalFile("./web/assets/js", true)))
	router.Use(static.Serve("/webfonts", static.LocalFile("./web/assets/webfonts", true)))

	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/check-email-availability", userHandler.IsEmailAvailability)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)
	api.GET("/users/fetch", authMiddleware(authService, userService), userHandler.FetchUser)

	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.GET("/campaigns/:id", campaignHandler.GetCampaign)
	api.POST("/campaigns", authMiddleware(authService, userService), campaignHandler.CreateCampaign)
	api.PUT("/campaigns/:id", authMiddleware(authService, userService), campaignHandler.UpdateCampaign)
	api.POST("/campaign-images", authMiddleware(authService, userService), campaignHandler.CreateCampaignImage)

	api.GET("/campaigns/:id/transactions", authMiddleware(authService, userService), transactionHandler.GetCampaignTransaction)
	api.GET("/transactions", authMiddleware(authService, userService), transactionHandler.GetUserTransactions)
	api.POST("/transactions", authMiddleware(authService, userService), transactionHandler.CreateTransaction)
	api.POST("/transactions/notification", transactionHandler.PaymentNotification)

	web := router.Group("/web")
	web.GET("/users", authAdminMiddleware(), userWebHandler.Index)
	web.GET("/users/create", authAdminMiddleware(), userWebHandler.Create)
	web.POST("/users", authAdminMiddleware(), userWebHandler.Store)
	web.GET("/users/:id/edit", authAdminMiddleware(), userWebHandler.Edit)
	web.POST("/users/:id", authAdminMiddleware(), userWebHandler.Update)
	web.GET("/users/:id/avatar", authAdminMiddleware(), userWebHandler.EditAvatar)
	web.POST("/users/:id/avatar", authAdminMiddleware(), userWebHandler.UpdateAvatar)

	web.GET("/campaigns", authAdminMiddleware(), campaignWebHandler.Index)
	web.GET("/campaigns/create", authAdminMiddleware(), campaignWebHandler.Create)
	web.POST("/campaigns", authAdminMiddleware(), campaignWebHandler.Store)
	web.GET("/campaigns/:id", authAdminMiddleware(), campaignWebHandler.Show)
	web.GET("/campaigns/:id/edit", authAdminMiddleware(), campaignWebHandler.Edit)
	web.POST("/campaigns/:id", authAdminMiddleware(), campaignWebHandler.Update)
	web.GET("/campaigns/:id/image", authAdminMiddleware(), campaignWebHandler.FormUploadImage)
	web.POST("/campaigns/:id/image", authAdminMiddleware(), campaignWebHandler.UploadImage)

	web.GET("/transactions", authAdminMiddleware(), transactionWebHandler.Index)

	web.GET("/login", webAuthHandler.LoginForm)
	web.POST("/login", webAuthHandler.LoginAction)
	web.GET("/logout", webAuthHandler.Logout)

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

func authAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		userIdSession := session.Get("userId")

		if userIdSession == nil {
			c.Redirect(http.StatusFound, "/web/login")
			return
		}
	}
}

func loadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	layouts, err := filepath.Glob(templatesDir + "/layouts/*")
	if err != nil {
		panic(err.Error())
	}

	includes, err := filepath.Glob(templatesDir + "/**/*")
	if err != nil {
		panic(err.Error())
	}

	// Generate our templates map from our layouts/ and includes/ directories
	for _, include := range includes {
		layoutCopy := make([]string, len(layouts))
		copy(layoutCopy, layouts)
		files := append(layoutCopy, include)
		r.AddFromFiles(filepath.Base(include), files...)
	}
	return r
}
