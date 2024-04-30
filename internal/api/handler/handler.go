package handler

import (
	"encoding/json"
	"html/template"
	"io"
	"io/fs"
	"log"
)

type AuthProvider string

const (
	NONE  AuthProvider = "none"
	BASIC AuthProvider = "basic"
)

type Config struct {
	Base         string
	Addr         string
	Version      string
	Hostname     string
	AuthProvider AuthProvider
}

type Handler struct {
	config    *Config
	indexTmpl *template.Template
}

func NewHandler(config *Config, assets fs.FS) *Handler {
	file, err := assets.Open("index.html")
	if err != nil {
		log.Fatal(err)
	}
	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Panic(err)
	}

	tmpl, err := template.New("index.html").Funcs(template.FuncMap{
		"marshal": func(v interface{}) template.JS {
			var p []byte
			p, _ = json.Marshal(v)
			return template.JS(p)
		},
	}).Parse(string(bytes))
	if err != nil {
		log.Panic(err)
	}

	return &Handler{
		config:    config,
		indexTmpl: tmpl,
	}
}
