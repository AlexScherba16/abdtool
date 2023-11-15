package misc

import (
	"abdtool/utils/errors"
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

func CreateErrorWithStackTrace(severity errors.SeverityLevel, message string, errorSourcePoint string, trace ...string) *errors.CustomError {
	err := errors.NewCustomError(errors.Critical, message, errorSourcePoint)

	for _, t := range trace {
		err.AppendStackTrace(t)
	}
	return err
}
