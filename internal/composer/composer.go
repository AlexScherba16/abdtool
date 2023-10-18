package composer

import (
	"abdtool/internal/application"
	"abdtool/utils/errors"
)

// NewComposedApplication - constructs and returns an instance of a composed application
// that implements the IApplication interface and ready to lauch
//
// Returns:
//   - application.IApplication: The constructed application instance or nil if there's an error or it's not yet implemented
//   - *errors.CustomError: A custom error that provides detailed information if something went wrong during the construction, nil if there's no error
func NewComposedApplication() (application.IApplication, *errors.CustomError) {
	// Set global error verbosity
	errors.SetErrorVerbosity(errors.Full)

	err := errors.NewCustomError(errors.Critical, "implement NewComposedApplication", "NewComposedApplication")
	return nil, err
}
