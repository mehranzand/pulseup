package handler

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
	config *Config
}

func NewHandler(config *Config) *Handler {
	return &Handler{config: config}
}
