package controller

import (
	"net/http"

	"github.com/YanxinTang/blog-server/ent"
	"github.com/YanxinTang/blog-server/internal/app/common"
	"github.com/YanxinTang/blog-server/internal/app/service"
	"github.com/YanxinTang/blog-server/internal/pkg/e"
	"github.com/YanxinTang/blog-server/internal/pkg/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type SignupReqBody struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

func Signup(c *gin.Context) {
	var signupReqBody SignupReqBody
	if err := c.Bind(&signupReqBody); err != nil {
		return
	}

	user := ent.User{
		Username: signupReqBody.Username,
		Email:    signupReqBody.Email,
		Password: signupReqBody.Password,
	}

	if err := service.InitUserAndSetting(user); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

type SigninReqBody struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
	Remember bool   `json:"remember" form:"remember"`
}

func Signin(c *gin.Context) {
	var signinReqBody SigninReqBody
	if err := c.BindJSON(&signinReqBody); err != nil {
		return
	}

	user, err := model.Authentication(common.Context, common.Client)(signinReqBody.Username, signinReqBody.Password)
	if err != nil {
		c.Error(e.ERROR_INVALID_AUTH)
		return
	}

	session := sessions.Default(c)
	session.Set("login", true)
	session.Set("userID", user.ID)
	session.Set("user", user)
	if !signinReqBody.Remember {
		session.Options(sessions.Options{Path: "/"})
	}
	session.Save()
	c.JSON(http.StatusOK, user)
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
