package e

import "net/http"

var (
	ERROR_INTERVAL_ERROR     = New(http.StatusInternalServerError, "服务器内部错误")
	ERROR_BAD_REQUEST        = New(http.StatusBadRequest, "异常请求")
	ERROR_RESOURCE_NOT_FOUND = New(http.StatusNotFound, "资源未找到")
	ERROR_BEGIN_TX           = New(http.StatusInternalServerError, "创建事务失败")
	ERROR_COMMIT_TX          = New(http.StatusInternalServerError, "提交事务失败")
	ERROR_EMPTY_ARGUMENTS    = New(http.StatusBadRequest, "用户名和密码不能为空")
	ERROR_INVALID_AUTH       = New(http.StatusBadRequest, "用户名或密码错误")
	ERROR_SESSION_EXPIRED    = New(http.StatusUnauthorized, "登录会话过期")
	ERROR_POPULATE_USER      = New(http.StatusBadRequest, "新增用户失败")
	ERROR_TYPE_MISMATCH      = New(http.StatusBadRequest, "类型和数据不匹配")
	ERROR_UPDAET_SETTING     = New(http.StatusInternalServerError, "修改设置项失败")
)
