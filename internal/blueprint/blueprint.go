package blueprint

import (
	"abdtool/internal/blueprint/parser"
	"abdtool/internal/blueprint/parser/yaml"
	"abdtool/internal/blueprint/types"
	"abdtool/utils/errors"
	"fmt"
	"path/filepath"
)

// buildParser - private fabric-function, constructs and returns an instance of blueprint file parser
//
// Parameters:
//   - path: blueprint file path
//
// Returns:
//   - IBlueprintParser: A concrete instance of IBlueprintParser
//   - *errors.CustomError: A custom error that provides detailed information if something went wrong during parser creation, nil if there's no error
func buildParser(path string) (parser.IBlueprintParser, *errors.CustomError) {
	source := "buildParser"

	// Get file extension
	ext := filepath.Ext(path)

	switch ext {
	case ".yaml":
		return yaml.NewParser(), nil

	default:
		message := fmt.Sprintf("Invalid blueprint file extention, %q is not implemented", ext)
		return nil, errors.NewCustomError(errors.Critical, message, source)
	}
}

// validateBlueprint - private function validates critical fields in application blueprint
//
// Parameters:
//   - b: parsed blueprint entity
//
// Returns:
//   - *errors.CustomError: A custom error that provides detailed information if some fields has wrong content, nil if there's no error
func validateBlueprint(b *types.Blueprint) *errors.CustomError {
	source := "validateBlueprint"

	// Validate Project entry
	if b.Project.Name == "" {
		message := fmt.Sprintf("%q must be provided", "project.name")
		return errors.NewCustomError(errors.Critical, message, source)
	}

	if b.Project.Version == "" {
		message := fmt.Sprintf("%q must be provided", "project.version")
		return errors.NewCustomError(errors.Critical, message, source)
	}

	// Validate nonempty routines collection
	if len(b.Routines) == 0 {
		message := fmt.Sprintf("%q must be provided", "routines")
		return errors.NewCustomError(errors.Critical, message, source)
	}

	for _, r := range b.Routines {
		// Validate Routine entry
		if r.Name == "" {
			message := fmt.Sprintf("%q must be provided", "routine.name")
			return errors.NewCustomError(errors.Critical, message, source)
		}
		if len(r.Steps) == 0 {
			message := fmt.Sprintf("%q must be provided for %q routine", "routine.steps", r.Name)
			return errors.NewCustomError(errors.Critical, message, source)
		}

		for _, s := range r.Steps {
			// Validate Step entry
			if s.Name == "" {
				message := fmt.Sprintf("%q must be provided in %q routine", "step.name", r.Name)
				return errors.NewCustomError(errors.Critical, message, source)
			}
			if s.Description == "" {
				message := fmt.Sprintf("%q must be provided for %q step", "step.description", s.Name)
				return errors.NewCustomError(errors.Critical, message, source)
			}
			if s.Timeout_s <= 0 {
				message := fmt.Sprintf("%q must be >= 0 in %q step", "step.timeout_s", s.Name)
				return errors.NewCustomError(errors.Critical, message, source)
			}
			if s.Attempts <= 0 {
				message := fmt.Sprintf("%q must be >= 0 in %q step", "step.attempts", s.Name)
				return errors.NewCustomError(errors.Critical, message, source)
			}
		}
	}

	return nil
}

// NewBlueprint - constructs and returns an instance of application blueprint,
// whole application configurator
//
// Parameters:
//   - path: blueprint file path
//
// Returns:
//   - Blueprint: The constructed Blueprint instance or empty struct if there's an error
//   - *errors.CustomError: A custom error that provides detailed information if something went wrong during the construction, nil if there's no error
func NewBlueprint(path string) (types.Blueprint, *errors.CustomError) {
	source := "NewBlueprint"

	// Create blueprint file parser
	parser, err := buildParser(path)
	if err != nil {
		err.AppendStackTrace(source)
		return types.Blueprint{}, err
	}

	// Parse blueprint file
	blueprint, err := parser.GetBlueprint(path)
	if err != nil {
		err.AppendStackTrace(source)
		return types.Blueprint{}, err
	}

	// Validate blueprint result
	if err := validateBlueprint(&blueprint); err != nil {
		err.AppendStackTrace(source)
		return types.Blueprint{}, err
	}

	return blueprint, nil
}
