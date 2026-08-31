package backend

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHttpCommonGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ack := CommonAck[testData]{
			Code: CodeSuccess,
			Msg:  "ok",
			Data: testData{Name: "test", Value: 42},
		}
		json.NewEncoder(w).Encode(ack)
	}))
	defer server.Close()

	result, code, err := HttpCommonGet[testData](context.Background(), server.URL)
	if err != nil {
		t.Fatalf("HttpCommonGet error: %v", err)
	}
	if code != CodeSuccess {
		t.Fatalf("code: got %d, want %d", code, CodeSuccess)
	}
	if result.Name != "test" || result.Value != 42 {
		t.Fatalf("data: got %+v, want {test 42}", result)
	}
}

func TestHttpCommonGetError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ack := CommonAck[testData]{
			Code: 1001,
			Msg:  "not found",
		}
		json.NewEncoder(w).Encode(ack)
	}))
	defer server.Close()

	_, code, err := HttpCommonGet[testData](context.Background(), server.URL)
	if err == nil {
		t.Fatal("expected error")
	}
	if code != 1001 {
		t.Fatalf("code: got %d, want 1001", code)
	}
}

func TestHttpCommonPost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ack := CommonAck[testData]{
			Code: CodeSuccess,
			Msg:  "ok",
			Data: testData{Name: "posted", Value: 100},
		}
		json.NewEncoder(w).Encode(ack)
	}))
	defer server.Close()

	body, _ := json.Marshal(map[string]string{"key": "value"})
	result, code, err := HttpCommonPost[testData](context.Background(), server.URL, body)
	if err != nil {
		t.Fatalf("HttpCommonPost error: %v", err)
	}
	if code != CodeSuccess {
		t.Fatalf("code: got %d, want %d", code, CodeSuccess)
	}
	if result.Name != "posted" || result.Value != 100 {
		t.Fatalf("data: got %+v, want {posted 100}", result)
	}
}

type testData struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}
