package middler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gopkg.in/go-playground/validator.v9"
	"strings"
)

//type User struct {
//	Test string `json:"test"`
//}
//
//func MatchStruct(handlerName string) (res interface{}) {
//	if handlerName == "main.test" {
//		return &User{}
//	} else {
//		return
//	}
//}

func ParseParams(matchStruct func(name string) interface{}, logger *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		params := matchStruct(c.HandlerName())
		if err := c.ShouldBind(&params); err != nil {
			logger.Err(err).
				Str("client_ip", c.ClientIP()).
				Str("method", c.Request.Method).
				Str("handler", c.HandlerName()).
				Msgf("parse %#v failed", &params)
			c.Abort()
		}
		err := validate(params)
		if err != nil {
			logger.Err(err).
				Str("client_ip", c.ClientIP()).
				Str("method", c.Request.Method).
				Str("handler", c.HandlerName()).
				Msgf("parse %#v failed", &params)
			c.Abort()
		}
		c.Set("params", params)
		c.Next()
		// after request
	}
}

var vd *validator.Validate

// validate:"required"
func validate(req interface{}) (err error) {
	if err = vd.Struct(req); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return errors.New("InvalidValidationError")
		}
		var fields []string
		for _, errF := range err.(validator.ValidationErrors) {
			fields = append(fields, errF.Field())
		}
		return fmt.Errorf("ValidationErrors: [%s]", strings.Join(fields, ","))
	}
	return
}
