package main

import (
	"log"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/panithee/internship_day2/controller"
	"github.com/panithee/internship_day2/dto"
	"github.com/panithee/internship_day2/middleware"
	"github.com/panithee/internship_day2/models"
	"github.com/panithee/internship_day2/service"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type RegisterFormation struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginFomation struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Message struct {
	Message string `json:"message"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func register(db *gorm.DB, c *gin.Context) {
	var registerForm RegisterFormation

	if err := c.BindJSON(&registerForm); err != nil {
		log.Println("Cannot bind received JSON to registerForm")
	}

	hashPassword, err := HashPassword(registerForm.Password)
	if err != nil {
		log.Println(err)
	}

	log.Println(hashPassword)

	// match := CheckPasswordHash(registerForm.Password, hashPassword)
	// log.Println(match)

	user := models.Users{
		Email: registerForm.Email, Password: hashPassword,
	}

	data := db.Create(&user)
	if data.Error != nil {
		c.IndentedJSON(400, gin.H{"message": data.Error.Error()})
	} else {
		c.IndentedJSON(200, gin.H{"message": "register success"})
	}

}

func Post(db *gorm.DB, c *gin.Context) {
	var message dto.Message

	if err := c.BindJSON(&message); err != nil {
		log.Println(err)
	}

	claims, exists := c.Get("claims")
	if !exists {
		log.Println("Claims not found in context")
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	// Type assertion and validation for 'name' claim
	nameClaim, ok := claims.(jwt.MapClaims)["name"]
	if !ok {
		log.Println("Missing 'name' claim")
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	userEmail, emailOk := nameClaim.(string)
	if !emailOk {
		log.Println("Invalid 'name' claim type")
		c.IndentedJSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	log.Println(userEmail)

	var user models.Users

	// get user data
	result := db.Where("email = ?", userEmail).First(&user)
	if result.Error != nil {
		c.IndentedJSON(404, gin.H{"message": result.Error.Error()})
	}

	// assign user.Id into UserId for references
	postStruct := models.Posts{
		Message: message.Message,
		UserId:  user.ID,
	}
	// save data in posts table in database
	db.Create(&postStruct)
	c.IndentedJSON(200, gin.H{"data": message.Message, "message": "post success"})

}

func GetAllUserPosts(db *gorm.DB, c *gin.Context) {

	var postsResponse []struct {
		Email     string
		Message   string
		CreatedAt string
	}
	result := db.Table("posts").
		Select("users.email, posts.message, posts.created_at").
		Joins("JOIN users ON posts.user_id = users.id").
		Order("posts.created_at DESC").
		Find(&postsResponse)

	if result.Error != nil {
		c.IndentedJSON(500, gin.H{"message": result.Error.Error()})
	}
	c.IndentedJSON(200, postsResponse)

}

func main() {

	// var loginService service.LoginService = service.DynamicLoginService()
	var jwtService service.JWTService = service.JWTAuthService()
	var loginController controller.LoginController = controller.LoginHandler(
		// loginService,
		jwtService)

	// connect db with insert user as user and pass as password and db is database name in mysql
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := "user:pass@tcp(127.0.0.1:3306)/db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Panic(err)
	}

	// Migrate the schema
	db.AutoMigrate(&models.Users{}, &models.Posts{})

	router := gin.Default()
	router.POST("/register", func(ctx *gin.Context) {
		register(db, ctx)
	})

	router.POST("/login", func(ctx *gin.Context) {
		loginController.Login(ctx, db)

	})

	// v1 := router.Group("/v1")
	// v1.Use(middleware.AuthorizeJWT())
	// {
	// 	v1.POST("/post", func(ctx *gin.Context) {
	// 		Post(db, ctx)
	// 	})
	// }

	router.GET("/post", func(ctx *gin.Context) {
		GetAllUserPosts(db, ctx)
	})

	router.POST("/post", middleware.AuthorizeJWT(), func(ctx *gin.Context) {
		Post(db, ctx)
	})

	router.Run(":8080")

}
