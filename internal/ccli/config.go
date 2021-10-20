package scli

import (
	"fmt"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
)

func GetConfigPath(configFile string) (string, error) {
	if configFile != "" {
		// Use config file from the flag.
		_, err := os.Stat(configFile)
		if err != nil {
			return "", err
		}

		return configFile, nil
	}

	// Try to find config in home directory
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	homeConfigFile := filepath.Join(home, ".scli")
	_, err = os.Stat(homeConfigFile)
	if err == nil {
		return homeConfigFile, nil
	}

	etcConfigFile := "/etc/scli/config.yaml"
	_, err = os.Stat(etcConfigFile)
	if err == nil {
		return etcConfigFile, nil
	}

	return "", fmt.Errorf("config file not found in default locations: %s, %s", homeConfigFile, etcConfigFile)
}
