package conf

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

type Duration time.Duration

func (d *Duration) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	dur, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*d = Duration(dur)
	return nil
}

func (d *Duration) UnmarshalText(text []byte) error {
	dur, err := time.ParseDuration(string(text))
	if err != nil {
		return err
	}
	*d = Duration(dur)
	return nil
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

func (d Duration) AsDuration() time.Duration {
	return time.Duration(d)
}

type Bootstrap struct {
	Server *Server `yaml:"server"`
	Data   *Data   `yaml:"data"`
}

type Server struct {
	GRPC *ServerGRPC `yaml:"grpc"`
}

type ServerGRPC struct {
	Network string   `yaml:"network"`
	Addr    string   `yaml:"addr"`
	Timeout Duration `yaml:"timeout"`
}

type Data struct {
	Database *DataDatabase `yaml:"database"`
	Redis    *DataRedis    `yaml:"redis"`
}

type DataDatabase struct {
	Driver string `yaml:"driver"`
	Source string `yaml:"source"`
}

type DataRedis struct {
	Network      string   `yaml:"network"`
	Addr         string   `yaml:"addr"`
	ReadTimeout  Duration `yaml:"read_timeout"`
	WriteTimeout Duration `yaml:"write_timeout"`
}

// Load reads all yaml files in the config directory and returns a Bootstrap.
func Load(path string) (*Bootstrap, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	bc := &Bootstrap{}
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".yaml" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(path, entry.Name()))
		if err != nil {
			return nil, err
		}
		if err := yaml.Unmarshal(data, bc); err != nil {
			return nil, err
		}
	}
	return bc, nil
}
