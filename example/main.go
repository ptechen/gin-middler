package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ptechen/config"
	"github.com/ptechen/gin-middler/middler"
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
	conf := config.Flag().SetConfigFileType("toml").SetConfigFileDir("./config").SetEnv("test")
	data := &Config{}
	conf.ParseFile(data)
	log = data.Log.InitParams().InitLogger()
	r := gin.New()
	r.Use(middler.MiddleLogger(log, gin.LoggerConfig{SkipPaths:[]string{"fsfs", "fsfsf"}}))
	r.POST("/test", func(c *gin.Context) {
		var user User
		if c.ShouldBind(&user) == nil {
			c.Set("req", user)
			c.Set("res", user)
		} else {

		}
	})
	r.Run(":8080")
}
