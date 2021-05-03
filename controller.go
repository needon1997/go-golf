package golf

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"runtime"
	"strings"
)

//var R = &Controller{handlerMap: make(map[*regexp.Regexp]*handlerInfo)}

type Handler interface {
	ServeHTTP(ctx *Context)
}

type HandlerImpl func(ctx *Context)

func (fn HandlerImpl) ServeHTTP(ctx *Context) {
	fn(ctx)
	return
}

type Controller struct {
	prefix     string
	handlerMap map[*regexp.Regexp]*handlerInfo
}
type handlerInfo struct {
	params  map[int]string
	method  string
	handler Handler
}

func (this *Controller) Init(prefix string) {
	this.prefix = prefix
	this.handlerMap = make(map[*regexp.Regexp]*handlerInfo)
	//fmt.Println("init to the expanded by user if needed, such as adding handler")
}
func (this *Controller) BeforeServe(ctx *Context) {
	return
}
func (this *Controller) AfterServe(ctx *Context) {
	return
}
func (this *Controller) GetFuncMapping(pattern string, fn func(ctx *Context)) {
	this.AddFunc(pattern, "GET", fn)
}
func (this *Controller) PostFuncMapping(pattern string, fn func(ctx *Context)) {
	this.AddFunc(pattern, "POST", fn)
}
func (this *Controller) DeleteFuncMapping(pattern string, fn func(ctx *Context)) {
	this.AddFunc(pattern, "DELETE", fn)
}
func (this *Controller) PutFuncMapping(pattern string, fn func(ctx *Context)) {
	this.AddFunc(pattern, "PUT", fn)
}
func (this *Controller) AddFunc(pattern string, method string, fn func(ctx *Context)) {
	this.Add(pattern, method, HandlerImpl(fn))
}
func (this *Controller) GetMapping(pattern string, handler Handler) {
	this.Add(pattern, "GET", handler)
}
func (this *Controller) PostMapping(pattern string, handler Handler) {
	this.Add(pattern, "POST", handler)
}
func (this *Controller) DeleteMapping(pattern string, handler Handler) {
	this.Add(pattern, "DELETE", handler)
}
func (this *Controller) PutMapping(pattern string, handler Handler) {
	this.Add(pattern, "PUT", handler)
}
func (this *Controller) Add(pattern string, method string, handler Handler) {
	parts := strings.Split(pattern, "/")
	j := 0
	params := make(map[int]string)
	for i, part := range parts {
		if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") && len(part) > 2 {
			params[j] = part[1 : len(part)-1]
			parts[i] = "([A-Za-z0-9]+)"
			j++
		}
	}
	pattern = strings.Join(parts, "/")
	pattern += "$"
	regex, regexErr := regexp.Compile(pattern)
	if regexErr != nil {

		//TODO add error handling here to avoid panic
		panic(regexErr)
		return
	}
	for key, val := range this.handlerMap {
		if fmt.Sprintf("%v", key) == fmt.Sprintf("%v", regex) && val.method == method {
			panic("mapping conflict")
		}
	}
	// testPattern := "/122131"
	// fmt.Println(regex)
	// fmt.Println(regex.MatchString(testPattern))
	// fmt.Println(regex.FindStringSubmatch(testPattern))
	handlerInfo := &handlerInfo{}
	handlerInfo.handler = handler
	handlerInfo.method = method
	handlerInfo.params = params
	this.handlerMap[regex] = handlerInfo
	return
}

func (this *Controller) ServeHTTP(ctx *Context) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			for i := 1; ; i += 1 {
				_, file, line, ok := runtime.Caller(i)
				if !ok {
					break
				}
				fmt.Println(file, line)
			}
		}
	}()
	if ctx.End {
		return nil
	}
	requestPath := ctx.R.URL.Path
	//doublechek
	if !strings.HasPrefix(requestPath, this.prefix) {
		return errors.New("prefix error")
	}
	if this.prefix != "/" {
		requestPath = requestPath[len(this.prefix):]
	}
	if requestPath == "" {
		requestPath = "/"
	}
	found := false
	for key, val := range this.handlerMap {
		if !key.MatchString(requestPath) || ctx.R.Method != val.method {
			continue
		}
		matches := key.FindStringSubmatch(requestPath)
		if len(matches[0]) != len(requestPath) {
			continue
		}
		found = true
		params := make(map[string]string)
		if len(val.params) > 0 {
			values := ctx.R.URL.Query()
			for i, value := range matches[1:] {
				values.Add(val.params[i], value)
				params[val.params[i]] = value
			}
			ctx.R.URL.RawQuery = url.Values(values).Encode() + "&" + ctx.R.URL.RawQuery
		}
		val.handler.ServeHTTP(ctx)
		break
	}

	if found == false {
		return errors.New("no match url")
	}
	return nil
}
