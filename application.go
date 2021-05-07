package golf

import (
	"fmt"
	"github.com/needon1997/go-golf/session"
	"log"
	"net/http"
)

var App *Application

type DaemonGo func()

func init() {
	App = NewApplication("./config/app.conf")
}

type Application struct {
	router     *Router
	Config     *Config
	daemonList []DaemonGo
}

func NewApplication(path string) *Application {
	this := &Application{}
	c, err := LoadConfig(path)
	if err != nil {
		panic("config not found")
	}
	this.Config = c
	this.router = NewRouter()
	this.Initialize()
	App = this
	return this
}
func (this *Application) Initialize() {
	session_on, _ := this.Config.Bool("session.on")
	if session_on {
		session_type := this.Config.String("session.type")
		if session_type == "" {
			session_type = "memory"
		}
		session_name := this.Config.String("session.name")
		if session_name == "" {
			session_name = "X-Session-Id"
		}
		session_lifetime, err := this.Config.Int("session.lifetime")
		if err != nil {
			session_lifetime = 1800
		}
		session.GlobalSessions, err = session.NewManager(session_type, session_name, int64(session_lifetime))
		if err != nil {
			session.GlobalSessions = nil
		}
		go session.GlobalSessions.GC()
	}
}
func (this *Application) Run() {
	host := this.Config.String("host")
	http.Handle("/", this.router)
	for i := 0; i < len(this.daemonList); i++ {
		go this.daemonList[i]()
	}
	fmt.Println("server listen on " + host)
	err := http.ListenAndServe(host, nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
func (this *Application) Add(prefix string, c ControllerInterface) *Application {
	this.router.Add(prefix, c)
	return this
}

func (this *Application) RegisterDeamonGo(daemonGo DaemonGo) {
	this.daemonList = append(this.daemonList, daemonGo)
}

func (this *Application) AddStaticPath(prefix, path string) *Application {
	this.router.AddStaticPath(prefix, path)
	return this
}
