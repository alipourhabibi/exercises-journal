package server

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type ServerConfiguration func(*ServerService) error

type ServerService struct {
	host    string
	timeout time.Duration
	logger  *zap.Logger
}

func New(cfgs ...ServerConfiguration) (*ServerService, error) {
	ss := &ServerService{}
	for _, cfg := range cfgs {
		err := cfg(ss)
		if err != nil {
			return nil, err
		}
	}
	return ss, nil
}

func WithHost(host string) ServerConfiguration {
	return func(ss *ServerService) error {
		ss.host = host
		return nil
	}
}

func WithTimetout(timeout time.Duration) ServerConfiguration {
	return func(ss *ServerService) error {
		ss.timeout = timeout
		return nil
	}
}

func WithLogger(logger *zap.Logger) ServerConfiguration {
	return func(ss *ServerService) error {
		ss.logger = logger
		return nil
	}
}

func (ss *ServerService) Send(body []byte, headers map[string][]string) error {
	ss.logger.Sugar().Infow("Sending to server", "destination", ss.host, "timeout", ss.timeout)
	reader := bytes.NewBuffer(body)
	req, err := http.NewRequest(http.MethodPost, ss.host, reader)
	if err != nil {
		// not using logger as they will be logged by the asyncCheckFeed function
		return err
	}
	req.Header = headers
	client := http.Client{
		Timeout: ss.timeout,
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		err := fmt.Errorf("Status Code not 200: %d", resp.StatusCode)
		return err
	}
	return nil
}
