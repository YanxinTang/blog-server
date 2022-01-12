package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/YanxinTang/blog-server/e"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh_Hans"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

type BadRequestMeta struct {
	Url string
}

var (
	uni *ut.UniversalTranslator
)

func init() {
	en := en.New()
	zh_Hans := zh_Hans.New()
	uni = ut.New(zh_Hans, en, zh_Hans)

	enTrans, _ := uni.GetTranslator("en")
	zhTrans, _ := uni.GetTranslator("zh_Hans")

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		en_translations.RegisterDefaultTranslations(v, enTrans)
		zh_translations.RegisterDefaultTranslations(v, zhTrans)
	}
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			lastError := c.Errors.Last().Err
			switch err := lastError.(type) {
			case e.ApiError:
				c.AbortWithStatusJSON(err.Code, gin.H{
					"message": err.Message,
				})
			case validator.ValidationErrors:
				languages := acceptedLanguages(c)
				trans, _ := uni.FindTranslator(languages...)
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"message": "参数异常",
					"errors":  err.Translate(trans),
				})

			default:
				log.Println(err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": "内部错误",
				})
			}
		}
	}
}

// acceptedLanguages returns an array of accepted languages denoted by
// the Accept-Language header sent by the browser
func acceptedLanguages(c *gin.Context) (languages []string) {
	accepted := c.GetHeader("Accept-Language")
	if accepted == "" {
		return
	}
	options := strings.Split(accepted, ",")
	l := len(options)
	languages = make([]string, l)

	for i := 0; i < l; i++ {
		locale := strings.SplitN(options[i], ";", 2)
		languages[i] = strings.Trim(locale[0], " ")
	}
	return
}
