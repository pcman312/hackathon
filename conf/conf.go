package conf

import (
	"net/url"
	"time"

	"github.com/joho/godotenv"
	"github.com/pcman312/env"
)

// Opts for the configurator
type Opts struct {
	Port            int           `env:"grpc.port"         default:"9090"  min:"1" max:"65535"`
	GRPCConnTimeout time.Duration `env:"grpc.conn.timeout" default:"120s"`

	ESHost       *url.URL      `env:"es.host"       required:"true"`
	ESSkipVerify bool          `env:"es.skipVerify"`
	ESTimeout    time.Duration `env:"es.timeout"    default:"30s"`
	ESIndex      string        `env:"es.index"      default:"hackathon"`

	ShutdownWait time.Duration `env:"shutdownWait" default:"5s"  min:"1s"`
}

// LoadOpts from the given environment file
func LoadOpts(envFile string) (opts Opts, err error) {
	// Load the envFile into the environment, but don't overwrite any values that are set in the env
	err = godotenv.Load(envFile)
	if err != nil {
		return opts, err
	}

	// Parse it into the struct
	err = env.Parse(&opts)
	if err != nil {
		return opts, err
	}

	return opts, nil
}
