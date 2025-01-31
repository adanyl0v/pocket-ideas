package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"os"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
}

type Reader interface {
	Read() (*Config, error)
}

type fileReader struct {
	path string
}

func NewFileReader(path string) Reader {
	return &fileReader{path: path}
}

func (f *fileReader) Read() (*Config, error) {
	cfg := new(Config)
	if err := cleanenv.ReadConfig(f.path, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func DefaultFilePath() string {
	return os.Getenv("CONFIG_FILE")
}

func ReadFile(path string) (*Config, error) {
	cfg, err := NewFileReader(path).Read()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func MustReadFile(path string) *Config {
	cfg, err := ReadFile(path)
	if err != nil {
		panic(err)
	}

	return cfg
}
