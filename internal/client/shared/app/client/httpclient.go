package client

import (
	"net/http"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

func NewHTTPClient(baseURL string, logger *zap.Logger) *resty.Client {
	httpClient := &http.Client{
		Transport: CustomRoundTripper{proxy: http.DefaultTransport, logger: logger},
	}
	client := resty.NewWithClient(httpClient)
	client.BaseURL = baseURL
	return client
}

type CustomRoundTripper struct {
	proxy  http.RoundTripper
	logger *zap.Logger
}

func (t CustomRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	t.logger.Info("request",
		zap.String("uri", request.URL.String()),
		zap.String("method", request.Method),
	)
	response, err := t.proxy.RoundTrip(request)
	if err != nil {
		t.logger.Error("request error", zap.String("error", err.Error()))
		return nil, err
	}

	t.logger.Info("command_response",
		zap.String("uri", request.URL.String()),
		zap.String("method", request.Method),
		zap.String("status", response.Status),
	)

	return response, nil
}
