package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ptechen/config"
	"github.com/ptechen/logger"
	"github.com/rs/zerolog"
)

type User struct {
	Name string `from:"name" binding:"required"`
}

type Config struct {
	Log *logger.LogParams
}

var log *zerolog.Logger

func main() {
	conf := config.Flag().SetConfigFileType("toml").SetConfigFileDir("../config").SetEnv("test")
	data := &Config{}
	conf.ParseFile(data)
	log = data.Log.InitParams().InitLogger()
	r := gin.New()
	r.Use(MiddleLogger(log))
	r.GET("/test", func(c *gin.Context) {
		var user User
		if c.ShouldBindQuery(&user) == nil {
			c.Set("params", user)
		}
	})
	r.Run(":8080")
}



