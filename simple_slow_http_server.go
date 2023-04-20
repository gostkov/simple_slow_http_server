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
	logger *zap.SugaredLogger
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
	service := &Service{logger: logger}
	mux := http.NewServeMux()
	mux.HandleFunc("/slow/", service.slowRequest)
	mux.HandleFunc("/fast/", service.fastRequest)
	mux.HandleFunc("/error/", service.errorRequest)
	logger.Infoln("Starting http service")
	err = http.ListenAndServe(httpIp+":"+httpPort, mux)
	if err != nil {
		logger.Errorln("Can't run service: ", err)
		return
	}
}
