package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"runtime"
)

type configuration struct {
	Bootstrappers []*bootstrapperConf `json:"bootstrappers"`
	Services      []*serviceConf      `json:"services"`
}

type bootstrapperConf struct {
	EcdsaHexPrivateKey string `json:"ecdsaHexPrivateKey,omitempty"`
}

type serviceConf struct {
	EcdsaHexPrivateKey string `json:"ecdsaHexPrivateKey,omitempty"`
	EcdsaHexPublicKey  string `json:"ecdsaHexPublicKey,omitempty"`
}

func readConf() (configuration, error) {
	var conf configuration

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return conf, fmt.Errorf("No caller information")
	}
	if !filepath.IsAbs(configFilePath) {
		configFilePath = path.Join(path.Dir(filename), "..", configFilePath)
	}

	jsonBytes, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return conf, err
	}
	if err := json.Unmarshal(jsonBytes, &conf); err != nil {
		return conf, err
	}

	return conf, nil
}
