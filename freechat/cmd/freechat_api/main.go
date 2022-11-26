package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"

	freechat_user "Open_IM/freechat/internal/freechat/user"
	"Open_IM/open_im/pkg/common/config"
	"Open_IM/open_im/pkg/common/constant"
	"Open_IM/open_im/pkg/common/log"
	"Open_IM/open_im/pkg/utils"

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
	fmt.Println("start freechat api server, port: ", *ginPort)
	address := "0.0.0.0:" + strconv.Itoa(*ginPort)
	if config.FreechatConfig.Api.ListenIP != "" {
		address = config.FreechatConfig.Api.ListenIP + ":" + strconv.Itoa(*ginPort)
	}
	address = config.FreechatConfig.CmsApi.ListenIP + ":" + strconv.Itoa(*ginPort)
	fmt.Println("start freechat api server address: ", address)
	err := r.Run(address)
	if err != nil {
		log.Error("", "run failed ", *ginPort, err.Error())
	}
}
