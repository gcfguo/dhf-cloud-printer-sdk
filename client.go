package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"dhf-cloud-printer-sdk/model"
)

type Client struct {
	serverURL  string
	email      string
	password   string
	httpClient *http.Client
	debug      bool
	token      atomic.Value
	stamp      atomic.Int64
}

func NewClient(opts ...ClientOption) *Client {
	defaultSettings := &ClientSettings{
		ServerURL:  "https://xprinter.96101210.com",
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		Debug:      false,
	}
	for _, opt := range opts {
		_ = opt.Apply(defaultSettings)
	}
	return &Client{
		serverURL:  defaultSettings.ServerURL,
		email:      defaultSettings.Email,
		password:   defaultSettings.Password,
		httpClient: defaultSettings.HTTPClient,
		debug:      defaultSettings.Debug,
		token:      atomic.Value{},
		stamp:      atomic.Int64{},
	}
}

func (c *Client) doRequest(authorize bool, method string, url string, body any, reader any) error {
	req, err := http.NewRequest(method, url, c.buildReqBody(body))
	if err != nil {
		return err
	}
	if authorize {
		var nowStamp = time.Now().Unix()
		var tokenType = "Bearer"
		var tokenValue string
		var ok bool
		a := time.Now().Add(time.Hour * -2).Unix()
		b := c.stamp.Load()
		if a != 0 && a < b {
			v := c.token.Load()
			tokenValue, ok = v.(string)
			if !ok {
				tokenValue, err = c.getAndLoadToken(nowStamp)
			}
		} else {
			tokenValue, err = c.getAndLoadToken(nowStamp)
		}
		if err != nil {
			return err
		}
		value := fmt.Sprintf("%s %s", tokenType, tokenValue)
		req.Header.Set("Authorization", value)
	}

	var res *http.Response
	for i := 0; i < 3; i++ {
		res, err = c.httpClient.Do(req)
		if err == nil {
			break
		}
		if _, ok := err.(net.Error); ok {
			continue
		}

		return err
	}

	if res == nil {
		return fmt.Errorf("invalid response")
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	_ = res.Body.Close()

	err = c.handleResponse(b, reader)
	if err != nil {
		return err
	}

	if c.debug {
		fmt.Printf("[pkg.transport] request_url:  %s\n", url)
		fmt.Printf("[pkg.transport] request_inf:  %s\n", jsonify(body))
		fmt.Printf("[pkg.transport] response_inf: %s\n", string(b))
	}

	return nil
}

func (c *Client) buildReqBody(body any) io.Reader {
	switch body.(type) {
	case string:
		return strings.NewReader(body.(string))

	case []byte:
		return bytes.NewReader(body.([]byte))

	default:
		b, _ := json.Marshal(body)
		return bytes.NewReader(b)
	}
}

func (c *Client) handleResponse(b []byte, reader any) error {
	type Response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data any    `json:"data"`
	}
	var resp Response
	resp.Data = reader
	err := json.Unmarshal(b, &resp)
	if err != nil {
		return err
	}

	if resp.Code != 0 {
		return fmt.Errorf(resp.Msg)
	}

	return nil
}

func (c *Client) getAndLoadToken(nowStamp int64) (string, error) {
	got, err := c.SignIn(&model.SignInReq{
		Email:    c.email,
		Password: c.password,
	})
	if err != nil {
		return "", err
	}
	c.token.Store(got.TokenValue)
	c.stamp.Store(nowStamp)
	return got.TokenValue, nil
}

// Verify 验证
func (c *Client) Verify(reqInf *model.VerifyReq) (*model.VerifyRes, error) {
	reqURL := c.concatURL("/v1/open/user/verify")
	var resInf model.VerifyRes
	err := c.doRequest(false, http.MethodPost, reqURL, reqInf, &resInf)
	if err != nil {
		return nil, err
	}

	return &resInf, nil
}

// SignUp 注册
// 注册前请先调用Verify获取验证码
func (c *Client) SignUp(reqInf *model.SignUpReq) (*model.SignUpRes, error) {
	reqURL := c.concatURL("/v1/open/user/signup")
	var resInf model.SignUpRes
	err := c.doRequest(false, http.MethodPost, reqURL, reqInf, &resInf)
	if err != nil {
		return nil, err
	}

	return &resInf, nil
}

// SignIn 登录
func (c *Client) SignIn(reqInf *model.SignInReq) (*model.SignInRes, error) {
	reqURL := c.concatURL("/v1/open/user/signin")
	var resInf model.SignInRes
	err := c.doRequest(false, http.MethodPost, reqURL, reqInf, &resInf)
	if err != nil {
		return nil, err
	}

	return &resInf, nil
}

// BindPrinter 绑定打印机
func (c *Client) BindPrinter(reqInf *model.PrinterBindReq) (*model.PrinterBindRes, error) {
	reqURL := c.concatURL("/v1/api/printer/bind")
	var resInf model.PrinterBindRes
	err := c.doRequest(true, http.MethodPost, reqURL, reqInf, nil)
	if err != nil {
		return nil, err
	}

	return &resInf, nil
}

// UnBindPrinter 解绑打印机
func (c *Client) UnBindPrinter(reqInf *model.PrinterUnBindReq) (*model.PrinterUnBindRes, error) {
	reqURL := c.concatURL("/v1/api/printer/unbind")
	var resInf model.PrinterUnBindRes
	err := c.doRequest(true, http.MethodPost, reqURL, reqInf, nil)
	if err != nil {
		return nil, err
	}

	return &resInf, nil
}

// ListPrinters 分页获取打印机列表
func (c *Client) ListPrinters(reqInf *model.ListPrinterReq) (*model.ListPrinterRes, error) {
	reqURL := c.concatURL("/v1/api/printer/list")
	var resInf model.ListPrinterRes
	err := c.doRequest(true, http.MethodPost, reqURL, reqInf, &resInf)
	if err != nil {
		return nil, err
	}

	return &resInf, nil
}

// Print 打印
func (c *Client) Print(reqInf *model.PrintReq) (*model.PrintRes, error) {
	reqURL := c.concatURL("/v1/api/printer/print")
	var resInf model.PrintRes
	err := c.doRequest(true, http.MethodPost, reqURL, reqInf, &resInf)
	if err != nil {
		return nil, err
	}

	return &resInf, nil
}

// ClearToken 清除当前授权的token
// 当你的授权信息发生变化时,请先清除token
func (c *Client) ClearToken() {
	c.token.Store("")
	c.stamp.Store(0)
}

func (c *Client) concatURL(uri string) string {
	return c.serverURL + uri
}

func jsonify(any any) string {
	b, _ := json.Marshal(any)
	return string(b)
}
