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
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type ParamsRegister struct {
	UserID           string `gorm:"column:userID;primary_key;size:64" json:"userID" binding:"required"`
	Email            string `json:"email" binding:"required,email"`
	VerificationCode string `json:"verificationCode" binding:"required"`
	Nickname         string `json:"nickname"`
	PhoneNumber      string `json:"phoneNumber"`
	Password         string `json:"password"`
	Platform         int32  `json:"platform" binding:"required,min=1,max=7"`
	Ex               string `json:"ex"`
	FaceURL          string `json:"faceURL"`
	OperationID      string `json:"operationID" binding:"required"`
	AreaCode         string `json:"areaCode"`
}

// Register register in uses and registers table
func Register(c *gin.Context) {
	params := ParamsRegister{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError(params.OperationID, utils.GetSelfFuncName(), "bind json failed", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": constant.FormattingError, "errMsg": err.Error()})
		return
	}
	if params.Nickname == "" {
		rand.Seed(time.Now().UnixNano())
		rd := 100000 + rand.Intn(900000)
		params.Nickname = "freechat" + strconv.Itoa(rd)
	}
	account := params.Email
	if params.VerificationCode != config.Config.Demo.SuperCode {
		accountKey := account + "_" + constant.VerificationCodeForRegisterSuffix
		v, err := db.DB.GetAccountCode(accountKey)
		if err != nil || v != params.VerificationCode {
			log.NewError(params.OperationID, "password Verification code error", account, params.VerificationCode)
			c.JSON(http.StatusOK, gin.H{"errCode": constant.CodeInvalidOrExpired, "errMsg": "Verification code error!"})
			return
		}
	}

	//userID := utils.Md5(params.OperationID + strconv.FormatInt(time.Now().UnixNano(), 10))
	//bi := big.NewInt(0)
	//bi.SetString(userID[0:8], 16)
	//userID = bi.String()
	userID := params.UserID
	url := config.Config.Demo.ImAPIURL + "/auth/user_register"
	openIMRegisterReq := api.UserRegisterReq{}
	openIMRegisterReq.OperationID = params.OperationID
	openIMRegisterReq.Platform = params.Platform
	openIMRegisterReq.UserID = userID
	openIMRegisterReq.Email = params.Email
	openIMRegisterReq.Nickname = params.Nickname
	openIMRegisterReq.Secret = config.Config.Secret
	openIMRegisterReq.FaceURL = params.FaceURL
	openIMRegisterResp := api.UserRegisterResp{}
	bMsg, err := http2.Post(url, openIMRegisterReq, 2)
	if err != nil {
		log.NewError(params.OperationID, "request openIM register error", account, "err", err.Error())
		c.JSON(http.StatusOK, gin.H{"errCode": constant.RegisterFailed, "errMsg": "register failed"})
		return
	}
	err = json.Unmarshal(bMsg, &openIMRegisterResp)
	if err != nil || openIMRegisterResp.ErrCode != 0 {
		log.NewError(params.OperationID, "request openIM register error", account, "err", "resp: ",
			openIMRegisterResp.ErrCode, openIMRegisterResp.ErrMsg)
		if err != nil {
			log.NewError(params.OperationID, utils.GetSelfFuncName(), err.Error())
		}
		c.JSON(http.StatusOK, gin.H{"errCode": constant.RegisterFailed, "errMsg": "register failed"})
		return
	}
	log.Info(params.OperationID, "begin store mysql", account, params.Password, "info", params.FaceURL, params.Nickname)

	err = im_mysql_model.InsertRegister(account, params.Password, params.Ex, userID, params.AreaCode)
	if err != nil {
		log.NewError(params.OperationID, "set phone number password error", account, "err", err.Error())
		c.JSON(http.StatusOK, gin.H{"errCode": constant.RegisterFailed, "errMsg": "register failed"})
		return
	}
	log.Info(params.OperationID, "end InsertRegister", account, userID)
	// demo onboarding
	onboardingProcess(params.OperationID, userID, params.Nickname, params.FaceURL, params.AreaCode+params.PhoneNumber, params.Email)
	c.JSON(http.StatusOK, gin.H{"errCode": constant.NoError, "errMsg": "", "data": openIMRegisterResp.UserToken})
	return
}
