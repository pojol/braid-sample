package proto

// RenameReq 更新用户名请求
type RenameReq struct {
	accountid string
	nickname  string
}

// RenameRes 更新用户名返回
type RenameRes struct {
	nickname string
}
