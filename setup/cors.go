package setup

import (
	"net/http"
	"os"

	"github.com/rs/cors"
	"gopkg.in/yaml.v2"
)

type OriginConfig struct {
	Origin  string   `yaml:"origin"`
	Methods []string `yaml:"methods"`
}

type CorsConfig struct {
	Origins []OriginConfig `yaml:"origins"`
}

func LoadCorsConfig(filename string) (CorsConfig, error) {
	var config CorsConfig

	data, err := os.ReadFile(filename)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(data, &config)
	return config, err
}

func CorsHandlerMiddleware(config CorsConfig) http.Handler {
	corsOptions := cors.Options{
		AllowedOrigins: []string{},
		AllowedMethods: []string{},
		AllowedHeaders: []string{"Content-Type"},
	}

	for _, originConfig := range config.Origins {
		corsOptions.AllowedOrigins = append(corsOptions.AllowedOrigins, originConfig.Origin)
		corsOptions.AllowedMethods = append(corsOptions.AllowedMethods, originConfig.Methods...)
	}

	c := cors.New(corsOptions)
	return c.Handler(http.DefaultServeMux)
}
