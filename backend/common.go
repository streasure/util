package backend

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	backendHttpTimeout = 5 * time.Second
	CodeSuccess        = 0
	CodeInternal       = 1
)

func httpGetReadBody(ctx context.Context, url string) ([]byte, error) {
	httpClient := http.Client{
		Timeout: backendHttpTimeout,
	}

	httpResp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = httpResp.Body.Close()
	}()

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status code: %d", httpResp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	return bodyBytes, nil
}

func httpPostReadBody(ctx context.Context, url string, body io.Reader) ([]byte, error) {
	httpClient := http.Client{
		Timeout: backendHttpTimeout,
	}

	httpResp, err := httpClient.Post(url, "application/json", body)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = httpResp.Body.Close()
	}()

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status code: %d", httpResp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	return bodyBytes, nil
}

func HttpCommonGet[T any](ctx context.Context, url string) (*T, int, error) {
	bodyBytes, err := httpGetReadBody(ctx, url)
	if err != nil {
		return nil, CodeInternal, err
	}

	ack := &CommonAck[T]{}
	err = json.Unmarshal(bodyBytes, &ack)
	if err != nil {
		return nil, CodeInternal, err
	}

	if ack.Code != CodeSuccess {
		return nil, ack.Code, fmt.Errorf(ack.Msg)
	}

	return &ack.Data, ack.Code, nil
}

func HttpCommonPost[T any](ctx context.Context, url string, paramBytes []byte) (*T, int, error) {
	reader := bytes.NewReader(paramBytes)

	ack := &CommonAck[T]{}
	bodyBytes, err := httpPostReadBody(ctx, url, reader)
	if err != nil {
		return nil, CodeInternal, err
	}

	err = json.Unmarshal(bodyBytes, &ack)
	if err != nil {
		return nil, CodeInternal, err
	}

	if ack.Code != CodeSuccess {
		return nil, ack.Code, fmt.Errorf(ack.Msg)
	}

	return &ack.Data, CodeSuccess, nil
}

func ActivationCode(ctx context.Context, url string) (*ActivationCodeAck, error) {
	bodyBytes, err := httpGetReadBody(ctx, url)
	if err != nil {
		return nil, err
	}

	ack := ActivationCodeAck{}
	err = json.Unmarshal(bodyBytes, &ack)
	if err != nil {
		return nil, err
	}

	return &ack, nil
}
