package utils

import (
	"fmt"
	"strconv"

	"github.com/YanxinTang/blog-server/internal/pkg/e"
	"github.com/YanxinTang/blog-server/internal/pkg/log"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetID(c *gin.Context, param string) (int, error) {
	param = c.Param(param)
	id, err := strconv.Atoi(param)
	if err != nil {
		log.Warn(fmt.Sprintf("failing converting %s", param), zap.String("param", param))
		c.Error(e.ERROR_RESOURCE_NOT_FOUND)
		return 0, err
	}
	return id, nil
}
