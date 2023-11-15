package blueprint

import (
	"abdtool/internal/blueprint/parser"
	"abdtool/internal/blueprint/parser/yaml"
	"abdtool/internal/blueprint/types"
	"abdtool/internal/misc"
	"abdtool/utils/errors"
	"fmt"
	"reflect"
	"testing"
)

func Test_validateBlueprint(t *testing.T) {
	errorSource := "validateBlueprint"
	routineName := "TEST_ROUTINE_NAME"
	stepName := "TEST_STEP_NAME"
	tests := []struct {
		name          string
		blueprint     *types.Blueprint
		expectedError *errors.CustomError
	}{
		{
			name:          "emptyProjectName",
			blueprint:     &types.Blueprint{},
			expectedError: errors.NewCustomError(errors.Critical, "\"project.name\" must be provided", errorSource),
		},
		{
			name:          "emptyProjectVersion",
			blueprint:     &types.Blueprint{Project: types.ProjectEntry{Name: "t"}},
			expectedError: errors.NewCustomError(errors.Critical, "\"project.version\" must be provided", errorSource),
		},
		{
			name:          "emptyRoutinesEntry",
			blueprint:     &types.Blueprint{Project: types.ProjectEntry{Name: "t", Version: "t"}},
			expectedError: errors.NewCustomError(errors.Critical, "\"routines\" must be provided", errorSource),
		},
		{
			name:          "emptyRoutineName",
			blueprint:     &types.Blueprint{Project: types.ProjectEntry{Name: "t", Version: "t"}, Routines: []types.RoutineEntry{{}}},
			expectedError: errors.NewCustomError(errors.Critical, "\"routine.name\" must be provided", errorSource),
		},
		{
			name: "emptyRoutineStepsEntry",
			blueprint: &types.Blueprint{Project: types.ProjectEntry{Name: "t", Version: "t"}, Routines: []types.RoutineEntry{{
				Name: routineName,
			}}},
			expectedError: errors.NewCustomError(errors.Critical, fmt.Sprintf("%q must be provided for %q routine", "routine.steps", routineName), errorSource),
		},
		{
			name: "emptyStepName",
			blueprint: &types.Blueprint{Project: types.ProjectEntry{Name: "t", Version: "t"}, Routines: []types.RoutineEntry{{
				Name: routineName, Steps: []types.StepEntry{{}},
			}}},
			expectedError: errors.NewCustomError(errors.Critical, fmt.Sprintf("%q must be provided in %q routine", "step.name", routineName), errorSource),
		},
		{
			name: "emptyStepDescription",
			blueprint: &types.Blueprint{Project: types.ProjectEntry{Name: "t", Version: "t"}, Routines: []types.RoutineEntry{{
				Name: routineName, Steps: []types.StepEntry{{Name: stepName}},
			}}},
			expectedError: errors.NewCustomError(errors.Critical, fmt.Sprintf("%q must be provided for %q step", "step.description", stepName), errorSource),
		},
		{
			name: "emptyStepTimeout",
			blueprint: &types.Blueprint{Project: types.ProjectEntry{Name: "t", Version: "t"}, Routines: []types.RoutineEntry{{
				Name: routineName, Steps: []types.StepEntry{{Name: stepName, Description: "t"}},
			}}},
			expectedError: errors.NewCustomError(errors.Critical, fmt.Sprintf("%q must be >= 0 in %q step", "step.timeout_s", stepName), errorSource),
		},
		{
			name: "negativeStepTimeout",
			blueprint: &types.Blueprint{Project: types.ProjectEntry{Name: "t", Version: "t"}, Routines: []types.RoutineEntry{{
				Name: routineName, Steps: []types.StepEntry{{Name: stepName, Description: "t", Timeout_s: -1}},
			}}},
			expectedError: errors.NewCustomError(errors.Critical, fmt.Sprintf("%q must be >= 0 in %q step", "step.timeout_s", stepName), errorSource),
		},
		{
			name: "emptyStepAttempts",
			blueprint: &types.Blueprint{Project: types.ProjectEntry{Name: "t", Version: "t"}, Routines: []types.RoutineEntry{{
				Name: routineName, Steps: []types.StepEntry{{Name: stepName, Description: "t", Timeout_s: 1}},
			}}},
			expectedError: errors.NewCustomError(errors.Critical, fmt.Sprintf("%q must be >= 0 in %q step", "step.attempts", stepName), errorSource),
		},
		{
			name: "negativeStepAttempts",
			blueprint: &types.Blueprint{Project: types.ProjectEntry{Name: "t", Version: "t"}, Routines: []types.RoutineEntry{{
				Name: routineName, Steps: []types.StepEntry{{Name: stepName, Description: "t", Timeout_s: 1, Attempts: -1}},
			}}},
			expectedError: errors.NewCustomError(errors.Critical, fmt.Sprintf("%q must be >= 0 in %q step", "step.attempts", stepName), errorSource),
		},
		{
			name: "validBlueprint",
			blueprint: &types.Blueprint{Project: types.ProjectEntry{Name: "t", Version: "t"}, Routines: []types.RoutineEntry{{
				Name: routineName, Steps: []types.StepEntry{{Name: stepName, Description: "t", Timeout_s: 1, Attempts: 1}},
			}}},
			expectedError: nil,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			/* ARRANGE */
			/* ACT */
			err := validateBlueprint(testCase.blueprint)

			/* ASSERT */
			// Assert expected error
			if testCase.expectedError != nil {
				if err == nil {
					t.Fatalf("[%s] : expected error -> %v, got -> %v", t.Name(), testCase.expectedError, nil)
				}

				// Assert expected error severity
				if err.Severity() != testCase.expectedError.Severity() {
					t.Fatalf("[%s] : expected severity -> %v, got -> %v", t.Name(), testCase.expectedError.Severity(), err.Severity())
				}

				// Assert expected error message
				if err.Error() != testCase.expectedError.Error() {
					t.Fatalf("[%s] : expected error message -> %v, got -> %v", t.Name(), testCase.expectedError.Error(), err.Error())
				}

			} else {
				// Assert unexpected error
				if err != nil {
					t.Fatalf("[%s] : unexpected error -> %v, got -> %v", t.Name(), nil, err.Error())
				}
			}
		})
	}
}

func Test_buildParser(t *testing.T) {
	tests := []struct {
		name          string
		path          string
		expectedError *errors.CustomError
		parser        parser.IBlueprintParser
	}{
		{
			name:          "invalidPath",
			path:          "qwe.path",
			expectedError: errors.NewCustomError(errors.Critical, fmt.Sprintf("Invalid blueprint file extention, %q is not implemented", ".path"), "buildParser"),
			parser:        nil,
		},
		{
			name:          "invalidYamlPath",
			path:          "qwe.yaml",
			expectedError: nil,
			parser:        yaml.NewParser(),
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			/* ARRANGE */
			/* ACT */
			parser, err := buildParser(testCase.path)

			/* ASSERT */
			// Assert expected error
			if testCase.expectedError != nil {
				if err == nil {
					t.Fatalf("[%s] : expected error -> %v, got -> %v", t.Name(), testCase.expectedError, nil)
				}

				// Assert expected error severity
				if err.Severity() != testCase.expectedError.Severity() {
					t.Fatalf("[%s] : expected severity -> %v, got -> %v", t.Name(), testCase.expectedError.Severity(), err.Severity())
				}

				// Assert expected error message
				if err.Error() != testCase.expectedError.Error() {
					t.Fatalf("[%s] : expected error message -> %v, got -> %v", t.Name(), testCase.expectedError.Error(), err.Error())
				}

				// Assert expected nil parser
				if parser != nil {
					t.Fatalf("[%s] : expected parser -> %v, got -> %v", t.Name(), nil, parser)
				}

			} else {
				// Assert unexpected error
				if err != nil {
					t.Fatalf("[%s] : unexpected error -> %v, got -> %v", t.Name(), nil, err.Error())
				}
			}
		})
	}
}

func Test_NewBlueprint_InvalidPath(t *testing.T) {
	path := "qwe.path"
	message := fmt.Sprintf("Invalid blueprint file extention, %q is not implemented", ".path")
	source := "buildParser"

	tests := []struct {
		name           string
		verbosity      errors.VerbosityLevel
		expectedError  *errors.CustomError
		path           string
		expectedResult types.Blueprint
	}{
		{
			name:           "errorVerbosity",
			verbosity:      errors.ErrorMessage,
			expectedError:  errors.NewCustomError(errors.Critical, message, source),
			path:           path,
			expectedResult: types.Blueprint{},
		},
		{
			name:           "stackTraceVerbosity",
			verbosity:      errors.StackTrace,
			expectedError:  misc.CreateErrorWithStackTrace(errors.Critical, message, source, "NewBlueprint"),
			path:           path,
			expectedResult: types.Blueprint{},
		},
		{
			name:           "fullVerbosity",
			verbosity:      errors.Full,
			expectedError:  misc.CreateErrorWithStackTrace(errors.Critical, message, source, "NewBlueprint"),
			path:           path,
			expectedResult: types.Blueprint{},
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			/* ARRANGE */
			errors.SetErrorVerbosity(testCase.verbosity)

			/* ACT */
			result, err := NewBlueprint(testCase.path)

			/* ASSERT */
			// Assert expected error
			if err == nil {
				t.Fatalf("[%s] : expected error -> %v, got -> %v", t.Name(), testCase.expectedError, nil)
			}

			// Assert expected error severity
			if err.Severity() != testCase.expectedError.Severity() {
				t.Fatalf("[%s] : expected severity -> %v, got -> %v", t.Name(), testCase.expectedError.Severity(), err.Severity())
			}

			// Assert expected error message
			if err.Error() != testCase.expectedError.Error() {
				t.Fatalf("[%s] : expected error message -> %v, got -> %v", t.Name(), testCase.expectedError.Error(), err.Error())
			}

			// Assert expected result
			if !reflect.DeepEqual(result, testCase.expectedResult) {
				t.Fatalf("[%s] : expected Blueprint -> %v, got -> %v", t.Name(), testCase.expectedResult, result)
			}
		})
	}
}

func Test_NewBlueprint_ParserError(t *testing.T) {
	path := misc.ValidYamlFile
	message := fmt.Sprintf("Failed to read blueprint file, open %s: The filename, directory name, or volume label syntax is incorrect.", path)
	source := "yaml|GetBlueprint"

	tests := []struct {
		name           string
		verbosity      errors.VerbosityLevel
		expectedError  *errors.CustomError
		path           string
		expectedResult types.Blueprint
	}{
		{
			name:           "errorVerbosity",
			verbosity:      errors.ErrorMessage,
			expectedError:  errors.NewCustomError(errors.Critical, message, source),
			path:           path,
			expectedResult: types.Blueprint{},
		},
		{
			name:           "stackTraceVerbosity",
			verbosity:      errors.StackTrace,
			expectedError:  misc.CreateErrorWithStackTrace(errors.Critical, message, source, "NewBlueprint"),
			path:           path,
			expectedResult: types.Blueprint{},
		},
		{
			name:           "fullVerbosity",
			verbosity:      errors.Full,
			expectedError:  misc.CreateErrorWithStackTrace(errors.Critical, message, source, "NewBlueprint"),
			path:           path,
			expectedResult: types.Blueprint{},
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			/* ARRANGE */
			errors.SetErrorVerbosity(testCase.verbosity)

			/* ACT */
			result, err := NewBlueprint(testCase.path)

			/* ASSERT */
			// Assert expected error
			if err == nil {
				t.Fatalf("[%s] : expected error -> %v, got -> %v", t.Name(), testCase.expectedError, nil)
			}

			// Assert expected error severity
			if err.Severity() != testCase.expectedError.Severity() {
				t.Fatalf("[%s] : expected severity -> %v, got -> %v", t.Name(), testCase.expectedError.Severity(), err.Severity())
			}

			// Assert expected error message
			if err.Error() != testCase.expectedError.Error() {
				t.Fatalf("[%s] : expected error message -> %v, got -> %v", t.Name(), testCase.expectedError.Error(), err.Error())
			}

			// Assert expected result
			if !reflect.DeepEqual(result, testCase.expectedResult) {
				t.Fatalf("[%s] : expected Blueprint -> %v, got -> %v", t.Name(), testCase.expectedResult, result)
			}
		})
	}
}

func Test_NewBlueprint_BlueprintValidationError(t *testing.T) {
	message := fmt.Sprintf("%q must be provided", "project.name")
	source := "validateBlueprint"

	tests := []struct {
		name           string
		verbosity      errors.VerbosityLevel
		expectedError  *errors.CustomError
		yaml           map[string]interface{}
		expectedResult types.Blueprint
	}{
		{
			name:           "errorVerbosity",
			verbosity:      errors.ErrorMessage,
			expectedError:  errors.NewCustomError(errors.Critical, message, source),
			yaml:           map[string]interface{}{},
			expectedResult: types.Blueprint{},
		},
		{
			name:           "stackTraceVerbosity",
			verbosity:      errors.StackTrace,
			expectedError:  misc.CreateErrorWithStackTrace(errors.Critical, message, source, "NewBlueprint"),
			yaml:           map[string]interface{}{},
			expectedResult: types.Blueprint{},
		},
		{
			name:           "fullVerbosity",
			verbosity:      errors.Full,
			expectedError:  misc.CreateErrorWithStackTrace(errors.Critical, message, source, "NewBlueprint"),
			yaml:           map[string]interface{}{},
			expectedResult: types.Blueprint{},
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			/* ARRANGE */
			errors.SetErrorVerbosity(testCase.verbosity)
			path, e := misc.CreateTmpYAML(misc.ValidYamlFile, testCase.yaml)
			if e != nil {
				t.Fatalf("[%s] : can't create YAML file, %s", t.Name(), e.Error())
			}
			defer misc.DeleteTmpYAML(path)

			/* ACT */
			result, err := NewBlueprint(path)

			/* ASSERT */
			// Assert expected error
			if err == nil {
				t.Fatalf("[%s] : expected error -> %v, got -> %v", t.Name(), testCase.expectedError, nil)
			}

			// Assert expected error severity
			if err.Severity() != testCase.expectedError.Severity() {
				t.Fatalf("[%s] : expected severity -> %v, got -> %v", t.Name(), testCase.expectedError.Severity(), err.Severity())
			}

			// Assert expected error message
			if err.Error() != testCase.expectedError.Error() {
				t.Fatalf("[%s] : expected error message -> %v, got -> %v", t.Name(), testCase.expectedError.Error(), err.Error())
			}

			// Assert expected result
			if !reflect.DeepEqual(result, testCase.expectedResult) {
				t.Fatalf("[%s] : expected Blueprint -> %v, got -> %v", t.Name(), testCase.expectedResult, result)
			}
		})
	}
}

func Test_NewBlueprint_ValidBlueprintFile(t *testing.T) {
	/* ARRANGE */
	projectName := "test"
	projectVersion := "2.2.8"
	projectEntry := map[string]interface{}{"name": projectName, "version": projectVersion}

	stepEntry := map[string]interface{}{
		"name":        "TEST_STEP_0",
		"description": "TEST_STEP_0_DESCRIPTION",
		"attempts":    100,
		"timeout_s":   1000,
		"post": []interface{}{map[string]interface{}{
			"event":   "TEST_STEP_0_POST_EVENT",
			"payload": []string{"PAYLOAD_0", "PAYLOAD_1"},
		}},
	}
	routineSteps := []interface{}{map[string]interface{}{"step": stepEntry}}
	routineEntry := map[string]interface{}{"name": "TEST_ROUTINE_0", "steps": routineSteps}
	routinesEntry := []interface{}{map[string]interface{}{"routine": routineEntry}}
	yamlContent := map[string]interface{}{"project": projectEntry, "routines": routinesEntry}

	expectedResult := types.Blueprint{
		Project: types.ProjectEntry{Name: projectName, Version: projectVersion},
		Routines: []types.RoutineEntry{
			{Name: "TEST_ROUTINE_0", Steps: []types.StepEntry{
				{
					Name:        "TEST_STEP_0",
					Description: "TEST_STEP_0_DESCRIPTION",
					Attempts:    100,
					Timeout_s:   1000,
					Post:        []types.PostEventEntry{{Event: "TEST_STEP_0_POST_EVENT", Payload: []string{"PAYLOAD_0", "PAYLOAD_1"}}},
				},
			}},
		},
	}

	path, e := misc.CreateTmpYAML(misc.ValidYamlFile, yamlContent)
	if e != nil {
		t.Fatalf("[%s] : can't create YAML file, %s", t.Name(), e.Error())
	}
	defer misc.DeleteTmpYAML(path)

	/* ACT */
	result, err := NewBlueprint(path)

	/* ASSERT */
	// Assert unexpected error
	if err != nil {
		t.Fatalf("[%s] : expected error -> %v, got -> %v", t.Name(), nil, err)
	}

	// Assert expected result
	if !reflect.DeepEqual(result, expectedResult) {
		t.Fatalf("[%s] : expected Blueprint -> %v, got -> %v", t.Name(), expectedResult, result)
	}
}
