package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/panithee/internship_day2/dto"
	"github.com/panithee/internship_day2/models"
	"github.com/panithee/internship_day2/service"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// login contorller interface
type LoginController interface {
	Login(ctx *gin.Context, db *gorm.DB)
}

type loginController struct {
	// loginService service.LoginService
	jWtService service.JWTService
}

func LoginHandler(
	// loginService service.LoginService,
	jWtService service.JWTService) LoginController {
	return &loginController{
		// loginService: loginService,
		jWtService: jWtService,
	}
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (controller *loginController) Login(ctx *gin.Context, db *gorm.DB) {
	var credential dto.LoginCredentials
	err := ctx.BindJSON(&credential)
	if err != nil {
		ctx.IndentedJSON(500, gin.H{"message": "no data found"})
	}
	log.Print(credential)

	var user models.Users
	// result := db.Where("email = ? AND password = ?", credential.Email, credential.Password).First(&user)

	result := db.Where("email = ?", credential.Email).First(&user)

	if result.Error != nil || !CheckPasswordHash(credential.Password, user.Password) {
		ctx.IndentedJSON(401, gin.H{"message": "email or password is not match"})
		return

	} else {
		token := controller.jWtService.GenerateToken(credential.Email, true)

		if token != "" {
			ctx.IndentedJSON(http.StatusOK, gin.H{
				"data":    token,
				"message": "login success",
			})
		} else {
			ctx.IndentedJSON(http.StatusUnauthorized, nil)
		}
	}

	// isUserAuthenticated := controller.loginService.LoginUser(credential.Email, credential.Password)
	// if isUserAuthenticated {
	// return controller.jWtService.GenerateToken(credential.Email, true)

	// }
	// return ""
}
