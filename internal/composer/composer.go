package composer

import (
	"abdtool/internal/application"
	"abdtool/internal/blueprint"
	"abdtool/utils/errors"
)

// NewComposedApplication - constructs and returns an instance of a composed application
// that implements the IApplication interface and ready to lauch
//
// Returns:
//   - application.IApplication: The constructed application instance or nil if there's an error or it's not yet implemented
//   - *errors.CustomError: A custom error that provides detailed information if something went wrong during the construction, nil if there's no error
func NewComposedApplication() (application.IApplication, *errors.CustomError) {
	source := "NewComposedApplication"

	// Set global error verbosity
	errors.SetErrorVerbosity(errors.Full)

	path := "blueprint_path_here"
	_, err := blueprint.NewBlueprint(path)
	if err != nil {
		err.AppendStackTrace(source)
		return nil, err
	}

	err = errors.NewCustomError(errors.Critical, "Implement NewComposedApplication", source)
	return nil, err
}
