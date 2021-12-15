package e

import "net/http"

var (
	ERROR_INTERVAL_ERROR     = New(http.StatusInternalServerError, "服务器内部错误")
	ERROR_RESOURCE_NOT_FOUND = New(http.StatusNotFound, "资源未找到")
	ERROR_EMPTY_ARGUMENTS    = New(http.StatusBadRequest, "用户名和密码不能为空")
	ERROR_INVALID_AUTH       = New(http.StatusBadRequest, "用户名或密码错误")
	ERROR_SESSION_EXPIRED    = New(http.StatusUnauthorized, "登录会话过期")
)
