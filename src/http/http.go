package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/patrickmn/go-cache"
	"github.com/ruairigibney/wolfpack-file-server/src/auth"
	"github.com/ruairigibney/wolfpack-file-server/src/config"
)

type cookieStore struct {
	*sessions.CookieStore
}

type getPassCodeRequestBody struct {
	Secret string
}

type Cfg struct {
	*config.Config
}

func HttpHandler(c *config.Config) {
	cfg := Cfg{c}
	archiveMux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(c.ArchiveDirectory))
	archiveMux.Handle("/", cfg.RestrictedHandler((fileServer)))
	log.Printf("Serving archive %s on HTTP port: %v\n", c.ArchiveDirectory, c.ArchivePort)
	go func() {
		log.Fatal(http.ListenAndServe(":"+c.ArchivePort, archiveMux))
	}()

	passCodeMux := http.NewServeMux()
	passCodeMux.HandleFunc("/passcode", cfg.GetPassCode)
	log.Printf("Serving passcodes on HTTP port: %v\n", c.PasscodePort)
	log.Fatal(http.ListenAndServe(":"+c.PasscodePort, passCodeMux))
}

func (c *Cfg) GetPassCode(resp http.ResponseWriter, req *http.Request) {
	log.Print("getting passcode")

	var requestBody getPassCodeRequestBody

	err := json.NewDecoder(req.Body).Decode(&requestBody)
	if err != nil {
		log.Printf("error: %v", err.Error())
		http.Error(resp, "Error decoding request", http.StatusBadRequest)
		return
	}

	if requestBody.Secret != c.Secret {
		http.Error(resp, "Invalid secret", http.StatusBadRequest)
		return
	}

	passcode, err := auth.GenerateSecureKey(32)
	if err != nil {
		http.Error(resp, "Error generating key", http.StatusInternalServerError)
		return
	}
	c.C.Add(passcode, nil, cache.DefaultExpiration)
	url := fmt.Sprintf("%s:%s/?passcode=%s", c.Host, c.ArchivePort, passcode)
	resp.Write([]byte(url))
	return
}

func (c *Cfg) RestrictedHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		session, _ := c.Store.Get(req, "wolfpack-file-server")
		if session.Values["token"] == nil {
			passcodeQueryParams := req.URL.Query()["passcode"]
			if len(passcodeQueryParams) == 0 {
				http.Error(resp, "No passcode", http.StatusForbidden)
				return
			}

			if _, found := c.C.Get(passcodeQueryParams[0]); !found {
				http.Error(resp, "Invalid OTP", http.StatusInternalServerError)
				return
			} else {
				var err error
				session.Values["token"], err = auth.GenerateSecureKey(32)
				if err != nil {
					http.Error(resp, "Error generating token", http.StatusInternalServerError)
					return
				}
				err = session.Save(req, resp)
				if err != nil {
					http.Error(resp, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}

		h.ServeHTTP(resp, req)
	})
}
