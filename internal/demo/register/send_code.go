package register

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"math/rand"
	"net/http"

	"time"
)

var sms SMS

func init() {
	var err error
	if config.Config.Demo.AliSMSVerify.Enable {
		sms, err = NewAliSMS()
		if err != nil {
			panic(err)
		}
	} else {
		sms, err = NewTencentSMS()
		if err != nil {
			panic(err)
		}
	}
}

type paramsVerificationCode struct {
	// 加一个userID donedone
	UserID      string `json:"userID" binding:"required"`
	Email       string `json:"email"`
	OperationID string `json:"operationID" binding:"required"`
	UsedFor     int    `json:"usedFor"`
}

func SendVerificationCode(c *gin.Context) {
	params := paramsVerificationCode{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("", "BindJSON failed", "err:", err.Error(), "userID", params.UserID)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": constant.FormattingError, "errMsg": err.Error()})
		return
	}
	operationID := params.OperationID
	if operationID == "" {
		operationID = utils.OperationIDGenerator()
	}
	if params.UsedFor == 0 {
		params.UsedFor = constant.VerificationCodeForLogin
	}
	var account string
	if params.Email != "" && params.UsedFor == constant.VerificationCodeForRegister {
		account = params.Email
	} else {
		// 发验证码的时候前端传给后端userId，后端自己去查邮箱然后发 donedone
		r, err := im_mysql_model.GetEmail(params.UserID)
		if err != nil || r.Email == "" {
			c.JSON(http.StatusOK, gin.H{"errCode": constant.NotRegistered, "errMsg": err.Error()})
		}
		account = r.Email
	}
	if account == "" {
		c.JSON(http.StatusOK, gin.H{"errCode": constant.MailSendCodeErr, "errMsg": "The email is empty"})
	}
	// 修改验证码存储的key  donedone
	var accountKey = account
	switch params.UsedFor {
	case constant.VerificationCodeForRegister:
		_, err := im_mysql_model.GetRegister(account, "")
		if err == nil {
			log.NewError(params.OperationID, "The email has been registered", params)
			c.JSON(http.StatusOK, gin.H{"errCode": constant.HasRegistered, "errMsg": "The email has been registered"})
			return
		}
		accountKey = accountKey + "_" + constant.VerificationCodeForRegisterSuffix
		ok, err := db.DB.JudgeAccountEXISTS(accountKey)
		if ok || err != nil {
			log.NewError(params.OperationID, "Repeat send code", params, accountKey)
			c.JSON(http.StatusOK, gin.H{"errCode": constant.RepeatSendCode, "errMsg": "Repeat send code"})
			return
		}
	case constant.VerificationCodeForReset:
		accountKey = accountKey + "_" + constant.VerificationCodeForResetSuffix
		ok, err := db.DB.JudgeAccountEXISTS(accountKey)
		if ok || err != nil {
			log.NewError(params.OperationID, "Repeat send code", params, accountKey)
			c.JSON(http.StatusOK, gin.H{"errCode": constant.RepeatSendCode, "errMsg": "Repeat send code"})
			return
		}
	case constant.VerificationCodeForLogin:
		accountKey = accountKey + "_" + constant.VerificationCodeForLoginSuffix
		ok, err := db.DB.JudgeAccountEXISTS(accountKey)
		if ok || err != nil {
			log.NewError(params.OperationID, "Repeat send code", params, accountKey)
			c.JSON(http.StatusOK, gin.H{"errCode": constant.RepeatSendCode, "errMsg": "Repeat send code"})
			return
		}

	}
	rand.Seed(time.Now().UnixNano())
	code := 100000 + rand.Intn(900000)
	log.NewInfo(params.OperationID, params.UsedFor, "begin store redis", accountKey, code)
	err := db.DB.SetAccountCode(accountKey, code, config.Config.Demo.CodeTTL)
	if err != nil {
		log.NewError(params.OperationID, "set redis error", accountKey, "err", err.Error())
		c.JSON(http.StatusOK, gin.H{"errCode": constant.SmsSendCodeErr, "errMsg": "Enter the superCode directly in the verification code box, SuperCode can be configured in config.xml"})
		return
	}
	log.NewDebug(params.OperationID, config.Config.Demo)

	m := gomail.NewMessage()
	m.SetHeader(`From`, config.Config.Demo.Mail.SenderMail)
	m.SetHeader(`To`, []string{account}...)
	m.SetHeader(`Subject`, config.Config.Demo.Mail.Title)
	m.SetBody(`text/html`, fmt.Sprintf("%d", code))
	if err := gomail.NewDialer(config.Config.Demo.Mail.SmtpAddr, config.Config.Demo.Mail.SmtpPort, config.Config.Demo.Mail.SenderMail, config.Config.Demo.Mail.SenderAuthorizationCode).DialAndSend(m); err != nil {
		log.Error(params.OperationID, "send mail error", account, err.Error())
		c.JSON(http.StatusOK, gin.H{"errCode": constant.MailSendCodeErr, "errMsg": ""})
		return
	}
	log.Debug(params.OperationID, "send sms success", code, accountKey)
	data := make(map[string]interface{})
	data["account"] = account
	c.JSON(http.StatusOK, gin.H{"errCode": constant.NoError, "errMsg": "Verification code has been set!", "data": data})
}
