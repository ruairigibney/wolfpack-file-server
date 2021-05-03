package config

import (
	"github.com/gorilla/sessions"
	"github.com/patrickmn/go-cache"
)

type Config struct {
	LogDirectory     string `envconfig:"LOG_DIRECTORY"`
	ArchiveDirectory string `envconfig:"ARCHIVE_DIRECTORY"`
	PasscodePort     string `envconfig:"PASSCODE_PORT"`
	ArchivePort      string `envconfig:"ARCHIVE_PORT"`
	FrontEndPort     string `envconfig:"FRONTEND_PORT"`
	Secret           string `envconfig:"SECRET"`
	Host             string `envconfig:"HOST"`
	Store            *sessions.CookieStore
	C                *cache.Cache
}
