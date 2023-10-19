package mics

import (
	"os"

	"gopkg.in/yaml.v3"
)

const (
	FileName             = "tmp.inv.abc.*.tmp"
	ValidYamlFile        = "tmp.*.yaml"
	InvalidYamlExtention = ".qwe"
)

func CreateTmpYAML(name string, data interface{}) (string, error) {
	tmpFile, err := os.CreateTemp("", name)
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	encoder := yaml.NewEncoder(tmpFile)
	if err := encoder.Encode(data); err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

func DeleteTmpYAML(path string) error {
	return os.Remove(path)
}
