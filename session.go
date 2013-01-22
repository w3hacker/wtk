package xgo

import (
	"fmt"
	"time"
)

type xgoSessionStorageInterface interface {
	Init(int64)
	GC()
	CreateSessionID() string
	Set(string, map[string]string)
	Get(string) map[string]string
	Delete(string)
}

type xgoSessionManager struct {
	sessionStorage xgoSessionStorageInterface
}

func (this *xgoSessionManager) RegisterStorage(storage xgoSessionStorageInterface) {
	if storage == nil {
		return
	}
	this.sessionStorage = storage
	storage.Init(SessionTTL)
}

func (this *xgoSessionManager) CreateSessionID() string {
	return this.sessionStorage.CreateSessionID()
}

func (this *xgoSessionManager) Set(sid string, data map[string]string) {
	this.sessionStorage.Set(sid, data)
}

func (this *xgoSessionManager) Get(sid string) map[string]string {
	return this.sessionStorage.Get(sid)
}

func (this *xgoSessionManager) Delete(sid string) {
	this.sessionStorage.Delete(sid)
}

type xgoSession struct {
	sessionManager *xgoSessionManager
	sessionId      string
	ctx            *xgoContext
	data           map[string]string
}

func (this *xgoSession) init() {
	if this.sessionId == "" {
		this.sessionId = this.sessionManager.CreateSessionID()
		this.ctx.SetCookie(SessionName, this.sessionId, 0)
	}
	if this.data == nil {
		this.data = this.sessionManager.Get(this.sessionId)
	}
}

func (this *xgoSession) Get(key string) string {
	this.init()
	if data, exist := this.data[key]; exist {
		return data
	}
	return ""
}

func (this *xgoSession) Set(key string, data string) {
	this.init()
	this.data[key] = data
	this.sessionManager.Set(this.sessionId, this.data)
}

func (this *xgoSession) Delete(key string) {
	this.init()
	delete(this.data, key)
	this.sessionManager.Set(this.sessionId, this.data)
}

type xgoDefaultSessionStorage struct {
	ttl   int64
	datas map[string]xgoDefaultSessionStorageData
}

type xgoDefaultSessionStorageData struct {
	expires int64
	data    map[string]string
}

func (this *xgoDefaultSessionStorage) Init(ttl int64) {
	if this.datas != nil {
		return
	}
	this.ttl = ttl
	this.datas = make(map[string]xgoDefaultSessionStorageData)
}

func (this *xgoDefaultSessionStorage) GC() {
	for {
		if len(this.datas) > 0 {
			now := time.Now().Unix()
			for sid, data := range this.datas {
				if data.expires <= now {
					delete(this.datas, sid)
				}
			}
		}
		time.Sleep(time.Second)
	}
}

func (this *xgoDefaultSessionStorage) CreateSessionID() string {
	t := time.Now()
	return "SESS" + fmt.Sprintf("%d%d", t.Unix(), t.Nanosecond())
}

func (this *xgoDefaultSessionStorage) Set(sid string, data map[string]string) {
	d := xgoDefaultSessionStorageData{
		expires: time.Now().Unix() + this.ttl,
		data:    data,
	}
	this.datas[sid] = d
}

func (this *xgoDefaultSessionStorage) Get(sid string) map[string]string {
	if data, exist := this.datas[sid]; exist {
		data.expires = time.Now().Unix() + this.ttl
		return data.data
	}
	return make(map[string]string)
}

func (this *xgoDefaultSessionStorage) Delete(sid string) {
	delete(this.datas, sid)
}
