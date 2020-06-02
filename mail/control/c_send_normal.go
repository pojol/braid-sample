package control

import (
	"braid-game/errcode"
	"braid-game/mail/model"
	"fmt"
)

// SendNormal 发送一封默认邮件给指定玩家
func SendNormal(accid string, body model.MailDat) errcode.Err {

	fmt.Println("send mail ", body.Title)

	return errcode.Succ
}
