package handle

import (
	"braid-game/login/control"
	"context"
	"encoding/json"
	"fmt"

	"github.com/pojol/braid/plugin/grpcclient/bproto"
)

// RouteServer 路由服务器
type RouteServer struct {
	bproto.ListenServer
}

// RouteHandle 路由函数句柄
type RouteHandle func(ctx context.Context, token string, reqBody []byte) (res interface{}, err error)

var (
	routeMap map[string]RouteHandle
)

// Routing 接收外界路由过来的（通常是gate）消息
func (rs *RouteServer) Routing(ctx context.Context, req *bproto.RouteReq) (*bproto.RouteRes, error) {
	res := new(bproto.RouteRes)
	fmt.Println("routing", req.Service)

	if _, ok := routeMap[req.Service]; ok {
		ires, err := routeMap[req.Service](ctx, req.Token, req.ReqBody)
		if err != nil {
			return nil, err
		}

		if ires == nil {
			res.ResBody = []byte{}
		} else {
			res.ResBody, err = json.Marshal(ires)
			if err != nil {
				return nil, err
			}
		}

	}

	return res, nil
}

func init() {
	routeMap = make(map[string]RouteHandle)

	routeMap["guest"] = control.GuestLogin
	routeMap["out"] = control.Loginout
}
