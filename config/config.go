package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"runtime"
)

type Configuration struct {
	Bootstrappers []*BootstrapperConf `json:"bootstrappers"`
	Services      []*ServiceConf      `json:"services"`
}

type BootstrapperConf struct {
	EcdsaHexPrivateKey string `json:"ecdsaHexPrivateKey,omitempty"`
}

type ServiceConf struct {
	EcdsaHexPrivateKey string `json:"ecdsaHexPrivateKey,omitempty"`
	EcdsaHexPublicKey  string `json:"ecdsaHexPublicKey,omitempty"`
}

func ReadConf(fpath string) (Configuration, error) {
	var conf Configuration

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return conf, fmt.Errorf("No caller information")
	}
	if !filepath.IsAbs(fpath) {
		fpath = path.Join(path.Dir(filename), "..", fpath)
	}

	jsonBytes, err := ioutil.ReadFile(fpath)
	if err != nil {
		return conf, fmt.Errorf("couldn't open %q, please ensure it exists", fpath)
	}
	if err := json.Unmarshal(jsonBytes, &conf); err != nil {
		return conf, err
	}

	return conf, nil
}
