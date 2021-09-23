package e

import "net/http"

var (
	ERROR_INVALID_AUTH       = New(http.StatusBadRequest, "用户名或密码错误")
	ERROR_SESSION_EXPIRED    = New(http.StatusUnauthorized, "登录会话过期")
	ERROR_RESOURCE_NOT_FOUND = New(http.StatusNotFound, "资源未找到")
)
