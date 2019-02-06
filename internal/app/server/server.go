package server

import (
	"bytes"
	"strings"
	"net/http"
	"encoding/json"
	
	"github.com/go-chi/chi"

	"github.com/kcq/poc-ipblock-pool/pkg/pool"
)

const (
	serverAddr = ":5555"
	paramDelay = "delay"
	paramPretty = "pretty"
	paramBlock = "block"
	paramKey = "key"
	pathPoolAllocation = "/pool/allocation"
)

type App struct {
	pm     *pool.Manager
	router *chi.Mux
}

func New(pmanager *pool.Manager) *App {
	app := &App{
		pm: pmanager,
	}

	app.init()
	return app
}

func (a *App) init() {
	a.router = chi.NewRouter()

	a.router.Get(pathPoolAllocation, func(w http.ResponseWriter, r *http.Request) {
		pretty := false
		if strings.ToLower(r.URL.Query().Get(paramPretty)) == "true" {
			pretty = true
		}

		var block string
		if r.URL.Query().Get(paramBlock) != "" {
			block = r.URL.Query().Get(paramBlock)
		}

		var key string
		if r.URL.Query().Get(paramKey) != "" {
			key = r.URL.Query().Get(paramKey)
		}

		blockInfo := a.pm.Lookup(block,key)

		if blockInfo == nil {
			reply(w, r,http.StatusNotFound)
		} else {
			replyJSON(w,r, blockInfo, http.StatusOK, pretty)
		}
	})

	a.router.Post(pathPoolAllocation, func(w http.ResponseWriter, r *http.Request) {
		delayUnlock := false
		if strings.ToLower(r.URL.Query().Get(paramDelay)) == "true" {
			delayUnlock = true
		}

		pretty := false
		if strings.ToLower(r.URL.Query().Get(paramPretty)) == "true" {
			pretty = true
		}

		key := ""
		if r.URL.Query().Get(paramKey) != "" {
			key = r.URL.Query().Get(paramKey)
		}

		blockInfo := a.pm.Allocate(key,delayUnlock)

		replyJSON(w,r, blockInfo, http.StatusOK, pretty)

	})

	a.router.Delete(pathPoolAllocation, func(w http.ResponseWriter, r *http.Request) {
		key := ""
		if r.URL.Query().Get(paramKey) != "" {
			key = r.URL.Query().Get(paramKey)
		}

		block := ""
		if r.URL.Query().Get(paramBlock) != "" {
			block = r.URL.Query().Get(paramBlock)
		}

		err := a.pm.Free(block,key)

		switch err {
			case pool.ErrBlockNotFound: reply(w, r,http.StatusNotFound)
			case nil: reply(w, r,http.StatusNoContent)
			default: reply(w, r,http.StatusInternalServerError)
		}
	})
}

func (a *App) Run() {
	if err := http.ListenAndServe(serverAddr, a.router); err != nil {
		panic(err)
	}	
}

func replyJSON(w http.ResponseWriter, 
			   r *http.Request, 
			   value interface{}, 
			   status int, 
			   pretty bool) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)

	if pretty {
		enc.SetIndent("","  ")
	}

	if err := enc.Encode(value); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if status > 0 {
		w.WriteHeader(status)
	}
	w.Write(buf.Bytes())
}

func reply(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
}

