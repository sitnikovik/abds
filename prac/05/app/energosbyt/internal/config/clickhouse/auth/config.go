package auth

// Config holds ClickHouse authentication credentials.
type Config struct {
	// username stores the login username for ClickHouse.
	username string

	// password stores the login password for ClickHouse.
	password string
}

// NewConfig creates a new Config with the given credentials.
func NewConfig(
	username string,
	password string,
) Config {
	return Config{
		username: username,
		password: password,
	}
}

// Username returns the login username.
func (c Config) Username() string {
	return c.username
}

// Password returns the login password.
func (c Config) Password() string {
	return c.password
}
