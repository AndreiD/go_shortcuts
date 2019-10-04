package gin

import (
	"github.com/go-delve/delve/pkg/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

func Init() {

	if conf.DebugMode {
		gin.SetMode(gin.DebugMode)
		log.SetLevel(log.DebugLevel)
	} else {
		gin.SetMode(gin.ReleaseMode)
		log.SetLevel(log.WarnLevel)
	}

	router = gin.New()
	router.Use(securitySetup())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(requestIDMiddleware())
	router.Use(corsMiddleware())
	router.Use(configurationMiddleware(conf))

	InitializeRouter()

	server := &http.Server{
		Addr:           conf.ServerHostname + ":" + strconv.Itoa(conf.ServerPort),
		Handler:        router,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1Mb
	}
	server.SetKeepAlivesEnabled(true)

	// Serve'em
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("initiated server shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("server shutdown:", err)
	}
	log.Println("server exiting. bye!")
}

func securitySetup() gin.HandlerFunc {

	secureMiddleware := secure.New(secure.Options{
		//AllowedHosts:          []string{"your_website\\.com", ".*\\.secured\\.com"},
		//HostsProxyHeaders:     []string{"X-Forwarded-Host"},
		//SSLRedirect:           true,
		//SSLHost:               "ssl.you.com",
		//SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
		//STSSeconds:            31536000,
		//STSIncludeSubdomains:  true,
		//STSPreload:            true,
		//FrameDeny:             true,
		//ContentTypeNosniff:    true,
		//BrowserXssFilter:      true,
		//ContentSecurityPolicy: "script-src $NONCE",
		//PublicKey:             `pin-sha256="base64+primary=="; pin-sha256="base64+backup=="; max-age=5184000; includeSubdomains; report-uri="https://www.xyz.com/hpkp-report"`,
		FrameDeny: true,
	})
	secureFunc := func() gin.HandlerFunc {
		return func(c *gin.Context) {
			err := secureMiddleware.Process(c.Writer, c.Request)
			if err != nil {
				c.Abort()
				return
			}
			if status := c.Writer.Status(); status > 300 && status < 399 {
				c.Abort()
			}
		}
	}()
	return secureFunc
}

// requestIDMiddleware adds x-request-id
func requestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Request-Id", uuid.NewV4().String())
		c.Next()
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers",
			"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// configurationMiddleware will add the configuration to the context
func configurationMiddleware(config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("configuration", config)
		c.Next()
	}
}
