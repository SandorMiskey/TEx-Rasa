// region: packages

package api

import (
	"fmt"
	"log/syslog"
	"net"
	"os"
	"strconv"
	"sync"

	"github.com/SandorMiskey/TEx-kit/cfg"
	// "github.com/SandorMiskey/TEx-kit/db"
	"github.com/SandorMiskey/TEx-kit/log"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

// endregion: packages
// region: global variables

var (
	Config cfg.Config
	// Db     *db.Db
	Logger log.Logger
)

const (
	LOG_ERR    syslog.Priority = log.LOG_ERR
	LOG_NOTICE syslog.Priority = log.LOG_NOTICE
	LOG_INFO   syslog.Priority = log.LOG_INFO
	LOG_DEBUG  syslog.Priority = log.LOG_DEBUG
	LOG_EMERG  syslog.Priority = log.LOG_EMERG
)

// endregion: globals

func obsolete() {

	// region: config and cli flags

	Config = *cfg.NewConfig(os.Args[0])
	flagSet := Config.NewFlagSet(os.Args[0])
	flagSet.Entries = map[string]cfg.Entry{
		// "dbAddr":        {Desc: "database address", Type: "string", Def: "/app/mgmt.db"},
		// "dbName":        {Desc: "database name", Type: "string", Def: "mgmt"},
		// "dbPasswd":      {Desc: "database user password", Type: "string", Def: ""},
		// "dbPasswd_file": {Desc: "database user password", Type: "string", Def: ""},
		// "dbType":        {Desc: "db type as in TEx-kit/db/db.go", Type: "int", Def: 4},
		// "dbUser":        {Desc: "database user", Type: "string", Def: "mgmt"},

		"httpEnabled":            {Desc: "enable http", Type: "bool", Def: true},
		"httpName":               {Desc: "server name in response header", Type: "string", Def: "TEx-Rasa management service"},
		"httpLogAllErrors":       {Desc: "enable http", Type: "bool", Def: true},
		"httpMaxRequestBodySize": {Desc: "http max request body size ", Type: "int", Def: 4 * 1024 * 1024},
		"httpNetworkProto":       {Desc: "network protocol must be 'tcp', 'tcp4', 'tcp6', 'unix' or 'unixpacket'", Type: "string", Def: "tcp"},
		"httpPort":               {Desc: "http port", Type: "int", Def: 5000},
		"httpStaticEnabled":      {Desc: "enable serving static files", Type: "bool", Def: true},
		"httpStaticRoot":         {Desc: "path to static files", Type: "string", Def: "/app/public"},
		"httpStaticIndex":        {Desc: "index file to serve during directory access", Type: "string", Def: "index.html"},
		"httpStaticError":        {Desc: "location to redirect in case of 404", Type: "string", Def: "index.html"},
		"httpTLSCert":            {Desc: "https certificate", Type: "string", Def: ""},
		"httpTLSCert_file":       {Desc: "https certificate file", Type: "string", Def: ""},
		"httpTLSEnabled":         {Desc: "enable https", Type: "bool", Def: false},
		"httpTLSKey":             {Desc: "private key for HTTPS certificate", Type: "string", Def: ""},
		"httpTLSKey_file":        {Desc: "httpTLSKey file", Type: "string", Def: ""},
		"httpTLSPort":            {Desc: "https port", Type: "int", Def: 5000},

		"logLevel": {Desc: "Logger min severity", Type: "int", Def: 7},
	}

	err := flagSet.ParseCopy()
	if err != nil {
		panic(err)
	}

	// endregion: cli flags
	// region: logger

	logLevel := syslog.Priority(Config.Entries["logLevel"].Value.(int))

	Logger = *log.NewLogger()
	defer Logger.Close()
	_, _ = Logger.NewCh(log.ChConfig{Severity: &logLevel})

	// Logger.Out(LOG_DEBUG, spew.Sdump(Config))

	// endregion: logger
	// region: db

	// if db.DbType(Config.Entries["dbType"].Value.(int)) == db.Postgres {
	// 	db.Defaults = db.DefaultsPostgres // TODO: this goes to TEx-kit/db/db.go
	// }

	// dbConfig := db.Config{
	// 	Addr:   Config.Entries["dbAddr"].Value.(string),
	// 	DBName: Config.Entries["dbName"].Value.(string),
	// 	Logger: Logger,
	// 	// Params: nil,
	// 	Passwd: Config.Entries["dbPasswd"].Value.(string),
	// 	Type:   db.DbType(Config.Entries["dbType"].Value.(int)),
	// 	User:   Config.Entries["dbUser"].Value.(string),
	// }

	// Db, err = dbConfig.Open()
	// defer Db.Close()

	// if err != nil {
	// 	Logger.Out(LOG_EMERG, err)
	// 	panic(err)
	// }

	// endregion: db
	// region: http routing

	httpRouterActual := fasthttprouter.New()
	if Config.Entries["httpStaticEnabled"].Value.(bool) {
		httpFS := &fasthttp.FS{
			Root:       Config.Entries["httpStaticRoot"].Value.(string),
			IndexNames: []string{Config.Entries["httpStaticIndex"].Value.(string)},
			PathNotFound: func(ctx *fasthttp.RequestCtx) {
				Logger.Out(LOG_NOTICE, "dead end", ctx)
				ctx.Redirect(Config.Entries["httpStaticError"].Value.(string), 303)
			},
			Compress:           true,
			AcceptByteRange:    true,
			GenerateIndexPages: false,
		}

		httpRouterActual.NotFound = httpFS.NewRequestHandler()
	}

	httpRouterPre := fasthttprouter.New()
	httpRouterPre.NotFound = func(ctx *fasthttp.RequestCtx) {
		Logger.Out(LOG_DEBUG, fmt.Sprintf("%s request on %s from %s with content type '%s' and body '%s' (%s)", ctx.Method(), ctx.Path(), ctx.RemoteAddr(), ctx.Request.Header.Peek("Content-Type"), ctx.PostBody(), ctx))
		Logger.Out(LOG_INFO, ctx)
		httpRouterActual.Handler(ctx)
	}

	// endregion: http routing
	// region: http and https

	var wg sync.WaitGroup

	if Config.Entries["httpEnabled"].Value.(bool) {
		http := &fasthttp.Server{
			// Logger:             Logger,
			Handler:            httpRouterPre.Handler,
			LogAllErrors:       Config.Entries["httpLogAllErrors"].Value.(bool),
			MaxRequestBodySize: Config.Entries["httpMaxRequestBodySize"].Value.(int),
			Name:               Config.Entries["httpName"].Value.(string),
		}
		ln, err := net.Listen(Config.Entries["httpNetworkProto"].Value.(string), ":"+strconv.Itoa(Config.Entries["httpPort"].Value.(int)))
		if err != nil {
			Logger.Out(LOG_ERR, "error while opening http listener", err)
		} else {

			wg.Add(1)
			go func() {
				Logger.Out(LOG_INFO, "listening for HTTP requests", Config.Entries["httpNetworkProto"].Value, Config.Entries["httpPort"].Value)
				http.Serve(ln)
			}()
		}
	}

	if Config.Entries["httpTLSEnabled"].Value.(bool) {
		https := &fasthttp.Server{
			// Logger:          Logger,
			Handler:            httpRouterPre.Handler,
			LogAllErrors:       Config.Entries["httpLogAllErrors"].Value.(bool),
			MaxRequestBodySize: Config.Entries["httpMaxRequestBodySize"].Value.(int),
			Name:               Config.Entries["httpName"].Value.(string),
		}
		ln, err := net.Listen(Config.Entries["httpNetworkProto"].Value.(string), ":"+strconv.Itoa(Config.Entries["httpTLSPort"].Value.(int)))
		if err != nil {
			Logger.Out(LOG_ERR, "error while opening https listener", err)
		} else {

			wg.Add(1)
			go func() {
				Logger.Out(LOG_INFO, "listening for HTTPS requests", Config.Entries["httpNetworkProto"].Value, Config.Entries["httpTLSPort"].Value)
				https.ServeTLSEmbed(ln, []byte(Config.Entries["httpTLSCert"].Value.(string)), []byte(Config.Entries["httpTLSKey"].Value.(string)))
			}()
		}
	}

	wg.Wait()

	// endregion: http and https

}
