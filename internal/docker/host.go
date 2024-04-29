package docker

import (
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type Host struct {
	Name       string   `json:"name"`
	ID         string   `json:"id"`
	URL        *url.URL `json:"-"`
	CertPath   string   `json:"-"`
	CACertPath string   `json:"-"`
	KeyPath    string   `json:"-"`
	ValidCerts bool     `json:"-"`
}

func ParseConnection(connection string) (Host, error) {
	remoteUrl, err := url.Parse(connection)
	if err != nil {
		return Host{}, err
	}

	basePath, err := filepath.Abs("./certs")
	if err != nil {
		log.Fatalf("error converting certs path to absolute: %s", err)
	}

	host := remoteUrl.Hostname()
	if _, err := os.Stat(filepath.Join(basePath, host)); !os.IsNotExist(err) {
		basePath = filepath.Join(basePath, host)
	}

	cacertPath := filepath.Join(basePath, "ca.pem")
	certPath := filepath.Join(basePath, "cert.pem")
	keyPath := filepath.Join(basePath, "key.pem")

	hasCerts := true
	if _, err := os.Stat(cacertPath); os.IsNotExist(err) {
		cacertPath = ""
		hasCerts = false
	}

	return Host{
		ID:         strings.ReplaceAll(remoteUrl.String(), "/", ""),
		Name:       connection,
		URL:        remoteUrl,
		CertPath:   certPath,
		CACertPath: cacertPath,
		KeyPath:    keyPath,
		ValidCerts: hasCerts,
	}, nil
}
