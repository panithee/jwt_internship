package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/panithee/internship_day2/dto"
	"github.com/panithee/internship_day2/service"
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

func (controller *loginController) Login(ctx *gin.Context, db *gorm.DB) {
	var credential dto.LoginCredentials
	err := ctx.BindJSON(&credential)
	if err != nil {
		ctx.IndentedJSON(500, gin.H{"message": "no data found"})
	}
	log.Print(credential)

	var user dto.Users
	result := db.Where("email = ? AND password = ?", credential.Email, credential.Password).First(&user)

	if result.Error != nil {
		ctx.IndentedJSON(401, gin.H{"message": "username or password is not matching"})
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
