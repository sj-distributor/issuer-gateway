package main

import (
	"cert-gateway/cert/internal/configs"
	"cert-gateway/utils"
	"flag"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func main() {

	var configFile = flag.String("f", "internal/configs/config.yaml", "the config file")

	utils.MustLoad(configFile, configs.C)

	//database.Init(configs.C)

	engine := gin.Default()

	engine.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "*"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "*"},
		AllowCredentials: false,
		AllowAllOrigins:  true,
		MaxAge:           12 * time.Hour,
	}))

	engine.GET("/cert", func(c *gin.Context) {

		certificate, err := utils.ReqCertificate(configs.C.Acme.Email, "anson.itst.cn")

		c.JSON(200, gin.H{
			"cert": certificate,
			"err":  err,
		})

	})

	engine.NoRoute(func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": 404,
			"msg":  "Not Found",
		})
	})

	if err := utils.GraceFul(time.Minute, &http.Server{
		Addr:    ":9527",
		Handler: engine,
	}).ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}
