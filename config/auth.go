package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/taask/taask-server/auth"
	yaml "gopkg.in/yaml.v2"
)

// ClientAuthConfigType describe consts for server config
const (
	MemberAuthConfigVersion = 1
	MemberAuthConfigType    = "com.taask.config.memberauth"
)

// ClientAuthConfig is the config for client auth
type ClientAuthConfig struct {
	Version     int
	Type        string
	MemberGroup auth.MemberGroup
}

func clientAuthConfigFromFile(filepath string) (*ClientAuthConfig, error) {
	raw, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to ReadFile")
	}

	config := &ClientAuthConfig{}
	if err := yaml.Unmarshal(raw, config); err != nil {
		if jsonErr := json.Unmarshal(raw, config); jsonErr != nil {
			return nil, errors.Wrap(jsonErr, errors.Wrap(err, "failed to yaml and json Unmarshal").Error()) // stupid, but whatever
		}
	}

	return config, nil
}

// WriteYAML writes the YAML marshalled config to disk
func (ca *ClientAuthConfig) WriteYAML(filepath string) error {
	rawYAML, err := yaml.Marshal(ca)
	if err != nil {
		return errors.Wrap(err, "failed to yaml.Marshal")
	}

	if err := ioutil.WriteFile(filepath, rawYAML, 0666); err != nil {
		return errors.Wrap(err, "failed to WriteFile")
	}

	return nil
}
