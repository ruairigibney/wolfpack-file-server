package main

import (
	"log"
	"os"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/kelseyhightower/envconfig"
	"github.com/patrickmn/go-cache"
	wpfsConfig "github.com/ruairigibney/wolfpack-file-server/src/config"
	api "github.com/ruairigibney/wolfpack-file-server/src/http"
)

type WpfsCfg struct {
	wpfsCfg *wpfsConfig.Config
}

func main() {
	c := cache.New(24*time.Hour, 10*time.Minute)
	cookieKey := securecookie.GenerateRandomKey(16)
	store := sessions.NewCookieStore(cookieKey)
	store.Options = &sessions.Options{
		Path:   "/",
		MaxAge: int(time.Hour * 24),
	}

	var cfg = wpfsConfig.Config{}
	err := envconfig.Process("FS", &cfg)
	if err != nil {
		log.Fatal("Error reading config")
	}
	cfg.Store = store
	cfg.C = c

	checkConfig(&cfg)

	logFile, err := os.OpenFile("wolfpack-file-server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logFile)

	api.HttpHandler(&cfg)
}

func checkConfig(cfg *wpfsConfig.Config) {
	if cfg.LogDirectory == "" || cfg.ArchiveDirectory == "" || cfg.Host == "" ||
		cfg.PasscodePort == "" || cfg.ArchivePort == "" || cfg.Secret == "" {
		log.Printf(`logDirectory: %s; archiveDirectory: %s; host: %s; passCodePort: %s
		archivePort: %s; secret: %s`, cfg.LogDirectory, cfg.ArchiveDirectory, cfg.Host,
			cfg.PasscodePort, cfg.ArchivePort, cfg.Secret)
		log.Fatal("Environment variables missing")
	}

	if _, err := os.Stat(cfg.ArchiveDirectory); os.IsNotExist(err) {
		log.Fatalf("Archive directory %s does not exist", cfg.ArchiveDirectory)
	}
}
