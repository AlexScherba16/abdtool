package types

import (
	"abdtool/utils/errors"
	"fmt"

	"gopkg.in/yaml.v3"
)

// ProjectEntry - represents essential details about project. It captures
// the project's name and its version
type ProjectEntry struct {
	// Name - contains specified title of the project
	Name string `yaml:"name"`

	// Version - denotes the version string of the project
	Version string `yaml:"version"`
}

// NewProjectEntry - constructs and returns an ProjectEntry instance,
// with project name and it's version
//
// Returns:
//   - ProjectEntry: The constructed ProjectEntry instance or empty struct if there's an error
//   - *errors.CustomError: A custom error that provides detailed information if something went wrong during the construction, nil if there's no error
func NewProjectEntry(data []byte) (ProjectEntry, *errors.CustomError) {
	source := "NewProjectEntry"

	// Using anonimous struct for yaml parser
	tmp := struct {
		Project ProjectEntry `yaml:"project"`
	}{}

	// Read project entry data
	err := yaml.Unmarshal(data, &tmp)
	if err != nil {
		message := fmt.Sprintf("Failed to unmarshall yaml file, %s", err.Error())
		err := errors.NewCustomError(errors.Critical, message, source)
		return ProjectEntry{}, err
	}

	// Validate provided fields
	if tmp.Project.Name == "" {
		message := fmt.Sprintf("%q must be provided", "project.name")
		err := errors.NewCustomError(errors.Critical, message, source)
		return ProjectEntry{}, err
	}

	if tmp.Project.Version == "" {
		message := fmt.Sprintf("%q must be provided", "project.version")
		err := errors.NewCustomError(errors.Critical, message, source)
		return ProjectEntry{}, err
	}
	return tmp.Project, nil
}
