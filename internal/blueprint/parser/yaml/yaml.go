package yaml

import (
	"abdtool/internal/blueprint/parser"
	"abdtool/internal/blueprint/types"
	"abdtool/utils/errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// yamlParser - simple GetBlueprint interface implementator
type yamlParser struct{}

// NewParser - creates and returns a new instance of IBlueprintParser implementator
//
// Returns:
//   - IBlueprintParser: A concrete instance of IBlueprintParser
func NewParser() parser.IBlueprintParser { return &yamlParser{} }

// Private structures designed to facilitate the direct parsing of
// specific YAML file structures into types* structs
type yamlProjectEntry struct {
	Project struct {
		Name    string `yaml:name`
		Version string `yaml:version`
	} `yaml:project`
}

type yamlStepEntry struct {
	Name        string          `yaml:name`
	Description string          `yaml:description`
	Attempts    int             `yaml:"attempts"`
	Timeout_s   int             `yaml:"timeout_s"`
	Post        []yamlPostEntry `yaml:"post"`
	Subscribe   []string        `yaml:"subscribe"`
}

type yamlPostEntry struct {
	Event   string   `yaml:"event"`
	Payload []string `yaml:"payload"`
}

type yamlRoutineEntry struct {
	Name  string `yaml:name`
	Steps []struct {
		Step yamlStepEntry `yaml:"step"`
	} `yaml:"steps"`
	Context map[string]interface{} `yaml:"context"`
}

// parsePostEvents - private function transforms a slice of intermediate event structures (postTmp) into a slice of PostEventEntry types
// each post in postTmp is extracted, processed, and appended to the resulting slice of PostEventEntry types
//
// Parameters:
//   - postTmp: A slice of anonymous structures that directly represent the YAML format
//
// Returns:
//   - A slice of PostEventEntry
func parsePostEvents(postTmp []yamlPostEntry) []types.PostEventEntry {
	var posts []types.PostEventEntry
	for _, post := range postTmp {
		posts = append(posts, types.PostEventEntry{
			Event:   post.Event,
			Payload: post.Payload,
		})
	}
	return posts
}

// parseSteps - private function transforms a slice of intermediate step structures (stepsTmp) into a slice of StepEntry types
// each step in stepsTmp is extracted, processed, and appended to the resulting slice of StepEntry types
//
// Parameters:
//   - stepsTmp: A slice of anonymous structures that directly represent the YAML format
//
// Returns:
//   - A slice of StepEntry
func parseSteps(stepsTmp []struct {
	Step yamlStepEntry `yaml:"step"`
}) []types.StepEntry {
	var steps []types.StepEntry
	for _, step := range stepsTmp {
		steps = append(steps, types.StepEntry{
			Name:        step.Step.Name,
			Description: step.Step.Description,
			Attempts:    step.Step.Attempts,
			Timeout_s:   step.Step.Timeout_s,
			Post:        parsePostEvents(step.Step.Post),
			Subscribe:   step.Step.Subscribe,
		})
	}
	return steps
}

// parseProjectEntry - private function constructs and returns an ProjectEntry instance,
// with project name and it's version
//
// Parameters:
//   - data: blueprint file content
//
// Returns:
//   - ProjectEntry: The constructed ProjectEntry instance or empty struct if there's an error
//   - *errors.CustomError: A custom error that provides detailed information if something went wrong during the construction, nil if there's no error
func parseProjectEntry(data []byte) (types.ProjectEntry, *errors.CustomError) {
	source := "parseProjectEntry"
	tmp := yamlProjectEntry{}

	// Read project yaml entry
	err := yaml.Unmarshal(data, &tmp)
	if err != nil {
		message := fmt.Sprintf("Failed to unmarshall yaml file, %s", err.Error())
		return types.ProjectEntry{}, errors.NewCustomError(errors.Critical, message, source)
	}

	return types.ProjectEntry{
		Name:    tmp.Project.Name,
		Version: tmp.Project.Version,
	}, nil
}

// parseRoutinesEntry - private function constructs and returns slice of RoutineEntry instances
//
// Parameters:
//   - data: blueprint file content
//
// Returns:
//   - []RoutineEntry: The constructed RoutineEntry slice or nil if there's an error
//   - *errors.CustomError: A custom error that provides detailed information if something went wrong during the construction, nil if there's no error
func parseRoutinesEntry(data []byte) ([]types.RoutineEntry, *errors.CustomError) {
	source := "parseRoutinesEntry"

	// Anon struct wrapper for parser
	tmp := struct {
		Routines []struct {
			Routine yamlRoutineEntry `yaml:routine`
		} `yaml:"routines"`
	}{}

	// Read routines yaml entry
	err := yaml.Unmarshal(data, &tmp)
	if err != nil {
		message := fmt.Sprintf("Failed to unmarshall yaml file, %s", err.Error())
		return nil, errors.NewCustomError(errors.Critical, message, source)
	}

	var routines []types.RoutineEntry
	for _, r := range tmp.Routines {
		routines = append(routines, types.RoutineEntry{
			Name:    r.Routine.Name,
			Steps:   parseSteps(r.Routine.Steps),
			Context: r.Routine.Context,
		})
	}

	return routines, nil
}

// GetBlueprint - implemented parser interface, related to .yaml blueprint format,
// returns parsed instance of Blueprint
//
// Parameters:
//   - path: blueprint file path
//
// Returns:
//   - Blueprint: The constructed Blueprint instance or empty struct if there's an error
//   - *errors.CustomError: A custom error that provides detailed information if something went wrong during the construction, nil if there's no error
func (p *yamlParser) GetBlueprint(path string) (types.Blueprint, *errors.CustomError) {
	source := "yaml|GetBlueprint"

	// Check for file exists
	_, e := os.Stat(path)
	if e != nil {
		err := errors.NewCustomError(errors.Critical, "Blueprint file doesn't exist", source)
		return types.Blueprint{}, err
	}

	// Get file extension
	ext := filepath.Ext(path)
	if ext != ".yaml" {
		message := fmt.Sprintf("Invalid blueprint file extention, expected .yaml, got %s", ext)
		return types.Blueprint{}, errors.NewCustomError(errors.Critical, message, source)
	}

	// Read blueprint file content
	data, e := os.ReadFile(path)
	if e != nil {
		message := fmt.Sprintf("Failed to read blueprint file, %s", e.Error())
		return types.Blueprint{}, errors.NewCustomError(errors.Critical, message, source)
	}

	// Retrieve project details
	project, err := parseProjectEntry(data)
	if err != nil {
		err.AppendStackTrace(source)
		return types.Blueprint{}, err
	}

	// Retrieve project details
	routines, err := parseRoutinesEntry(data)
	if err != nil {
		err.AppendStackTrace(source)
		return types.Blueprint{}, err
	}

	return types.Blueprint{
		Project:  project,
		Routines: routines,
	}, nil
}
