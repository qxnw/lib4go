package http

import (
	"net/http"
	"testing"
)

// func TestNewHTTPClientCert(t *testing.T) {
// 	certFile := "/home/champly/http.cer"
// 	keyFile := "/home/champly/http.key"
// 	caFile := "/home/champly/http.ca"
// 	_, err := NewHTTPClientCert(certFile, keyFile, caFile)
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

func TestGet(t *testing.T) {
	client := NewHTTPClient()
	content, status, err := client.Get("http://www.baidu.com", nil...)
	if err != nil {
		t.Errorf("Get error: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Get error status:%d", status)
	}
	if content == "" {
		t.Error("Get error with not content")
	}

	// content, status, err = client.Get("http://www.google.com", nil...)
	// if err != nil {
	// 	t.Errorf("Get error: %v", err)
	// }
	// if status != http.StatusRequestTimeout {
	// 	t.Errorf("Get error status:%d", status)
	// }
	// if content != "" {
	// 	t.Error("Get error with not content")
	// }

	content, status, err = client.Get("http://192.168.0.121:8013", nil...)
	if err != nil {
		t.Errorf("Get error: %v", err)
	}
	if status != http.StatusNotFound {
		t.Errorf("Get error status:%d", status)
	}
	if content == "" {
		t.Error("Get error with not content")
	}
}

func TestPost(t *testing.T) {
	client := NewHTTPClient()
	content, status, err := client.Post("http://www.baidu.com", "name=bob", "UTF-8")
	if err != nil {
		t.Errorf("Get error: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("Get error status:%d", status)
	}
	if content == "" {
		t.Error("Get error with not content")
	}

	content, status, err = client.Post("http://192.168.0.121:8013", "name=bob", "UTF-8")
	if err != nil {
		t.Errorf("Get error: %v", err)
	}
	if status != http.StatusNotFound {
		t.Errorf("Get error status:%d", status)
	}
	if content == "" {
		t.Error("Get error with not content")
	}

	// content, status, err = client.Post("http://www.google.com", "name=bob", "UTF-8")
	// if err != nil {
	// 	t.Errorf("Get error: %v", err)
	// }
	// if status != http.StatusRequestTimeout {
	// 	t.Errorf("Get error status:%d", status)
	// }
	// if content != "" {
	// 	t.Error("Get error with not content")
	// }
}

func TestRequest(t *testing.T) {
	method := "api/query_balance_by_platform_tag"
}
