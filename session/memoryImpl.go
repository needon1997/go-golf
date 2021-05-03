package session

import (
	"container/list"
	"sync"
	"time"
)

func init() {
	var provider = &MemoryProvider{list: list.New()}
	provider.sessions = make(map[string]*list.Element)
	Register("memory", provider)
}

type MemorySession struct {
	sid            string
	lastAccessTime time.Time
	value          map[interface{}]interface{}
	provider       *MemoryProvider
}

func (this *MemorySession) Set(key, value interface{}) error {
	this.value[key] = value
	this.lastAccessTime = time.Now()
	this.provider.SessionUpdateAcess(this.sid)
	return nil
}
func (this *MemorySession) Get(key interface{}) interface{} {
	this.lastAccessTime = time.Now()
	this.provider.SessionUpdateAcess(this.sid)
	if val, ok := this.value[key]; ok {
		return val
	} else {
		return nil
	}
}

func (this *MemorySession) Delete(key interface{}) error {
	this.lastAccessTime = time.Now()
	this.provider.SessionUpdateAcess(this.sid)
	delete(this.value, key)
	return nil
}

func (this *MemorySession) SessionID() string {
	return this.sid
}

type MemoryProvider struct {
	lock     sync.Mutex               //用来锁
	sessions map[string]*list.Element //用来存储在内存
	list     *list.List               //用来做gc
}

func (this *MemoryProvider) SessionInit(sid string) (Session, error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	session := &MemorySession{sid: sid, lastAccessTime: time.Now(), value: make(map[interface{}]interface{}), provider: this}
	ele := this.list.PushBack(session)
	this.sessions[sid] = ele
	return session, nil

}
func (this *MemoryProvider) SessionRead(sid string) (Session, error) {
	this.lock.Lock()
	ele, ok := this.sessions[sid]
	this.lock.Unlock()
	if !ok {
		return this.SessionInit(sid)
	} else {
		return ele.Value.(*MemorySession), nil
	}
}
func (this *MemoryProvider) SessionDestroy(sid string) error {
	this.lock.Lock()
	defer this.lock.Unlock()
	if ele, ok := this.sessions[sid]; ok {
		delete(this.sessions, sid)
		this.list.Remove(ele)
		return nil
	}
	return nil
}
func (this *MemoryProvider) SessionGC(maxLifeTime int64) {
	this.lock.Lock()
	defer this.lock.Unlock()
	for {
		element := this.list.Back()
		if element == nil {
			break
		}
		if (element.Value.(*MemorySession).lastAccessTime.Unix() + maxLifeTime) < time.Now().Unix() {
			this.list.Remove(element)
			delete(this.sessions, element.Value.(*MemorySession).sid)
		} else {
			break
		}
	}
}
func (this *MemoryProvider) SessionUpdateAcess(sid string) error {
	this.lock.Lock()
	defer this.lock.Unlock()
	if ele, ok := this.sessions[sid]; ok {
		this.list.MoveToFront(ele)
	}
	return nil
}
