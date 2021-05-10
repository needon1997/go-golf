package golf_test

import (
	"encoding/json"
	"github.com/needon1997/go-golf/httplib"
	"net/http"
	"testing"
)

//func TestDaemon(t *testing.T) {
//	golf.App.RegisterDaemonGo(func() {
//		for i := 0; i < 10; i++ {
//			fmt.Println(i)
//		}
//	})
//	golf.App.Run()
//}

func TestGet(t *testing.T) {
	statusCode, _, err := httplib.Get("http://localhost:10081/echo")
	if err != nil {
		t.Error(err)
	}
	if statusCode != http.StatusOK {
		t.Error("status not ok")
	}
}

func TestPost(t *testing.T) {
	reqBodyByte, err := json.Marshal(map[string]string{"msg": "hello"})
	statusCode, rspBodyByte, err := httplib.Post("http://localhost:10081/echo", reqBodyByte)
	if err != nil {
		t.Error(err)
	}
	if statusCode != http.StatusOK {
		t.Error("status not ok")
	}
	rspBody := make(map[string]string)
	json.Unmarshal(rspBodyByte, &rspBody)
	if msg, ok := rspBody["msg"]; !ok {
		t.Error("wrong response")
	} else {
		if msg != "hello" {
			t.Error("wrong response")
		}
	}
}

func TestPut(t *testing.T) {
	reqBodyByte, err := json.Marshal(map[string]string{"msg": "hello"})
	statusCode, rspBodyByte, err := httplib.Put("http://localhost:10081/echo", reqBodyByte)
	if err != nil {
		t.Error(err)
	}
	if statusCode != http.StatusOK {
		t.Error("status not ok")
	}
	rspBody := make(map[string]string)
	json.Unmarshal(rspBodyByte, &rspBody)
	if msg, ok := rspBody["msg"]; !ok {
		t.Error("wrong response")
	} else {
		if msg != "hello" {
			t.Error("wrong response")
		}
	}
}

func TestDelete(t *testing.T) {
	reqBodyByte, err := json.Marshal(map[string]string{"msg": "hello"})
	statusCode, rspBodyByte, err := httplib.Delete("http://localhost:10081/echo", reqBodyByte)
	if err != nil {
		t.Error(err)
	}
	if statusCode != http.StatusOK {
		t.Error("status not ok")
	}
	rspBody := make(map[string]string)
	json.Unmarshal(rspBodyByte, &rspBody)
	if msg, ok := rspBody["msg"]; !ok {
		t.Error("wrong response")
	} else {
		if msg != "hello" {
			t.Error("wrong response")
		}
	}
}
