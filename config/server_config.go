package config

// ServerConf is the config struct for the main server properties
type ServerConf struct {
	SSHPort    *string
	SSHKeyPath *string
	SSHKeyName *string

	HTTPPort          *string
	HTTPFileServerDir *string
}

// default values if not provided in config file
var (
	SSHPort    = "22222"
	SSHKeyPath = "./config/resources/"
	SSHKeyName = "sshtron.pem"

	HTTPPort          = "8080"
	HTTPFileServerDir = "./static/"
)
