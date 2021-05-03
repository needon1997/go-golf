package golf

import (
	"encoding/json"
	"fmt"
	"github.com/needon1997/go-golf/session"
	"io"
	"io/ioutil"
	"net/http"
)

type Context struct {
	W   http.ResponseWriter
	R   *http.Request
	End bool
}

func (this *Context) SessionStart() session.Session {
	if session.GlobalSessions != nil {
		return session.GlobalSessions.SessionStart(this.W, this.R)
	} else {
		fmt.Println("session service not been set to on, return nil")
		return nil
	}
}
func (this *Context) SessionDestroy() {
	if session.GlobalSessions != nil {
		session.GlobalSessions.SessionDestroy(this.W, this.R)
	} else {
		fmt.Println("session service not been set to on")
	}
}
func (this *Context) ReadBody() ([]byte, error) {
	bytes, err := ioutil.ReadAll(this.R.Body)
	defer this.R.Body.Close()
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
func (this *Context) ReadJson(i interface{}) error {
	bytes, err := this.ReadBody()
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, i)
	if err != nil {
		return err
	}
	return nil
}
func (this *Context) SetResponseJSON(statusCode int, res interface{}) {
	this.W.WriteHeader(statusCode)
	if res == nil {
		return
	}
	resStr, err := json.Marshal(res)
	if err != nil {
		this.SetResponseJSON(500, nil)
		return
	}
	io.WriteString(this.W, string(resStr))
}

//return the first valu
func (this *Context) GetFormValue(key string) (string, error) {
	err := this.R.ParseForm()
	if err != nil {
		fmt.Print(err)
		return "", err
	}
	return this.R.FormValue(key), nil
}

//return the value slice
func (this *Context) GetFormValues(key string) []string {
	return this.R.Form[key]
}
func (this *Context) Redirect(url string, code int) {
	http.Redirect(this.W, this.R, url, code)
}
