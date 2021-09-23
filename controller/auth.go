package controller

import (
	"log"
	"net/http"

	"github.com/YanxinTang/blog/server/e"
	"github.com/YanxinTang/blog/server/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var (
	ERROR_EMPTY_ARGUMENTS = e.New(http.StatusBadRequest, "用户名和密码不能为空")
	ERROR_INVALID_AUTH    = e.New(http.StatusBadRequest, "用户名和或密码错误")
	ERROR_INTERVAL_ERROR  = e.New(http.StatusInternalServerError, "服务器内部错误")
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
	if err := model.CreateUser(apiSignup.Username, apiSignup.Email, apiSignup.Password); err != nil {
		c.Error(err)
		return
	}
	model.SetSetting("signupEnable", false)
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
		log.Println(err)
		c.Error(ERROR_EMPTY_ARGUMENTS)
		return
	}

	user, err := model.Authentication(apiSignin.Username, apiSignin.Password)
	if err != nil {
		c.Error(ERROR_INVALID_AUTH)
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
