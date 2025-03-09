package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/oddmario/systemstats-agent/httpserver"
	"github.com/oddmario/systemstats-agent/utils"
	"github.com/oddmario/systemstats-agent/workers"
	_ "go.uber.org/automaxprocs"
)

func main() {
	debug.SetMaxStack(2 * 1024 * 1024 * 1024)
	debug.SetMaxThreads(100000)

	utils.LoadConfig(true)

	workers.Init()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	if utils.Config.Get("enable_pprof").Bool() {
		pprof.Register(r, "debug/"+utils.Config.Get("pprof_secret").String()+"/pprof")
	}
	if utils.Config.Get("corsAllowAll").Bool() {
		r.Use(cors.Default())
	}

	r.Use(func() gin.HandlerFunc {
		return func(c *gin.Context) {
			c.Writer.Header().Set("Server", "systemstats-agent")
		}
	}())

	r.UseRawPath = true
	r.UnescapePathValues = false

	initErrors(r)
	initRoutes(r)

	if utils.Config.Get("ssl.enabled").Bool() {
		SSL_KEY_PATH := utils.Config.Get("ssl.key_path").String()
		SSL_CERT_PATH := utils.Config.Get("ssl.cert_path").String()

		go func() {
			httpsListenerData := utils.Config.Get("ssl.listener").String()
			fmt.Println("Listening at " + httpsListenerData + "...")
			httpsLn, httpsLnErr := net.Listen("tcp", httpsListenerData)
			if httpsLnErr != nil {
				fmt.Println("Error starting listener:", httpsLnErr)

				return
			}

			httpserver.HttpsServer = &http.Server{
				Addr:              httpsListenerData,
				Handler:           r,
				ReadTimeout:       0,
				WriteTimeout:      0,
				IdleTimeout:       0,
				ReadHeaderTimeout: 1 * time.Minute,                                              // https://ieftimov.com/posts/make-resilient-golang-net-http-servers-using-timeouts-deadlines-context-cancellation/#server-timeouts---first-principles
				TLSNextProto:      make(map[string]func(*http.Server, *tls.Conn, http.Handler)), // disable http2
			}
			httpserver.HttpsServer.SetKeepAlivesEnabled(true)

			tcpListenerHTTPS := &customListener{httpsLn.(*net.TCPListener)}
			httpsServeErr := httpserver.HttpsServer.ServeTLS(tcpListenerHTTPS, SSL_CERT_PATH, SSL_KEY_PATH)
			if httpsServeErr != nil {
				fmt.Println(httpsServeErr)

				return
			}
		}()
	}

	listenerData := utils.Config.Get("listener").String()
	fmt.Println("Listening at " + listenerData + "...")
	ln, err := net.Listen("tcp", listenerData)
	if err != nil {
		fmt.Println("Error starting listener:", err)

		return
	}

	// disable timeouts to prevent interruptions during large file uploads & downloads.
	// ... see https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/ (the "About streaming" part)
	// ... and see https://ieftimov.com/posts/make-resilient-golang-net-http-servers-using-timeouts-deadlines-context-cancellation/
	httpserver.HttpServer = &http.Server{
		Addr:              listenerData,
		Handler:           r,
		ReadTimeout:       0,
		WriteTimeout:      0,
		IdleTimeout:       0,
		ReadHeaderTimeout: 1 * time.Minute,                                              // https://ieftimov.com/posts/make-resilient-golang-net-http-servers-using-timeouts-deadlines-context-cancellation/#server-timeouts---first-principles
		TLSNextProto:      make(map[string]func(*http.Server, *tls.Conn, http.Handler)), // disable http2
	}
	httpserver.HttpServer.SetKeepAlivesEnabled(true) // even though we hate using keepalives on our HTTP clients, but always make sure keep alives are supported and enabled on our HTTP server!!!

	tcpListener := &customListener{ln.(*net.TCPListener)}
	err = httpserver.HttpServer.Serve(tcpListener)
	if err != nil {
		fmt.Println(err)

		return
	}
}
