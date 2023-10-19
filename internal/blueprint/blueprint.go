package blueprint

import (
	"abdtool/internal/blueprint/types"
	"abdtool/utils/errors"
	"fmt"
	"os"
	"path/filepath"
)

// Blueprint - encapsulates the configuration blueprint of an application
type Blueprint struct {
	// project - represents the essential information about the project
	project types.ProjectEntry
}

// Project - getter function that returns ProjectEntry entity
//
// Returns:
//   - types.ProjectEntry: General project info
func (b Blueprint) Project() types.ProjectEntry {
	return b.project
}

// NewBlueprint - constructs and returns an instance of application blueprint,
// whole application configurator
//
// Returns:
//   - Blueprint: The constructed Blueprint instance or empty struct if there's an error
//   - *errors.CustomError: A custom error that provides detailed information if something went wrong during the construction, nil if there's no error
func NewBlueprint(path string) (Blueprint, *errors.CustomError) {
	source := "NewBlueprint"

	// Check for file exists
	_, e := os.Stat(path)
	if e != nil {
		err := errors.NewCustomError(errors.Critical, "Blueprint file doesn't exist", source)
		return Blueprint{}, err
	}

	// Get file extension
	ext := filepath.Ext(path)
	if ext != ".yaml" {
		message := fmt.Sprintf("Invalid plueprint file extention, expected .yaml, got %s", ext)
		err := errors.NewCustomError(errors.Critical, message, source)
		return Blueprint{}, err
	}

	// Read blueprint file content
	data, e := os.ReadFile(path)
	if e != nil {
		message := fmt.Sprintf("Failed to read blueprint file, %s", e.Error())
		err := errors.NewCustomError(errors.Critical, message, source)
		return Blueprint{}, err
	}

	project, err := types.NewProjectEntry(data)
	if err != nil {
		err.AppendStackTrace(source)
		return Blueprint{}, err
	}

	return Blueprint{
		project: project,
	}, nil
}
