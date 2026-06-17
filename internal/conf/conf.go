package conf

import (
	"encoding/json"
	"time"
)

// Duration is a duration that can be unmarshaled from a string like "1s".
// It works with both YAML (via TextUnmarshaler) and JSON (via UnmarshalJSON).
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

// MarshalJSON returns the duration as a JSON string.
func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

// AsDuration returns the duration as a time.Duration.
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
