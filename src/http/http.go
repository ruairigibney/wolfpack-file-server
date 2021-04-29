package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/gorilla/sessions"
	"github.com/patrickmn/go-cache"
	"github.com/rs/cors"
	"github.com/ruairigibney/wolfpack-file-server/src/auth"
	"github.com/ruairigibney/wolfpack-file-server/src/config"
)

type getPassCodeRequestBody struct {
	Secret string
}

type Cfg struct {
	*config.Config
}

type FileDetails struct {
	FileName string
	ModTime  time.Time
}

func HttpHandler(c *config.Config) {
	crs := cors.AllowAll()

	cfg := Cfg{c}
	archiveMux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(c.ArchiveDirectory))
	archiveMux.Handle("/", cfg.RestrictedHandler((fileServer)))
	log.Printf("Serving archive %s on HTTP port: %v\n", c.ArchiveDirectory, c.ArchivePort)
	handler := cors.Default().Handler(archiveMux)

	go func() {
		log.Fatal(http.ListenAndServe(":"+c.ArchivePort, crs.Handler(handler)))
	}()

	passCodeMux := http.NewServeMux()
	passCodeMux.HandleFunc("/passcode", cfg.GetPassCode)
	log.Printf("Serving passcodes on HTTP port: %v\n", c.PasscodePort)
	log.Fatal(http.ListenAndServe(":"+c.PasscodePort, crs.Handler((passCodeMux))))
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
	url := fmt.Sprintf("%s:4200?passcode=%s", c.Host, passcode)
	resp.Write([]byte(url))
}

func (c *Cfg) RestrictedHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		resp = setHeaders(resp, c.Host)

		session, _ := c.Store.Get(req, "wolfpack-file-server")
		if session.Values["token"] == nil {
			if req.URL.Path != "/api/token" {
				http.Error(resp, "Not authenticated", http.StatusForbidden)
				return
			}
			status, err := c.getToken(*session, resp, req)
			if err != nil {
				http.Error(resp, err.Error(), status)
			}

			resp.WriteHeader(status)
			return
		}

		switch req.URL.Path {
		case "/api/files/list":
			c.listFiles(resp, req)
		case "/api/files/content":
			c.getFile(resp, req)
		}

	})
}

func (c *Cfg) getToken(session sessions.Session, resp http.ResponseWriter, req *http.Request) (int, error) {
	passcodeQueryParams := req.URL.Query()["passcode"]
	if len(passcodeQueryParams) == 0 {
		return http.StatusForbidden, errors.New("no passcode")
	}

	if _, found := c.C.Get(passcodeQueryParams[0]); !found {
		return http.StatusInternalServerError, errors.New("invalid OTP")
	} else {
		var err error
		session.Values["token"], err = auth.GenerateSecureKey(32)
		if err != nil {
			return http.StatusInternalServerError, errors.New("error generating token")
		}
		err = session.Save(req, resp)
		if err != nil {
			return http.StatusInternalServerError, errors.New("error saving session")
		}
		c.C.Delete(passcodeQueryParams[0])
	}

	return http.StatusOK, nil
}

func (c *Cfg) getFile(resp http.ResponseWriter, req *http.Request) {
	filenameQueryParams := req.URL.Query()["filename"]
	if len(filenameQueryParams) == 0 {
		http.Error(resp, "No file specified", http.StatusBadRequest)
		return
	}

	file, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", c.ArchiveDirectory, filenameQueryParams[0]))
	if err != nil {
		http.Error(resp, "Error reading file", http.StatusInternalServerError)
		return
	}

	resp.Header().Set("Content-Type", "text/html")
	resp.WriteHeader(http.StatusOK)
	resp.Write(file)
}

func setHeaders(resp http.ResponseWriter, host string) http.ResponseWriter {
	resp.Header().Set("Access-Control-Allow-Origin", fmt.Sprintf("%s:4200", host))
	resp.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	resp.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	resp.Header().Set("Access-Control-Allow-Credentials", "true")

	return resp
}

func (c *Cfg) listFiles(resp http.ResponseWriter, req *http.Request) {
	files, err := ioutil.ReadDir(c.ArchiveDirectory)
	if err != nil {
		http.Error(resp, "Error reading dir", http.StatusInternalServerError)
		return
	}

	var fD []FileDetails

	for _, v := range files {
		if v.Name() != "" {

			file := FileDetails{
				FileName: v.Name(),
				ModTime:  v.ModTime(),
			}

			fD = append(fD, file)
		}
	}

	sort.Slice(fD, func(i, j int) bool { return fD[i].ModTime.Unix() > fD[j].ModTime.Unix() })

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusOK)
	body, err := json.Marshal(fD)
	if err != nil {
		http.Error(resp, "Error marshalling json", http.StatusInternalServerError)
		return
	}

	resp.Write(body)
}
