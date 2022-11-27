package register

import (
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	http2 "Open_IM/pkg/common/http"
	"Open_IM/pkg/common/log"
	pbFriend "Open_IM/pkg/proto/friend"
	"Open_IM/pkg/utils"
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ParamsRegister struct {
	UserID           string `gorm:"column:userID;primary_key;size:64" json:"userID" binding:"required,eth_addr|btc_addr|btc_addr_bech32"`
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
	InvitationCode   string `json:"invitationCode"`
	Gender           int32  `json:"gender"`
	Birth            int32  `json:"birth"`
}

// Register register in uses and registers table
func Register(c *gin.Context) {
	params := ParamsRegister{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError(params.OperationID, utils.GetSelfFuncName(), "bind json failed", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": constant.FormattingError, "errMsg": err.Error()})
		return
	}
	ip := c.Request.Header.Get("X-Forward-For")
	if ip == "" {
		ip = c.ClientIP()
	}
	log.NewDebug(params.OperationID, utils.GetSelfFuncName(), "ip:", ip)
	Limited, LimitError := imdb.IsLimitRegisterIp(ip)
	if LimitError != nil {
		log.Error(params.OperationID, utils.GetSelfFuncName(), LimitError, ip)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": constant.ErrDB.ErrCode, "errMsg": LimitError.Error()})
		return
	}
	if Limited {
		log.NewInfo(params.OperationID, utils.GetSelfFuncName(), "is limited", ip, "params:", params)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": constant.RegisterLimit, "errMsg": "limited"})
		return
	}
	if params.Nickname == "" {
		rand.Seed(time.Now().UnixNano())
		rd := 100000 + rand.Intn(900000)
		params.Nickname = "freechat" + strconv.Itoa(rd)
	}
	account := params.Email
	if (config.Config.Demo.UseSuperCode && params.VerificationCode != config.Config.Demo.SuperCode) ||
		!config.Config.Demo.UseSuperCode {
		accountKey := account + "_" + constant.VerificationCodeForRegisterSuffix
		v, err := db.DB.GetAccountCode(accountKey)
		if err != nil || v != params.VerificationCode {
			log.NewError(params.OperationID, "Verification code error", account, params.VerificationCode)
			c.JSON(http.StatusOK, gin.H{"errCode": constant.CodeInvalidOrExpired, "errMsg": "Verification code error!"})
			return
		}
	}

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
	openIMRegisterReq.Gender = params.Gender
	openIMRegisterReq.Birth = uint32(params.Birth)
	openIMRegisterResp := api.UserRegisterResp{}
	log.NewDebug(params.OperationID, utils.GetSelfFuncName(), "register req:", openIMRegisterReq)
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

	err = imdb.SetPassword(account, params.Password, params.Ex, userID, params.AreaCode, ip)
	if err != nil {
		log.NewError(params.OperationID, "set phone number password error", account, "err", err.Error())
		c.JSON(http.StatusOK, gin.H{"errCode": constant.RegisterFailed, "errMsg": "register failed"})
		return
	}
	if config.Config.Demo.NeedInvitationCode && params.InvitationCode != "" {
		//判断一下验证码的使用情况
		LockSucc := imdb.TryLockInvitationCode(params.InvitationCode, userID)
		if LockSucc {
			imdb.FinishInvitationCode(params.InvitationCode, userID)
		}
	}
	if err := imdb.InsertIpRecord(userID, ip); err != nil {
		log.NewError(params.OperationID, utils.GetSelfFuncName(), userID, ip, err.Error())
	}
	log.Info(params.OperationID, "end InsertRegister", account, userID)
	// demo onboarding
	if params.UserID == "" && config.Config.Demo.OnboardProcess {
		select {
		case Ch <- OnboardingProcessReq{
			OperationID: params.OperationID,
			UserID:      userID,
			NickName:    params.Nickname,
			FaceURL:     params.FaceURL,
			PhoneNumber: params.AreaCode + params.PhoneNumber,
			Email:       params.Email,
		}:
		case <-time.After(time.Second * 2):
			log.NewWarn(params.OperationID, utils.GetSelfFuncName(), "to ch timeOut")
		}
	}

	// register add friend
	select {
	case ChImportFriend <- &pbFriend.ImportFriendReq{
		OperationID: params.OperationID,
		FromUserID:  userID,
		OpUserID:    config.Config.Manager.AppManagerUid[0],
	}:
	case <-time.After(time.Second * 2):
		log.NewWarn(params.OperationID, utils.GetSelfFuncName(), "to ChImportFriend timeOut")
	}

	c.JSON(http.StatusOK, gin.H{"errCode": constant.NoError, "errMsg": "", "data": openIMRegisterResp.UserToken})
	return
}
