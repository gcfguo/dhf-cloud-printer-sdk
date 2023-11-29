package core

import "net/http"

type ClientOption interface {
	Apply(settings *ClientSettings) error
}

type httpClientOption struct {
	HttpClient *http.Client
}

func (o *httpClientOption) Apply(settings *ClientSettings) error {
	settings.HTTPClient = o.HttpClient
	return nil
}

// WithHttpClient 使用自定义的*http.Client
func WithHttpClient(httpClient *http.Client) ClientOption {
	return &httpClientOption{HttpClient: httpClient}
}

type instantAuthOption struct {
	Email    string
	Password string
}

func (o *instantAuthOption) Apply(settings *ClientSettings) error {
	settings.Email = o.Email
	settings.Password = o.Password
	return nil
}

// WithInstantAuth 使用
func WithInstantAuth(email, password string) ClientOption {
	return &instantAuthOption{Email: email, Password: password}
}

type serverURLOption struct {
	ServerURL string
}

func (o *serverURLOption) Apply(settings *ClientSettings) error {
	settings.ServerURL = o.ServerURL
	return nil
}

func WithServerURL(serverURL string) ClientOption {
	return &serverURLOption{ServerURL: serverURL}
}
