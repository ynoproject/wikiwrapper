package setup

import (
	"context"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
)

type contextKey string

const ConfigKey = contextKey("config")

type Game struct {
	Name         string            `yaml:"name"`
	Namespace    string            `yaml:"namespace"`
	Protagonists map[string]string `yaml:"protagonists"`
}

type Protagonist struct {
	Category string `yaml:"category"`
}

type WikiConfig struct {
	Games map[string]Game `yaml:"games"`
}

func LoadWikiConfig(filename string) (WikiConfig, error) {
	var config WikiConfig

	data, err := os.ReadFile(filename)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(data, &config)

	return config, err
}

func WikiConfigHandlerMiddleware(config WikiConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), ConfigKey, config)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
