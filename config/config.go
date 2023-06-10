/*
Package config is the internal configuration used for Efs2. This configuration is for the internal application execution. It exists to pave the way for non-cli instances of Efs2 in the future.
*/
package config

// Config provides a configuration structure used within the Efs2 application.
type Config struct {
	Verbose bool

	GatewayHost string

	GatewayToken string

	DebugSecretToken string

	GatewayMode string
}

// New will return Config populated with pre-defined defaults.
func New() Config {
	c := Config{}
	c.Verbose = true
	c.GatewayHost = "http://localhost:9992"
	c.GatewayToken = ""
	c.DebugSecretToken = "sekret"
	return c
}
