package register

import (
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/utils"
	"fmt"
	"freechat/pkg/common/config"
	"freechat/pkg/common/constant"
	"freechat/pkg/common/db/mysql_model/im_mysql_model"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
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
	UserID         string `json:"userID" binding:"required,eth_addr|btc_addr|btc_addr_bech32"`
	Email          string `json:"email" binding:"email"`
	OperationID    string `json:"operationID" binding:"required"`
	UsedFor        int    `json:"usedFor"`
	AreaCode       string `json:"areaCode"`
	InvitationCode string `json:"invitationCode"`
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
		r, err := im_mysql_model.GetEmail(params.UserID)
		if err != nil || r.Email == "" {
			c.JSON(http.StatusOK, gin.H{"errCode": constant.NotRegistered, "errMsg": err.Error()})
			return
		}
		account = r.Email
	}
	if account == "" {
		c.JSON(http.StatusOK, gin.H{"errCode": constant.MailSendCodeErr, "errMsg": "The email is empty"})
		return
	}
	var accountKey = account
	switch params.UsedFor {
	case constant.VerificationCodeForRegister:
		_, err := im_mysql_model.GetRegister(account, "", params.UserID)
		if err == nil {
			log.NewError(params.OperationID, "The account has been registered", params)
			c.JSON(http.StatusOK, gin.H{"errCode": constant.HasRegistered, "errMsg": "The email has been registered"})
			return
			//需要邀请码
			if config.Config.Demo.NeedInvitationCode {
				err = im_mysql_model.CheckInvitationCode(params.InvitationCode)
				if err != nil {
					log.NewError(params.OperationID, "邀请码错误", params)
					c.JSON(http.StatusOK, gin.H{"errCode": constant.InvitationError, "errMsg": "邀请码错误"})
					return
				}
			}
			accountKey = accountKey + "_" + constant.VerificationCodeForRegisterSuffix
			ok, err := db.DB.JudgeAccountEXISTS(accountKey)
			if ok || err != nil {
				log.NewError(params.OperationID, "Repeat send code", params, accountKey)
				c.JSON(http.StatusOK, gin.H{"errCode": constant.RepeatSendCode, "errMsg": "Repeat send code"})
				return
			}
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
	m.SetHeader(`Subject`, fmt.Sprintf("%s %d", config.Config.Demo.Mail.Title, code))
	m.SetBody(`text/html`, fmt.Sprintf("%s %d", config.Config.Demo.Mail.Content, code))
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
