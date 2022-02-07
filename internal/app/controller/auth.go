package controller

import (
	"net/http"

	"github.com/YanxinTang/blog-server/internal/app/service"
	"github.com/YanxinTang/blog-server/internal/pkg/e"
	"github.com/YanxinTang/blog-server/internal/pkg/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type apiSignupModel struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

func Signup(c *gin.Context) {
	var apiSignup apiSignupModel
	if err := c.Bind(&apiSignup); err != nil {
		return
	}

	user := model.User{
		Username:    apiSignup.Username,
		Email:       apiSignup.Email,
		RawPassword: apiSignup.Password,
	}

	if err := service.InitUserAndSetting(user); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

type apiSigninModel struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
	Remember bool   `json:"remember" form:"remember"`
}

func Signin(c *gin.Context) {
	var apiSignin apiSigninModel
	if err := c.ShouldBind(&apiSignin); err != nil {
		c.Error(e.ERROR_EMPTY_ARGUMENTS)
		return
	}

	user, err := model.Authentication(apiSignin.Username, apiSignin.Password)
	if err != nil {
		c.Error(e.ERROR_INVALID_AUTH)
		return
	}

	session := sessions.Default(c)
	session.Set("login", true)
	session.Set("userID", user.ID)
	session.Set("user", user)
	if !apiSignin.Remember {
		session.Options(sessions.Options{Path: "/"})
	}
	session.Save()
	c.Status(http.StatusOK)
}

func Signout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Status(http.StatusOK)
}

func GetLoginSession(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("login") != true {
		c.Error(e.ERROR_SESSION_EXPIRED)
		return
	}
	c.JSON(http.StatusOK, session.Get("user"))
}
