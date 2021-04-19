package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/patrickmn/go-cache"
)

type getPassCodeRequestBody struct {
	Secret string
}

var (
	store sessions.CookieStore
	c     *cache.Cache

	logDirectory     string
	archiveDirectory string
	passcodePort     string
	archivePort      string
	secret           string
	host             string
)

func main() {
	c = cache.New(24*time.Hour, 10*time.Minute)
	cookieKey := securecookie.GenerateRandomKey(16)
	store = *sessions.NewCookieStore(cookieKey)
	store.Options = &sessions.Options{
		Path:   "/",
		MaxAge: int(time.Hour * 24),
	}

	logDirectory = os.Getenv("FS_LOG_DIRECTORY")
	archiveDirectory = os.Getenv("FS_ARCHIVE_DIRECTORY")
	passcodePort = os.Getenv("FS_PASSCODE_PORT")
	archivePort = os.Getenv("FS_ARCHIVE_PORT")
	secret = os.Getenv("FS_SECRET")
	host = os.Getenv("FS_HOST")

	if logDirectory == "" || archiveDirectory == "" || host == "" ||
		passcodePort == "" || archivePort == "" || secret == "" {
		log.Printf(`logDirectory: %s; archiveDirectory: %s; host: %s; passCodePort: %s
			archivePort: %s; secret: %s`, logDirectory, archiveDirectory, host,
			passcodePort, archivePort, secret)
		log.Fatal("Environment variables missing")
	}

	if _, err := os.Stat(archiveDirectory); os.IsNotExist(err) {
		log.Fatalf("Archive directory %s does not exist", archiveDirectory)
	}

	logFile, err := os.OpenFile("wolfpack-file-server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logFile)

	archiveMux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(archiveDirectory))
	archiveMux.Handle("/", RestrictedHandler((fileServer)))
	log.Printf("Serving archive %s on HTTP port: %v\n", archiveDirectory, archivePort)
	go func() {
		log.Fatal(http.ListenAndServe(":"+archivePort, archiveMux))
	}()

	passCodeMux := http.NewServeMux()
	passCodeMux.HandleFunc("/passcode", getPassCode)
	log.Printf("Serving passcodes on HTTP port: %v\n", passcodePort)
	log.Fatal(http.ListenAndServe(":"+passcodePort, passCodeMux))
}

func getPassCode(resp http.ResponseWriter, req *http.Request) {
	log.Print("getting passcode")

	var requestBody getPassCodeRequestBody

	err := json.NewDecoder(req.Body).Decode(&requestBody)
	if err != nil {
		log.Printf("error: %v", err.Error())
		http.Error(resp, "Error decoding request", http.StatusBadRequest)
		return
	}

	if requestBody.Secret != secret {
		http.Error(resp, "Invalid secret", http.StatusBadRequest)
		return
	}

	passcode, err := generateSecureKey(20)
	if err != nil {
		http.Error(resp, "Error generating key", http.StatusInternalServerError)
		return
	}
	c.Add(passcode, nil, cache.DefaultExpiration)
	url := fmt.Sprintf("%s:%s/?passcode=%s", host, archivePort, passcode)
	resp.Write([]byte(url))
	return
}

func RestrictedHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		session, _ := store.Get(req, "wolfpack-file-server")
		if session.Values["token"] == nil {
			passcodeQueryParams := req.URL.Query()["passcode"]
			if len(passcodeQueryParams) == 0 {
				http.Error(resp, "No passcode", http.StatusForbidden)
				return
			}

			if _, found := c.Get(passcodeQueryParams[0]); !found {
				http.Error(resp, "Invalid OTP", http.StatusInternalServerError)
				return
			} else {
				var err error
				session.Values["token"], err = generateSecureKey(20)
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

func generateSecureKey(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", errors.New("Error generating key")
	}
	return hex.EncodeToString(b), nil
}
