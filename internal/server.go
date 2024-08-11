package internal

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/elazarl/goproxy"
)

type Server struct {
	port  string            // порт работы сервера
	users map[string]string // пользователи полученные из конфига для аутентификации
}

func (s *Server) authenticate(username, password string) bool {
	if validPass, ok := s.users[username]; ok {
		return validPass == password
	}
	return false
}

func (s *Server) proxyHandler(proxy *goproxy.ProxyHttpServer) {
	proxy.OnRequest().HandleConnectFunc(func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
		auth := ctx.Req.Header.Get("Proxy-Authorization")
		if auth == "" {
			ctx.Resp = &http.Response{
				StatusCode: http.StatusProxyAuthRequired,
				Header:     make(http.Header),
				Body:       http.NoBody,
			}
			ctx.Resp.Header.Set("Proxy-Authenticate", `Basic realm="Restricted"`)
			return goproxy.RejectConnect, host
		}

		payload, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(auth, "Basic "))
		if err != nil {
			return goproxy.RejectConnect, host
		}

		parts := strings.SplitN(string(payload), ":", 2)
		if len(parts) != 2 || !s.authenticate(parts[0], parts[1]) {
			return goproxy.RejectConnect, host
		}

		return goproxy.OkConnect, host
	})

	proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		auth := req.Header.Get("Proxy-Authorization")
		if auth == "" {
			resp := &http.Response{
				StatusCode: http.StatusProxyAuthRequired,
				Header:     make(http.Header),
				Body:       http.NoBody,
			}
			resp.Header.Set("Proxy-Authenticate", `Basic realm="Restricted"`)
			return req, resp
		}

		payload, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(auth, "Basic "))
		if err != nil {
			return req, nil
		}

		parts := strings.SplitN(string(payload), ":", 2)
		if len(parts) != 2 || !s.authenticate(parts[0], parts[1]) {
			return req, &http.Response{
				StatusCode: http.StatusProxyAuthRequired,
				Header:     make(http.Header),
				Body:       http.NoBody,
			}
		}

		return req, nil
	})
}

func NewServer() (*Server, error) {
	config, err := LoadConfig("config.yaml")
	if err != nil {
		return nil, err
	}

	var server Server
	server.users = make(map[string]string)
	for _, user := range config.Server.Users {
		server.users[user.Username] = user.Password
	}
	server.port = fmt.Sprintf(":%d", config.Server.Port)

	return &server, nil
}

func (s *Server) Run() {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true
	s.proxyHandler(proxy)

	log.Println("Starting proxy server on " + s.port)
	if err := http.ListenAndServe(s.port, proxy); err != nil {
		log.Fatal(err)
	}
}
