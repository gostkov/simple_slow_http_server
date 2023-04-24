package main

import (
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Service struct {
	logger            *zap.SugaredLogger
	basicAuthLogin    string
	basicAuthPassword string
}

func (s *Service) slowRequest(w http.ResponseWriter, r *http.Request) {
	timeout := r.URL.Query().Get("timeout")
	defaultTimeout := 10
	atoi, err := strconv.Atoi(timeout)
	if err != nil {
		s.logger.Warnf("Can't convert timeout to int, use default %d sec", defaultTimeout)
		atoi = defaultTimeout
	}
	s.logger.Infof("slowRequest, waiting %d sec", atoi)
	time.Sleep(time.Duration(atoi) * time.Second)
	w.Write([]byte("is slow response"))
	return
}

func (s *Service) fastRequest(w http.ResponseWriter, r *http.Request) {
	s.logger.Infoln("fastRequest, doesn't wait")
	w.Write([]byte("is fast response"))
	return
}

func (s *Service) errorRequest(w http.ResponseWriter, r *http.Request) {
	errorCode := r.URL.Query().Get("code")
	defaultError := http.StatusBadRequest
	atoi, err := strconv.Atoi(errorCode)
	if err != nil || (atoi < 100 && atoi > 599) {
		s.logger.Warnf("Not valid error_code, use default %d", defaultError)
		atoi = defaultError
	}
	w.WriteHeader(atoi)
	w.Write([]byte("is error response"))
	return
}

func (s *Service) basicAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.basicAuthPassword == "" || s.basicAuthLogin == "" {
			s.logger.Debugln("basic auth is not enabled")
			next.ServeHTTP(w, r)
			return
		}
		username, password, ok := r.BasicAuth()
		if ok {
			if username == s.basicAuthLogin && password == s.basicAuthPassword {
				s.logger.Debugln("basic auth has been passed")
				next.ServeHTTP(w, r)
				return
			}
		}
		s.logger.Debugln("basic auth has not been passed")
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
}

func main() {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	logger := zapLogger.Sugar()
	httpIp := os.Getenv("IP")
	if httpIp == "" {
		logger.Warnln("Use default IP 0.0.0.0")
		httpIp = "0.0.0.0"
	}
	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		logger.Warnln("Use default port 8080")
		httpPort = "8080"
	}
	authBasicLogin := os.Getenv("BASIC_AUTH_LOGIN")
	authBasicPassword := os.Getenv("BASIC_AUTH_PASSWORD")

	service := &Service{
		logger:            logger,
		basicAuthLogin:    authBasicLogin,
		basicAuthPassword: authBasicPassword,
	}
	mux := http.NewServeMux()

	mux.HandleFunc("/slow/", service.basicAuth(service.slowRequest))
	mux.HandleFunc("/fast/", service.basicAuth(service.fastRequest))
	mux.HandleFunc("/error/", service.basicAuth(service.errorRequest))
	logger.Infoln("Starting http service")
	err = http.ListenAndServe(httpIp+":"+httpPort, mux)
	if err != nil {
		logger.Errorln("Can't run service: ", err)
		return
	}
}
