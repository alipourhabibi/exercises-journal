package rss

import (
	"fmt"
	"io"
	"net/http"

	"github.com/alipourhabibi/exercises-journal/logging/logger"
)

type RssService struct {
	port   uint16
	logger *logger.Logger
}

func NewRssService(port uint16, logger *logger.Logger) (*RssService, error) {
	return &RssService{
		port:   port,
		logger: logger,
	}, nil
}

func (r *RssService) Run() error {
	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			r.logger.Error("error", err)
			return
		}
		r.logger.Info("body", string(body), "headers", req.Header)
	})
	r.logger.Info("msg", "Starting server", "port", r.port)
	return http.ListenAndServe(fmt.Sprintf(":%d", r.port), nil)
}
