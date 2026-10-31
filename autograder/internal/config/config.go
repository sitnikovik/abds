// Package config loads the settings the autograder needs to reach the student's services.
package config

import (
	"fmt"
	"strconv"
)

const (
	// keyAppPort is the environment variable with the HTTP port of the student's application.
	keyAppPort = "APP_PORT"
	// keyClickHouseHost is the environment variable with the ClickHouse host.
	keyClickHouseHost = "CLICKHOUSE_HOST"
	// keyClickHousePort is the environment variable with the ClickHouse native protocol port.
	keyClickHousePort = "CLICKHOUSE_PORT"
	// keyClickHouseDatabase is the environment variable with the ClickHouse database name.
	keyClickHouseDatabase = "CLICKHOUSE_DATABASE"
	// keyClickHouseUser is the environment variable with the ClickHouse user name.
	keyClickHouseUser = "CLICKHOUSE_USER"
	// keyClickHousePassword is the environment variable with the ClickHouse password.
	keyClickHousePassword = "CLICKHOUSE_PASSWORD"
)

// Config holds the settings of the student's application and ClickHouse server.
type Config struct {
	// appPort is the HTTP port of the student's application.
	appPort string
	// chHost is the host of the ClickHouse server.
	chHost string
	// chPort is the native protocol port of the ClickHouse server.
	chPort string
	// chDatabase is the database the student's tables live in.
	chDatabase string
	// chUser is the ClickHouse user name.
	chUser string
	// chPassword is the ClickHouse password, it may be empty.
	chPassword string
}

// Load reads the settings using lookup, which has the same contract as os.LookupEnv.
//
// It returns an error if a required variable is missing or a port is not a valid number.
func Load(lookup func(key string) (string, bool)) (Config, error) {
	required := []string{
		keyAppPort,
		keyClickHouseHost,
		keyClickHousePort,
		keyClickHouseDatabase,
		keyClickHouseUser,
	}
	env := make(map[string]string, len(required))
	for _, key := range required {
		val, ok := lookup(key)
		if !ok || val == "" {
			return Config{}, fmt.Errorf("environment variable %s is not set", key)
		}
		env[key] = val
	}
	for _, key := range []string{keyAppPort, keyClickHousePort} {
		if err := validatePort(env[key]); err != nil {
			return Config{}, fmt.Errorf("environment variable %s: %w", key, err)
		}
	}
	password, _ := lookup(keyClickHousePassword)
	return Config{
		appPort:    env[keyAppPort],
		chHost:     env[keyClickHouseHost],
		chPort:     env[keyClickHousePort],
		chDatabase: env[keyClickHouseDatabase],
		chUser:     env[keyClickHouseUser],
		chPassword: password,
	}, nil
}

// AppPort returns the HTTP port of the student's application.
func (c Config) AppPort() string {
	return c.appPort
}

// ClickHouseAddr returns the host and port of the ClickHouse native protocol endpoint,
// joined with a colon.
func (c Config) ClickHouseAddr() string {
	return c.chHost + ":" + c.chPort
}

// ClickHouseDatabase returns the name of the database the student's tables live in.
func (c Config) ClickHouseDatabase() string {
	return c.chDatabase
}

// ClickHouseUser returns the ClickHouse user name.
func (c Config) ClickHouseUser() string {
	return c.chUser
}

// ClickHousePassword returns the ClickHouse password, it may be empty.
func (c Config) ClickHousePassword() string {
	return c.chPassword
}

// validatePort returns an error if val is not a TCP port number from 1 to 65535.
func validatePort(val string) error {
	port, err := strconv.ParseUint(val, 10, 16)
	if err != nil {
		return fmt.Errorf("%q is not a number: %w", val, err)
	}
	if port == 0 {
		return fmt.Errorf("port must be from 1 to 65535, got %d", port)
	}
	return nil
}
