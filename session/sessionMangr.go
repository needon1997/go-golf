package session

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var GlobalSessions *Manager
var provides = make(map[string]Provider)

// Register makes a session provide available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, provider Provider) {
	if provider == nil {
		panic("session: Register provider is nil")
	}
	if _, dup := provides[name]; dup {
		panic("session: Register called twice for provider " + name)
	}
	provides[name] = provider
}

type Manager struct {
	cookieName  string     // private cookiename
	lock        sync.Mutex // protects session
	provider    Provider
	maxLifeTime int64
}

func NewManager(provideName, cookieName string, maxLifeTime int64) (*Manager, error) {
	provider, ok := provides[provideName]
	if !ok {
		return nil, fmt.Errorf("session: unknown provide %q (forgotten import?)", provideName)
	}
	man := &Manager{provider: provider, cookieName: cookieName, maxLifeTime: maxLifeTime}
	return man, nil
}

//generate a session Id
func (manager *Manager) sessionId() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

//get the session if exist or create one if it does not exist
func (manager *Manager) SessionStart(w http.ResponseWriter, r *http.Request) (session Session) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" { //if it does not exis
		sid := manager.sessionId()                                                                                                                //generate a unique session id
		session, _ = manager.provider.SessionInit(sid)                                                                                            // init the session in the provider
		cookie := http.Cookie{Name: manager.cookieName, Value: url.QueryEscape(sid), Path: "/", HttpOnly: true, MaxAge: int(manager.maxLifeTime)} //create one session
		http.SetCookie(w, &cookie)                                                                                                                //set the session in the cookie
	} else { // if it exist
		sid, _ := url.QueryUnescape(cookie.Value)      //get the sessionId from cookie
		session, _ = manager.provider.SessionRead(sid) //get the corresponde session
	}
	return
}

//Destroy sessionid
func (manager *Manager) SessionDestroy(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(manager.cookieName) //get the cookie that hold the sessionId
	if err != nil || cookie.Value == "" {
		return
	} else {
		manager.lock.Lock()
		defer manager.lock.Unlock()
		manager.provider.SessionDestroy(cookie.Value) //destroy the session base on sessionI
		expiration := time.Now()
		cookie := http.Cookie{Name: manager.cookieName, Path: "/", HttpOnly: true, Expires: expiration, MaxAge: -1} //set the expiration to now and inform the client to delete the cookie
		http.SetCookie(w, &cookie)
	}
}
func (manager *Manager) GC() {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	manager.provider.SessionGC(manager.maxLifeTime)
	time.AfterFunc(time.Duration(manager.maxLifeTime), func() { manager.GC() }) //this function will run into its routine, so it won't hold the lock forever
}
