package register

import (
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	http2 "Open_IM/pkg/common/http"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ParamsLogin struct {
	// 加一个userID donedone
	// 加一个验证码 donedone
	VerificationCode string `json:"verificationCode" binding:"required"`
	UserID           string `json:"userID" binding:"required"`
	Email            string `json:"email"`
	PhoneNumber      string `json:"phoneNumber"`
	Password         string `json:"password"`
	Platform         int32  `json:"platform"`
	OperationID      string `json:"operationID" binding:"required"`
	AreaCode         string `json:"areaCode"`
}

func Login(c *gin.Context) {
	params := ParamsLogin{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": constant.FormattingError, "errMsg": err.Error()})
		return
	}
	// 判断邮箱是否存在 donedone
	// TODO 调试的时候看下uid不存在是报err还是咋地
	eml, err := im_mysql_model.GetEmail(params.UserID)
	if err != nil || eml.Email == "" {
		c.JSON(http.StatusOK, gin.H{"errCode": constant.NotRegistered, "errMsg": "The Email has not been registered"})
		return
	}

	account := params.Email
	// 登录需要验证邮箱 donedone
	if params.VerificationCode != config.Config.Demo.SuperCode {
		accountKey := account + "_" + constant.VerificationCodeForLoginSuffix
		v, err := db.DB.GetAccountCode(accountKey)
		if err != nil || v != params.VerificationCode {
			log.NewError(params.OperationID, "password Verification code error", account, params.VerificationCode)
			c.JSON(http.StatusOK, gin.H{"errCode": constant.CodeInvalidOrExpired, "errMsg": "Verification code error!"})
			return
		}
	}
	r, err := im_mysql_model.GetRegister(account, params.AreaCode)
	if err != nil {
		log.NewError(params.OperationID, "user have not register", params.Password, account, err.Error())
		c.JSON(http.StatusOK, gin.H{"errCode": constant.NotRegistered, "errMsg": "Mobile phone number is not registered"})
		return
	}
	if r.Password != params.Password {
		log.NewError(params.OperationID, "password  err", params.Password, account, r.Password, r.Account)
		c.JSON(http.StatusOK, gin.H{"errCode": constant.PasswordErr, "errMsg": "password err"})
		return
	}
	var userID string
	if r.UserID != "" {
		userID = r.UserID
	} else {
		userID = r.Account
	}
	url := fmt.Sprintf("http://%s:%d/auth/user_token", utils.ServerIP, config.Config.Api.GinPort[0])
	openIMGetUserToken := api.UserTokenReq{}
	openIMGetUserToken.OperationID = params.OperationID
	openIMGetUserToken.Platform = params.Platform
	openIMGetUserToken.Secret = config.Config.Secret
	openIMGetUserToken.UserID = userID
	openIMGetUserTokenResp := api.UserTokenResp{}
	bMsg, err := http2.Post(url, openIMGetUserToken, 2)
	if err != nil {
		log.NewError(params.OperationID, "request openIM get user token error", account, "err", err.Error())
		c.JSON(http.StatusOK, gin.H{"errCode": constant.GetIMTokenErr, "errMsg": err.Error()})
		return
	}
	err = json.Unmarshal(bMsg, &openIMGetUserTokenResp)
	if err != nil || openIMGetUserTokenResp.ErrCode != 0 {
		log.NewError(params.OperationID, "request get user token", account, "err", "")
		c.JSON(http.StatusOK, gin.H{"errCode": constant.GetIMTokenErr, "errMsg": ""})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errCode": constant.NoError, "errMsg": "", "data": openIMGetUserTokenResp.UserToken})

}
