package xgo

import (
	"os"
)

var (
	app               *xgoApp
	apps              map[int]*xgoApp
	appIdGen          *AutoIncr
	util              xgoUtil
	cfg               *xgoConfig
	cfgFile           string
	ListenAddr        string
	ListenPort        int
	RunMode           string
	EnableStats       bool
	CookieSecret      string
	SessionName       string
	SessionTTL        int64
	EnablePprof       bool
	EnableGzip        bool
	GzipMinLength     int
	GzipTypes         []string
	SslCertificate    string
	SslCertificateKey string
)

func init() {
	defaultCfg := &xgoDefaultConfig{
		ListenAddr:        "",
		ListenPort:        80,
		RunMode:           "http",
		EnableStats:       true,
		CookieSecret:      "foobar",
		SessionName:       "XGOSESSID",
		SessionTTL:        60 * 15,
		EnablePprof:       true,
		EnableGzip:        true,
		GzipMinLength:     1024,
		GzipTypes:         []string{"text", "javascript", "css", "xml"},
		SslCertificate:    "",
		SslCertificateKey: "",
	}

	cfgFile = "app.conf"
	// Check the first argument of cmd line,
	// if it is not a flag (begin with '-'),
	// try to use it as the config file path.
	if len(os.Args) > 1 {
		arg := os.Args[1]
		if arg[0] != '-' {
			cfgFile = arg
		}
	}

	cfg = &xgoConfig{}
	cfg.LoadFile(cfgFile)
	cfg.RegisterConfig(defaultCfg)
	apps = make(map[int]*xgoApp)
	appIdGen = NewAutoIncr(1, 1)
	app = NewApp()
	util = xgoUtil{}
}

func NewApp() *xgoApp {
	id := appIdGen.Fetch()
	a := new(xgoApp).init(id)
	apps[id] = a
	return a
}

func GetMainApp() *xgoApp {
	return app
}

func RegisterHandler(pattern string, c xgoHandlerInterface) {
	app.RegisterHandler(pattern, c)
}

func RegisterHandlerHook(event string, hookFunc HookHandlerFunc) {
	app.RegisterHandlerHook(event, hookFunc)
}

func SetStaticPath(sPath, fPath string) {
	app.SetStaticPath(sPath, fPath)
}

func SetStaticFileType(ext string) {
	app.SetStaticFileType(ext)
}

func RegisterSessionStorage(storage SessionStorageInterface) {
	app.RegisterSessionStorage(storage)
}

func RegisterCustomHttpStatus(code int, filePath string) {
	app.RegisterCustomHttpStatus(code, filePath)
}

func Run() error {
	return app.Run(RunMode, ListenAddr, ListenPort)
}

func LoadConfig(conf interface{}) error {
	return cfg.RegisterConfig(conf)
}

func ReloadConfig() error {
	return cfg.ReloadFile()
}
