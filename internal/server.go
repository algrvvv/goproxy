package internal

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
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

func (s *Server) basicAuth(r *http.Request) bool {
	auth := r.Header.Get("Proxy-Authorization")
	if auth == "" {
		return false
	}

	payload, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(auth, "Basic "))
	if err != nil {
		return false
	}

	// получаем из username:password -> [0] -> username, [1] -> password
	parts := strings.SplitN(string(payload), ":", 2)
	if len(parts) != 2 || !s.authenticate(parts[0], parts[1]) {
		return false
	}

	return true
}

func (s *Server) handleHTTP(w http.ResponseWriter, r *http.Request) {
	if !s.basicAuth(r) {
		w.Header().Set("Proxy-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "Proxy authentication required", http.StatusProxyAuthRequired)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest(r.Method, r.URL.String(), r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req.Header = r.Header
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	for k, v := range resp.Header {
		w.Header()[k] = v
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (s *Server) handleHTTPS(w http.ResponseWriter, r *http.Request) {
	if !s.basicAuth(r) {
		w.Header().Set("Proxy-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "Proxy authentication required", http.StatusProxyAuthRequired)
		return
	}

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "hijacker not supported", http.StatusInternalServerError)
		return
	}

	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer clientConn.Close()

	servConn, err := net.Dial("tcp", r.Host)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer servConn.Close()

	clientConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))

	go io.Copy(servConn, clientConn)
	io.Copy(clientConn, servConn)
}

func (s *Server) handleProxy(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodConnect {
		s.handleHTTPS(w, r)
	} else {
		s.handleHTTP(w, r)
	}
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
	http.HandleFunc("/", s.handleProxy)
	log.Println("Starting proxy server on " + s.port)
}
