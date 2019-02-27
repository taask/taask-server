package config

import (
	"fmt"
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

	ClientAuthConfigFilename  = "client-auth.yaml"
	PartnerAuthConfigFilename = "partner-auth.yaml"
	RunnerAuthConfigFilename  = "runner-auth.yaml"

	missingAdminGroupConfigWarning = `
########################################################
	ADMIN GROUP NOT CONFIGURED.
	EMPTY PASSPHRASE WILL BE USED.
RUN 'taaskctl init server' TO GENERATE CREDENTIALS
EMPTY PASSPHRASE IS INSECURE. DO NOT RUN IN PRODUCTION.
########################################################
`
)

var clientConfigPath = path.Join(DefaultServerConfigDir(), ClientAuthConfigFilename)
var partnerConfigPath = path.Join(DefaultServerConfigDir(), PartnerAuthConfigFilename)
var runnerConfigPath = path.Join(DefaultServerConfigDir(), RunnerAuthConfigFilename)

// ServerConfig is the config for the server
type ServerConfig struct {
	ClientAuth  *ClientAuthConfig
	PartnerAuth *ClientAuthConfig
	RunnerAuth  *ClientAuthConfig
}

// ServerConfigFromDefaultDir reads the server config from the default directory
func ServerConfigFromDefaultDir() (*ServerConfig, error) {
	clientAuthConfig, err := clientAuthConfigFromFile(clientConfigPath)
	if err != nil {
		log.LogWarn(missingAdminGroupConfigWarning)
		log.LogInfo(fmt.Sprintf("writing insecure client config to: %s", runnerConfigPath))

		clientAuthConfig = generateInsecureAdminGroup()

		if err := clientAuthConfig.WriteYAML(path.Join(DefaultServerConfigDir(), ClientAuthConfigFilename)); err != nil {
			log.LogWarn(errors.Wrap(err, "failed to WriteYaml for generated admin config ").Error())
		}
	} else {
		log.LogInfo("loaded client auth from file")
	}

	partnerAuthConfig, err := clientAuthConfigFromFile(partnerConfigPath)
	if err != nil {
		log.LogWarn("partner auth file not found, federation will not be available")
	} else {
		log.LogInfo("loaded partner auth from file")
	}

	runnerAuthConfig, err := clientAuthConfigFromFile(runnerConfigPath)
	if err != nil {
		runnerAuthConfig = generateDefaultRunnerGroup()

		log.LogInfo(fmt.Sprintf("writing runner config to: %s", runnerConfigPath))

		if err := runnerAuthConfig.WriteYAML(runnerConfigPath); err != nil {
			log.LogWarn(errors.Wrap(err, "failed to WriteYaml for generated runner config ").Error())
		}
	} else {
		log.LogInfo("loaded runner auth from file")
	}

	config := &ServerConfig{
		ClientAuth:  clientAuthConfig,
		PartnerAuth: partnerAuthConfig,
		RunnerAuth:  runnerAuthConfig,
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

func generateInsecureAdminGroup() *ClientAuthConfig {
	joinCode := auth.GenerateJoinCode()
	authHash := auth.GroupAuthHash(joinCode, "")

	group := auth.MemberGroup{
		UUID:     auth.AdminGroupUUID,
		Name:     "admin",
		JoinCode: joinCode,
		AuthHash: authHash,
	}

	clientAuthConfig := &ClientAuthConfig{
		Version:     MemberAuthConfigVersion,
		Type:        MemberAuthConfigType,
		MemberGroup: group,
	}

	return clientAuthConfig
}

func generateDefaultRunnerGroup() *ClientAuthConfig {
	joinCode := auth.GenerateJoinCode()
	authHash := auth.GroupAuthHash(joinCode, "")

	group := auth.MemberGroup{
		UUID:     auth.DefaultGroupUUID,
		Name:     "default",
		JoinCode: joinCode,
		AuthHash: authHash,
	}

	clientAuthConfig := &ClientAuthConfig{
		Version:     MemberAuthConfigVersion,
		Type:        MemberAuthConfigType,
		MemberGroup: group,
	}

	return clientAuthConfig
}
