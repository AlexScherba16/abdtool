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

	// Using anonymous struct for yaml parser
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

type StepEntry struct {
	// Name - contains specified name of step
	Name string `yaml:"name"`

	// Description - contains more detailed info about step
	Description string `yaml:"description"`

	// Attempts - specifies the number of retry attempts to execute the step in case of failure
	Attempts int `yaml:"attempts"`

	// Timeout_s - time limit to execute step
	Timeout_s int `yaml:"timeout_s"`
}

// RoutineEntry - represents essential details about routine runner. It captures
// the routine's name and set of execution Steps, that logically composed and
// shared context for all steps
type RoutineEntry struct {
	// Name - contains specified name of routine
	Name string

	// Steps - slice of steps related to current routine
	Steps []StepEntry

	// Context - represents a shared data point for all steps
	// It can store key-value pairs to be used across different steps
	Context map[string]interface{}
}

// NewRoutinesEntry - constructs and returns slice of RoutineEntry instances
//
// Returns:
//   - []RoutineEntry: The constructed RoutineEntry slice or nil if there's an error
//   - *errors.CustomError: A custom error that provides detailed information if something went wrong during the construction, nil if there's no error
func NewRoutinesEntry(data []byte) ([]RoutineEntry, *errors.CustomError) {
	source := "NewRoutinesEntry"

	// Using anonymous struct for yaml parser
	tmp := struct {
		Routines []struct {
			Routine struct {
				Name  string `yaml:"name"`
				Steps []struct {
					Step StepEntry `yaml:"step"`
				} `yaml:"steps"`
				Context map[string]interface{} `yaml:"context"`
			} `yaml:"routine"`
		} `yaml:"routines"`
	}{}

	// Read project entry data
	err := yaml.Unmarshal(data, &tmp)
	if err != nil {
		message := fmt.Sprintf("Failed to unmarshall yaml file, %s", err.Error())
		err := errors.NewCustomError(errors.Critical, message, source)
		return nil, err
	}

	// All messy parsing stuff should be here))
	parseSteps := func(stepsTmp []struct {
		Step StepEntry `yaml:"step"`
	}) []StepEntry {
		var steps []StepEntry
		for _, s := range stepsTmp {
			steps = append(steps, s.Step)
		}
		return steps
	}

	parseRoutine := func(routineTmp struct {
		Routine struct {
			Name  string `yaml:"name"`
			Steps []struct {
				Step StepEntry `yaml:"step"`
			} `yaml:"steps"`
			Context map[string]interface{} `yaml:"context"`
		} `yaml:"routine"`
	}) RoutineEntry {
		return RoutineEntry{
			Name:    routineTmp.Routine.Name,
			Steps:   parseSteps(routineTmp.Routine.Steps),
			Context: routineTmp.Routine.Context,
		}
	}

	var routines []RoutineEntry
	for _, r := range tmp.Routines {
		routine := parseRoutine(r)

		// Validate routine fields
		if routine.Name == "" {
			message := fmt.Sprintf("%q must be provided", "routine.name")
			return nil, errors.NewCustomError(errors.Critical, message, source)
		}
		if len(routine.Steps) == 0 {
			message := fmt.Sprintf("%q must be provided", "routine.steps")
			return nil, errors.NewCustomError(errors.Critical, message, source)
		}

		// Validate step fields
		for _, step := range routine.Steps {
			if step.Name == "" {
				message := fmt.Sprintf("%q must be provided", "step.name")
				return nil, errors.NewCustomError(errors.Critical, message, source)
			}
			if step.Description == "" {
				message := fmt.Sprintf("%q must be provided", "step.description")
				return nil, errors.NewCustomError(errors.Critical, message, source)
			}
			if step.Timeout_s == 0 {
				message := fmt.Sprintf("%q must be provided", "step.timeout_s")
				return nil, errors.NewCustomError(errors.Critical, message, source)
			}
			if step.Attempts == 0 {
				step.Attempts = 1
			}
		}
		routines = append(routines, routine)
	}

	return routines, nil
}
