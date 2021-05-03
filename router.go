package golf

import (
	"net/http"
	"runtime"
	"strings"
)

type Router struct {
	staticDir     map[string]string
	controllerMap map[string]ControllerInterface
}

type ControllerInterface interface {
	Init(prefix string)
	BeforeServe(ctx *Context)
	ServeHTTP(ctx *Context) error
	AfterServe(ctx *Context)
}

func NewRouter() *Router {
	this := &Router{}
	this.controllerMap = make(map[string]ControllerInterface)
	this.staticDir = make(map[string]string)
	return this
}
func (this *Router) Add(prefix string, c ControllerInterface) *Router {
	_, ok := this.controllerMap[prefix]
	if !ok {
		c.Init(prefix)
		this.controllerMap[prefix] = c
	} else {
		panic("duplicate mapping")
	}
	return this
}
func (this *Router) AddStaticPath(prefix, path string) *Router {
	this.staticDir[prefix] = path
	return this
}

func (this *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			Critical(err)
			for i := 1; ; i += 1 {
				_, file, line, ok := runtime.Caller(i)
				if !ok {
					break
				}
				Critical(file, line)
			}
		}
	}()
	found := false
	requestPath := r.URL.Path
	for prefix, path := range this.staticDir {
		if strings.HasPrefix(requestPath, prefix) {
			filepath := path + requestPath[len(prefix):]
			http.ServeFile(w, r, filepath)
			found = true
			return
		}
	}

	for prefix, controller := range this.controllerMap {
		if strings.HasPrefix(requestPath, prefix) && (len(requestPath) == len(prefix) || strings.HasPrefix(requestPath[len(prefix):], "/")) {
			found = true
			ctx := &(Context{W: w, R: r})
			controller.BeforeServe(ctx)
			err := controller.ServeHTTP(ctx)
			if err != nil {
				http.Error(w, err.Error(), 400)
			}
			controller.AfterServe(ctx)
			return
		}
	}
	controller, ok := this.controllerMap["/"]
	if ok {
		found = true
		ctx := &(Context{W: w, R: r})
		controller.BeforeServe(ctx)
		err := controller.ServeHTTP(ctx)
		if err != nil {
			http.Error(w, err.Error(), 400)
		}
		controller.AfterServe(ctx)
		return
	}
	if found == false {
		http.NotFound(w, r)
	}

}
