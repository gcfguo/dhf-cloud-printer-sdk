package core

import "net/http"

type ClientSettings struct {
	//ServerURL 服务端地址
	ServerURL string
	//Email 邮箱号
	Email string
	//Password 密码
	Password string
	//HTTPClient http客户端请求对象
	HTTPClient *http.Client
	//Debug 调试开关
	Debug bool
}
