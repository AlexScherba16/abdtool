package parser

import (
	"abdtool/internal/blueprint/types"
	"abdtool/utils/errors"
)

// IBlueprintParser - defines the interface for blueprint file parsers
// implementing this interface should provide functionality specific to each file format
type IBlueprintParser interface {
	// GetBlueprint - is a method of interface which is expected to return an instance
	// of Blueprint type or error if there were issues during the file parsing process
	//
	// Parameters:
	//  - path: blueprint file path
	//
	// Returns:
	//  - Blueprint: The constructed Blueprint instance or empty struct if there's an error
	//  - *errors.CustomError: A custom error that provides detailed information
	//      if something went wrong during the construction, nil if there's no error
	GetBlueprint(path string) (types.Blueprint, *errors.CustomError)
}
