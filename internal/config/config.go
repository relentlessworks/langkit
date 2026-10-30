package config

import (
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
)

// Config holds all configuration for the langkit service.
type Config struct {
	Addr   string
	NoAuth bool
	Secret string
}

// Default returns a config with sensible defaults.
func Default() *Config {
	return &Config{
		Addr:   ":7700",
		NoAuth: false,
		Secret: "",
	}
}

// Load parses flags and env vars, layered: defaults < env < flags.
func Load() *Config {
	c := Default()

	// Env vars
	if v := os.Getenv("LANGKIT_ADDR"); v != "" {
		c.Addr = v
	}
	if v := os.Getenv("LANGKIT_NO_AUTH"); v == "true" || v == "1" {
		c.NoAuth = true
	}
	if v := os.Getenv("LANGKIT_SECRET"); v != "" {
		c.Secret = v
	}

	// Flags
	flag.StringVar(&c.Addr, "addr", c.Addr, "listen address")
	flag.BoolVar(&c.NoAuth, "no-auth", c.NoAuth, "disable auth (dev mode)")
	flag.StringVar(&c.Secret, "secret", c.Secret, "token signing secret (auto-generated if empty)")
	flag.Parse()

	// Generate random secret if not provided
	if c.Secret == "" {
		b := make([]byte, 32)
		rand.Read(b)
		c.Secret = hex.EncodeToString(b)
	}

	return c
}

func (c *Config) String() string {
	return fmt.Sprintf("addr=%s no-auth=%v", c.Addr, c.NoAuth)
}
