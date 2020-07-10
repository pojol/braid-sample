package control

import (
	"braid-game/proto"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// GuestLogin 游客登录
func GuestLogin(token string, reqBody []byte) (interface{}, error) {
	res := proto.GuestLoginRes{}
	req := proto.GuestLoginReq{}
	var err error

	err = json.Unmarshal(reqBody, &req)
	if err != nil {
		return nil, err
	}

	fmt.Println("guest login token :", token)

	res.Token = "test_token_" + strconv.Itoa(int(time.Now().Unix()))

	return res, err
}
