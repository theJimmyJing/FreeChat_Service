package main

import (
	freechat_user "Open_IM/internal/freechat/user"
	"Open_IM/pkg/utils"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"

	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"github.com/gin-gonic/gin"
)

func main() {
	log.NewPrivateLog(constant.LogFileName)
	gin.SetMode(gin.ReleaseMode)
	f, _ := os.Create("../logs/freechat_api.log")
	gin.DefaultWriter = io.MultiWriter(f)

	r := gin.Default()
	r.Use(utils.CorsHandler())

	fUserRouterGroup := r.Group("/user")
	{
		fUserRouterGroup.POST("/select_face", freechat_user.SelectFace)
	}

	defaultPorts := config.FreechatConfig.Freechat.Port
	ginPort := flag.Int("port", defaultPorts[0], "get ginServerPort from cmd,default 10005 as port")
	flag.Parse()
	fmt.Println("start demo api server, port: ", *ginPort)
	address := "0.0.0.0:" + strconv.Itoa(*ginPort)
	if config.FreechatConfig.Api.ListenIP != "" {
		address = config.FreechatConfig.Api.ListenIP + ":" + strconv.Itoa(*ginPort)
	}
	address = config.FreechatConfig.CmsApi.ListenIP + ":" + strconv.Itoa(*ginPort)
	fmt.Println("start demo api server address: ", address)
	err := r.Run(address)
	if err != nil {
		log.Error("", "run failed ", *ginPort, err.Error())
	}
}
