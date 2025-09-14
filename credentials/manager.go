package credentials

import (
	"encoding/json"
	"os"

	"github.com/adrg/xdg"
)

// Credentials represents the authentication credentials
type Credentials struct {
	Token   string `json:"token"`
	TokenId string `json:"token_id"`
}

// Save saves credentials to the XDG config directory
func (c *Credentials) Save() error {
	// Get the config file path according to XDG specification
	configFilePath, err := xdg.ConfigFile("syojctl/credentials.json")
	if err != nil {
		return err
	}

	// Marshal credentials to JSON
	credsJSON, err := json.Marshal(c)
	if err != nil {
		return err
	}

	// Write credentials to file
	return os.WriteFile(configFilePath, credsJSON, 0600)
}

// Load loads credentials from the XDG config directory
func Load() (*Credentials, error) {
	// Search for the config file according to XDG specification
	configFilePath, err := xdg.SearchConfigFile("syojctl/credentials.json")
	if err != nil {
		return nil, err
	}

	// Read the file
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON to credentials struct
	var creds Credentials
	err = json.Unmarshal(data, &creds)
	if err != nil {
		return nil, err
	}

	return &creds, nil
}

// Delete removes the credentials file
func Delete() error {
	// Search for the config file according to XDG specification
	configFilePath, err := xdg.SearchConfigFile("syojctl/credentials.json")
	if err != nil {
		return err
	}

	// Remove the file
	return os.Remove(configFilePath)
}