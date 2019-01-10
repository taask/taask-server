package config

import (
	"os"
	"os/user"
	"path"

	log "github.com/cohix/simplog"
	"github.com/pkg/errors"
	"github.com/taask/taask-server/auth"
)

// ConfigServerBaseDir and others are config related consts
const (
	ConfigServerBaseDir = ".taask/server/config/"

	ClientAuthConfigFilename = "client-auth.yaml"

	missingAdminGroupConfigWarning = `
########################################################
	ADMIN GROUP NOT CONFIGURED.
	EMPTY PASSPHRASE WILL BE USED.
	JOIN CODE WILL BE PRINTED.
RUN 'taaskctl generate auth' TO GENERATE CREDENTIALS
EMPTY PASSPHRASE IS INSECURE. DO NOT RUN IN PRODUCTION.
########################################################
`
)

// ServerConfig is the config for the server
type ServerConfig struct {
	ClientAuth *ClientAuthConfig
}

// ServerConfigFromDefaultDir reads the server config from the default directory
func ServerConfigFromDefaultDir() (*ServerConfig, error) {
	clientAuthConfig, err := clientAuthConfigFromFile(path.Join(DefaultServerConfigDir(), ClientAuthConfigFilename))
	if err != nil {
		log.LogWarn(missingAdminGroupConfigWarning)

		joinCode := auth.GenerateJoinCode()
		authHash := auth.GroupAuthHash(joinCode, "")

		group := auth.MemberGroup{
			UUID:     auth.AdminGroupUUID,
			Name:     "admin",
			JoinCode: joinCode,
			AuthHash: authHash,
		}

		clientAuthConfig = &ClientAuthConfig{
			Version:     ClientAuthConfigVersion,
			Type:        ClientAuthConfigType,
			MemberGroup: group,
		}

		if err := clientAuthConfig.WriteYAML(path.Join(DefaultServerConfigDir(), ClientAuthConfigFilename)); err != nil {
			log.LogWarn(errors.Wrap(err, "failed to WriteYaml for generated admin config ").Error())
		}
	}

	config := &ServerConfig{
		ClientAuth: clientAuthConfig,
	}

	return config, nil
}

// DefaultServerConfigDir returns ~/.taask/server/config unless XDG_CONFIG_HOME is set
func DefaultServerConfigDir() string {
	u, err := user.Current()
	if err != nil {
		return ""
	}

	root := u.HomeDir
	xdgConfig, useXDG := os.LookupEnv("XDG_CONFIG_HOME")
	if useXDG {
		root = xdgConfig
	}

	return path.Join(root, ConfigServerBaseDir)
}
