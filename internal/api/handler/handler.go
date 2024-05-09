package handler

import (
	"encoding/json"
	"html/template"
	"io"
	"io/fs"
	"log"

	"github.com/mehranzand/pulseup/internal/docker"
	"gorm.io/gorm"
)

type AuthProvider string

const (
	NONE  AuthProvider = "none"
	BASIC AuthProvider = "basic"
)

type Config struct {
	Base         string       `json:"base"`
	Adderss      string       `json:"address"`
	Version      string       `json:"version"`
	Hostname     string       `json:"host_name"`
	AuthProvider AuthProvider `json:"auth_provider"`
}

type Handler struct {
	config    *Config
	indexTmpl *template.Template
	clients   map[string]docker.Client
	db        *gorm.DB
}

func NewHandler(config *Config, clients map[string]docker.Client, assets fs.FS, db *gorm.DB) *Handler {
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
		clients:   clients,
		db:        db,
	}
}
