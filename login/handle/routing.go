package handle

import (
	"braid-game/login/control"
	"context"
	"encoding/json"

	"github.com/pojol/braid/module/rpc/client/bproto"
)

// RouteServer 路由服务器
type RouteServer struct {
	bproto.ListenServer
}

// RouteHandle 路由函数句柄
type RouteHandle func(reqBody []byte) (res interface{}, err error)

var (
	routeMap map[string]RouteHandle
)

// Routing 接收外界路由过来的（通常是gate）消息
func (rs *RouteServer) Routing(ctx context.Context, req *bproto.RouteReq) (*bproto.RouteRes, error) {
	res := new(bproto.RouteRes)

	if _, ok := routeMap[req.Service]; ok {
		ires, err := routeMap[req.Service](req.ReqBody)
		if err != nil {
			return nil, err
		}
		res.ResBody, err = json.Marshal(ires)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func init() {
	routeMap = make(map[string]RouteHandle)

	routeMap["guest"] = control.GuestLogin
}
