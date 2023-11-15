package yaml

import (
	"abdtool/internal/blueprint/types"
	"abdtool/internal/misc"
	"abdtool/utils/errors"
	"fmt"
	"os"
	"reflect"
	"testing"
)

// common global variables
const expectedName string = "test"
const expectedVersion string = "2.2.8"

func Test_parseProjectEntry(t *testing.T) {
	tests := []struct {
		name           string
		yaml           map[string]interface{}
		expectedError  *errors.CustomError
		expectedResult types.ProjectEntry
	}{
		{
			name:           "noProjectEntry",
			yaml:           map[string]interface{}{},
			expectedError:  nil,
			expectedResult: types.ProjectEntry{},
		},
		{
			name:           "emptyProjectEntry",
			yaml:           map[string]interface{}{"project": map[string]interface{}{}},
			expectedError:  nil,
			expectedResult: types.ProjectEntry{},
		},
		{
			name:           "emptyProjectName",
			yaml:           map[string]interface{}{"project": map[string]interface{}{"name": nil}},
			expectedError:  nil,
			expectedResult: types.ProjectEntry{},
		},
		{
			name:          "emptyProjectVersion",
			yaml:          map[string]interface{}{"project": map[string]interface{}{"name": expectedName}},
			expectedError: nil,
			expectedResult: types.ProjectEntry{
				Name: expectedName,
			},
		},
		{
			name:           "emptyNameValidVersion",
			yaml:           map[string]interface{}{"project": map[string]interface{}{"version": expectedVersion}},
			expectedError:  nil,
			expectedResult: types.ProjectEntry{Version: expectedVersion},
		},
		{
			name: "validProjectEntry",
			yaml: map[string]interface{}{
				"project": map[string]interface{}{
					"name":    expectedName,
					"version": expectedVersion,
				},
			},
			expectedError: nil,
			expectedResult: types.ProjectEntry{
				Name:    expectedName,
				Version: expectedVersion,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			/* ARRANGE */
			// Mock blueprint file
			path, e := misc.CreateTmpYAML(misc.ValidYamlFile, test.yaml)
			if e != nil {
				t.Fatalf("[%s] : can't create YAML file, %s", t.Name(), e.Error())
			}
			defer misc.DeleteTmpYAML(path)

			data, e := os.ReadFile(path)
			if e != nil {
				t.Fatalf("[%s] : failed to read blueprint file, %s", t.Name(), e.Error())
			}

			/* ACT */
			result, err := parseProjectEntry(data)
			/* ASSERT */
			// Assert expected error
			if test.expectedError != nil {
				if err == nil {
					t.Fatalf("[%s] : expected error -> %v, got -> %v", t.Name(), test.expectedError, nil)
				}

				// Assert expected error severity
				if err.Severity() != test.expectedError.Severity() {
					t.Fatalf("[%s] : expected severity -> %v, got -> %v", t.Name(), test.expectedError.Severity(), err.Severity())
				}

				// Assert expected error message
				if err.Error() != test.expectedError.Error() {
					t.Fatalf("[%s] : expected error message -> %v, got -> %v", t.Name(), test.expectedError.Error(), err.Error())
				}
			} else {
				// Assert unexpected error
				if err != nil {
					t.Fatalf("[%s] : unexpected error -> %v, got -> %v", t.Name(), nil, err.Error())
				}
			}

			// Assert expected result
			eq := reflect.DeepEqual(result, test.expectedResult)
			if !eq {
				t.Fatalf("[%s] : expected project entry -> %v, got -> %v", t.Name(), test.expectedResult, result)
			}
		})
	}
}

func Test_parseRoutinesEntry(t *testing.T) {
	tests := []struct {
		name           string
		yaml           map[string]interface{}
		expectedError  *errors.CustomError
		expectedResult []types.RoutineEntry
	}{
		{
			name:           "noRoutinesEntry",
			yaml:           map[string]interface{}{},
			expectedError:  nil,
			expectedResult: []types.RoutineEntry{},
		},
		{
			name:           "emptyRoutinesEntry",
			yaml:           map[string]interface{}{"routines": []interface{}{}},
			expectedError:  nil,
			expectedResult: []types.RoutineEntry{},
		},
		{
			name: "singleEmptyRoutineEntry",
			yaml: map[string]interface{}{
				"routines": []interface{}{
					map[string]interface{}{"routine": map[string]interface{}{}},
				},
			},
			expectedError:  nil,
			expectedResult: []types.RoutineEntry{{}},
		},
		{
			name: "singleRoutineEntry",
			yaml: map[string]interface{}{
				"routines": []interface{}{map[string]interface{}{
					"routine": map[string]interface{}{
						"name": "TEST_ROUTINE_0",
					}}},
			},
			expectedError:  nil,
			expectedResult: []types.RoutineEntry{{Name: "TEST_ROUTINE_0"}},
		},
		{
			name: "singleRoutineEntryEmptySteps",
			yaml: map[string]interface{}{
				"routines": []interface{}{map[string]interface{}{
					"routine": map[string]interface{}{
						"name":  "TEST_ROUTINE_0",
						"steps": []interface{}{},
					}}},
			},
			expectedError:  nil,
			expectedResult: []types.RoutineEntry{{Name: "TEST_ROUTINE_0"}},
		},
		{
			name: "singleRoutineSingleEmptyStep",
			yaml: map[string]interface{}{
				"routines": []interface{}{map[string]interface{}{
					"routine": map[string]interface{}{
						"name": "TEST_ROUTINE_0",
						"steps": []interface{}{map[string]interface{}{
							"step": map[string]interface{}{}}},
					}}},
			},
			expectedError: nil,
			expectedResult: []types.RoutineEntry{{
				Name:  "TEST_ROUTINE_0",
				Steps: []types.StepEntry{{}},
			}},
		},
		{
			name: "singleRoutineSingleStepEmptyDescription",
			yaml: map[string]interface{}{
				"routines": []interface{}{map[string]interface{}{
					"routine": map[string]interface{}{
						"name": "TEST_ROUTINE_0",
						"steps": []interface{}{map[string]interface{}{
							"step": map[string]interface{}{
								"name": "TEST_STEP_0",
							}}},
					}}},
			},
			expectedError: nil,
			expectedResult: []types.RoutineEntry{{
				Name:  "TEST_ROUTINE_0",
				Steps: []types.StepEntry{{Name: "TEST_STEP_0"}},
			}},
		},
		{
			name: "singleRoutineSingleStepEmptyAttempts",
			yaml: map[string]interface{}{
				"routines": []interface{}{map[string]interface{}{
					"routine": map[string]interface{}{
						"name": "TEST_ROUTINE_0",
						"steps": []interface{}{map[string]interface{}{
							"step": map[string]interface{}{
								"name":        "TEST_STEP_0",
								"description": "TEST_STEP_0_DESCRIPTION",
							}}},
					}}},
			},
			expectedError: nil,
			expectedResult: []types.RoutineEntry{{
				Name:  "TEST_ROUTINE_0",
				Steps: []types.StepEntry{{Name: "TEST_STEP_0", Description: "TEST_STEP_0_DESCRIPTION"}},
			}},
		},
		{
			name: "singleRoutineSingleStepEmptyTimeout",
			yaml: map[string]interface{}{
				"routines": []interface{}{map[string]interface{}{
					"routine": map[string]interface{}{
						"name": "TEST_ROUTINE_0",
						"steps": []interface{}{map[string]interface{}{
							"step": map[string]interface{}{
								"name":        "TEST_STEP_0",
								"description": "TEST_STEP_0_DESCRIPTION",
								"attempts":    123,
							}}},
					}}},
			},
			expectedError: nil,
			expectedResult: []types.RoutineEntry{{
				Name:  "TEST_ROUTINE_0",
				Steps: []types.StepEntry{{Name: "TEST_STEP_0", Description: "TEST_STEP_0_DESCRIPTION", Attempts: 123}},
			}},
		},
		{
			name: "singleRoutineSingleStepEmptyPost",
			yaml: map[string]interface{}{
				"routines": []interface{}{map[string]interface{}{
					"routine": map[string]interface{}{
						"name": "TEST_ROUTINE_0",
						"steps": []interface{}{map[string]interface{}{
							"step": map[string]interface{}{
								"name":        "TEST_STEP_0",
								"description": "TEST_STEP_0_DESCRIPTION",
								"attempts":    228,
								"timeout_s":   1337,
							}}},
					}}},
			},
			expectedError: nil,
			expectedResult: []types.RoutineEntry{{
				Name:  "TEST_ROUTINE_0",
				Steps: []types.StepEntry{{Name: "TEST_STEP_0", Description: "TEST_STEP_0_DESCRIPTION", Attempts: 228, Timeout_s: 1337}},
			}},
		},
		{
			name: "singleRoutineSingleStepEmptyPostEntry",
			yaml: map[string]interface{}{
				"routines": []interface{}{map[string]interface{}{
					"routine": map[string]interface{}{
						"name": "TEST_ROUTINE_0",
						"steps": []interface{}{map[string]interface{}{
							"step": map[string]interface{}{
								"name":        "TEST_STEP_0",
								"description": "TEST_STEP_0_DESCRIPTION",
								"attempts":    228,
								"timeout_s":   1337,
								"post":        []interface{}{},
							}}},
					}}},
			},
			expectedError: nil,
			expectedResult: []types.RoutineEntry{{
				Name:  "TEST_ROUTINE_0",
				Steps: []types.StepEntry{{Name: "TEST_STEP_0", Description: "TEST_STEP_0_DESCRIPTION", Attempts: 228, Timeout_s: 1337}},
			}},
		},
		{
			name: "singleRoutineSingleStepSinglePostEmptyPayload",
			yaml: map[string]interface{}{
				"routines": []interface{}{map[string]interface{}{
					"routine": map[string]interface{}{
						"name": "TEST_ROUTINE_0",
						"steps": []interface{}{map[string]interface{}{
							"step": map[string]interface{}{
								"name":        "TEST_STEP_0",
								"description": "TEST_STEP_0_DESCRIPTION",
								"attempts":    228,
								"timeout_s":   1337,
								"post": []interface{}{map[string]interface{}{
									"event": "TEST_STEP_0_POST_EVENT",
								}},
							}}},
					}}},
			},
			expectedError: nil,
			expectedResult: []types.RoutineEntry{{
				Name: "TEST_ROUTINE_0",
				Steps: []types.StepEntry{
					{
						Name: "TEST_STEP_0", Description: "TEST_STEP_0_DESCRIPTION", Attempts: 228, Timeout_s: 1337,
						Post: []types.PostEventEntry{{Event: "TEST_STEP_0_POST_EVENT"}}}},
			}},
		},
		{
			name: "singleRoutineSingleStepSinglePostEmptyPayloadEntry",
			yaml: map[string]interface{}{
				"routines": []interface{}{map[string]interface{}{
					"routine": map[string]interface{}{
						"name": "TEST_ROUTINE_0",
						"steps": []interface{}{map[string]interface{}{
							"step": map[string]interface{}{
								"name":        "TEST_STEP_0",
								"description": "TEST_STEP_0_DESCRIPTION",
								"attempts":    228,
								"timeout_s":   1337,
								"post": []interface{}{map[string]interface{}{
									"event":   "TEST_STEP_0_POST_EVENT",
									"payload": []string{},
								}},
							}}},
					}}},
			},
			expectedError: nil,
			expectedResult: []types.RoutineEntry{{
				Name: "TEST_ROUTINE_0",
				Steps: []types.StepEntry{
					{
						Name: "TEST_STEP_0", Description: "TEST_STEP_0_DESCRIPTION", Attempts: 228, Timeout_s: 1337,
						Post: []types.PostEventEntry{{Event: "TEST_STEP_0_POST_EVENT", Payload: []string{}}}}},
			}},
		},
		{
			name: "singleRoutineSingleStepSinglePostPayloadEntry",
			yaml: map[string]interface{}{
				"routines": []interface{}{map[string]interface{}{
					"routine": map[string]interface{}{
						"name": "TEST_ROUTINE_0",
						"steps": []interface{}{map[string]interface{}{
							"step": map[string]interface{}{
								"name":        "TEST_STEP_0",
								"description": "TEST_STEP_0_DESCRIPTION",
								"attempts":    228,
								"timeout_s":   1337,
								"post": []interface{}{map[string]interface{}{
									"event":   "TEST_STEP_0_POST_EVENT",
									"payload": []string{"PAYLOAD_0", "PAYLOAD_1"},
								}},
							}}},
					}}},
			},
			expectedError: nil,
			expectedResult: []types.RoutineEntry{{
				Name: "TEST_ROUTINE_0",
				Steps: []types.StepEntry{
					{
						Name: "TEST_STEP_0", Description: "TEST_STEP_0_DESCRIPTION", Attempts: 228, Timeout_s: 1337,
						Post: []types.PostEventEntry{{Event: "TEST_STEP_0_POST_EVENT", Payload: []string{"PAYLOAD_0", "PAYLOAD_1"}}}}},
			}},
		},
		{
			name: "singleRoutineMultipleStepsMultiplePostPayloadEntries",
			yaml: map[string]interface{}{
				"routines": []interface{}{map[string]interface{}{
					"routine": map[string]interface{}{
						"name": "TEST_ROUTINE_0",
						"steps": []interface{}{
							map[string]interface{}{
								"step": map[string]interface{}{
									"name":        "TEST_STEP_0",
									"description": "TEST_STEP_0_DESCRIPTION",
									"attempts":    100,
									"timeout_s":   1000,
									"post": []interface{}{map[string]interface{}{
										"event":   "TEST_STEP_0_POST_EVENT",
										"payload": []string{"PAYLOAD_0", "PAYLOAD_1"},
									}},
								}},
							map[string]interface{}{
								"step": map[string]interface{}{
									"name":        "TEST_STEP_1",
									"description": "TEST_STEP_1_DESCRIPTION",
									"attempts":    200,
									"timeout_s":   2000,
									"post": []interface{}{map[string]interface{}{
										"event":   "TEST_STEP_1_POST_EVENT",
										"payload": []string{"PAYLOAD_0", "PAYLOAD_1", "PAYLOAD_2"},
									}},
								}}},
					}},
				},
			},
			expectedError: nil,
			expectedResult: []types.RoutineEntry{{
				Name: "TEST_ROUTINE_0",
				Steps: []types.StepEntry{{
					Name: "TEST_STEP_0", Description: "TEST_STEP_0_DESCRIPTION", Attempts: 100, Timeout_s: 1000,
					Post: []types.PostEventEntry{
						{Event: "TEST_STEP_0_POST_EVENT", Payload: []string{"PAYLOAD_0", "PAYLOAD_1"}},
					},
				}, {
					Name: "TEST_STEP_1", Description: "TEST_STEP_1_DESCRIPTION", Attempts: 200, Timeout_s: 2000,
					Post: []types.PostEventEntry{
						{Event: "TEST_STEP_1_POST_EVENT", Payload: []string{"PAYLOAD_0", "PAYLOAD_1", "PAYLOAD_2"}}},
				}},
			}},
		},
		{
			name: "multipleRoutinesMultipleStepsMultiplePostPayloadEntries",
			yaml: map[string]interface{}{
				"routines": []interface{}{map[string]interface{}{
					"routine": map[string]interface{}{
						"name": "TEST_ROUTINE_0",
						"steps": []interface{}{map[string]interface{}{
							"step": map[string]interface{}{
								"name":        "TEST_STEP_0",
								"description": "TEST_STEP_0_DESCRIPTION",
								"attempts":    100,
								"timeout_s":   1000,
								"post": []interface{}{map[string]interface{}{
									"event":   "TEST_STEP_0_POST_EVENT",
									"payload": []string{"PAYLOAD_0", "PAYLOAD_1"},
								}},
							}}, map[string]interface{}{
							"step": map[string]interface{}{
								"name":        "TEST_STEP_1",
								"description": "TEST_STEP_1_DESCRIPTION",
								"attempts":    200,
								"timeout_s":   2000,
								"post": []interface{}{map[string]interface{}{
									"event":   "TEST_STEP_1_POST_EVENT",
									"payload": []string{"PAYLOAD_0", "PAYLOAD_1", "PAYLOAD_2"},
								}},
							}}},
					}}, map[string]interface{}{
					"routine": map[string]interface{}{
						"name": "TEST_ROUTINE_1",
						"steps": []interface{}{map[string]interface{}{
							"step": map[string]interface{}{
								"name":        "TEST_STEP_2",
								"description": "TEST_STEP_2_DESCRIPTION",
								"attempts":    300,
								"timeout_s":   3000,
								"post": []interface{}{map[string]interface{}{
									"event":   "TEST_STEP_2_POST_EVENT",
									"payload": []string{"PAYLOAD_2", "PAYLOAD_5"},
								}},
							}}},
					}},
				},
			},
			expectedError: nil,
			expectedResult: []types.RoutineEntry{{
				Name: "TEST_ROUTINE_0",
				Steps: []types.StepEntry{{
					Name: "TEST_STEP_0", Description: "TEST_STEP_0_DESCRIPTION", Attempts: 100, Timeout_s: 1000,
					Post: []types.PostEventEntry{
						{Event: "TEST_STEP_0_POST_EVENT", Payload: []string{"PAYLOAD_0", "PAYLOAD_1"}},
					},
				}, {
					Name: "TEST_STEP_1", Description: "TEST_STEP_1_DESCRIPTION", Attempts: 200, Timeout_s: 2000,
					Post: []types.PostEventEntry{
						{Event: "TEST_STEP_1_POST_EVENT", Payload: []string{"PAYLOAD_0", "PAYLOAD_1", "PAYLOAD_2"}}},
				}},
			}, {
				Name: "TEST_ROUTINE_1",
				Steps: []types.StepEntry{{
					Name: "TEST_STEP_2", Description: "TEST_STEP_2_DESCRIPTION", Attempts: 300, Timeout_s: 3000,
					Post: []types.PostEventEntry{
						{Event: "TEST_STEP_2_POST_EVENT", Payload: []string{"PAYLOAD_2", "PAYLOAD_5"}},
					},
				}},
			}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			/* ARRANGE */
			// Mock blueprint file
			path, e := misc.CreateTmpYAML(misc.ValidYamlFile, test.yaml)
			if e != nil {
				t.Fatalf("[%s] : can't create YAML file, %s", t.Name(), e.Error())
			}
			defer misc.DeleteTmpYAML(path)

			data, e := os.ReadFile(path)
			if e != nil {
				t.Fatalf("[%s] : failed to read blueprint file, %s", t.Name(), e.Error())
			}

			/* ACT */
			result, err := parseRoutinesEntry(data)
			/* ASSERT */
			// Assert expected error
			if test.expectedError != nil {
				if err == nil {
					t.Fatalf("[%s] : expected error -> %v, got -> %v", t.Name(), test.expectedError, nil)
				}

				// Assert expected error severity
				if err.Severity() != test.expectedError.Severity() {
					t.Fatalf("[%s] : expected severity -> %v, got -> %v", t.Name(), test.expectedError.Severity(), err.Severity())
				}

				// Assert expected error message
				if err.Error() != test.expectedError.Error() {
					t.Fatalf("[%s] : expected error message -> %v, got -> %v", t.Name(), test.expectedError.Error(), err.Error())
				}
			} else {
				// Assert unexpected error
				if err != nil {
					t.Fatalf("[%s] : unexpected error -> %v, got -> %v", t.Name(), nil, err.Error())
				}
			}

			// Assert expected len
			if len(result) != len(test.expectedResult) {
				t.Fatalf("[%s] : expected routines len -> %v, got -> %v", t.Name(), len(test.expectedResult), len(result))
			}

			// Assert expected result
			for i, r := range result {
				eq := reflect.DeepEqual(r, test.expectedResult[i])
				if !eq {
					t.Fatalf("[%s] : expected routine entry -> %v, got -> %v", t.Name(), test.expectedResult[i], r)
				}
			}
		})
	}
}

func Test_GetBlueprint_InvalidFilePath(t *testing.T) {
	yamlError := fmt.Sprintf("open %s: The filename, directory name, or volume label syntax is incorrect.", misc.ValidYamlFile)
	errorPrefix := "[Error] : Failed to read blueprint file"
	expectedMessage := fmt.Sprintf("%s, %s", errorPrefix, yamlError)
	expectedTrace := "[StackTrace] : yaml|GetBlueprint"

	tests := []struct {
		name             string
		verbosity        errors.VerbosityLevel
		expectedSeverity errors.SeverityLevel
		expectedMessage  string
	}{
		{
			name:             "errorVerbosity",
			verbosity:        errors.ErrorMessage,
			expectedSeverity: errors.Critical,
			expectedMessage:  fmt.Sprintf("%s\n", expectedMessage),
		}, {
			name:             "stackTraceVerbosity",
			verbosity:        errors.StackTrace,
			expectedSeverity: errors.Critical,
			expectedMessage:  fmt.Sprintf("%s\n", expectedTrace),
		}, {
			name:             "fullVerbosity",
			verbosity:        errors.Full,
			expectedSeverity: errors.Critical,
			expectedMessage:  fmt.Sprintf("%s | %s\n", expectedMessage, expectedTrace),
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			/* ARRANGE */
			errors.SetErrorVerbosity(testCase.verbosity)
			parser := NewParser()

			/* ACT */
			_, err := parser.GetBlueprint(misc.ValidYamlFile)

			/* ASSERT */
			// Assert expected error
			if err == nil {
				t.Fatalf("[%s] : expected error, not nil", t.Name())
			}

			// Assert expected error severity
			if err.Severity() != testCase.expectedSeverity {
				t.Fatalf("[%s] : expected severity -> %v, got -> %v", t.Name(), testCase.expectedSeverity, err.Severity())
			}

			// Assert expected error message
			if err.Error() != testCase.expectedMessage {
				t.Fatalf("[%s] : expected message -> %v, got -> %v", t.Name(), testCase.expectedMessage, err.Error())
			}
		})
	}
}

func Test_GetBlueprint_InvalidFileExtention(t *testing.T) {
	// Invalid blueprint file extention, expected .yaml, got .qwe
	expectedMessage := fmt.Sprintf("[Error] : Invalid blueprint file extention, expected .yaml, got %s", misc.InvalidYamlExtention)
	expectedTrace := "[StackTrace] : yaml|GetBlueprint"

	tests := []struct {
		name             string
		verbosity        errors.VerbosityLevel
		expectedSeverity errors.SeverityLevel
		expectedMessage  string
	}{
		{
			name:             "errorVerbosity",
			verbosity:        errors.ErrorMessage,
			expectedSeverity: errors.Critical,
			expectedMessage:  fmt.Sprintf("%s\n", expectedMessage),
		}, {
			name:             "stackTraceVerbosity",
			verbosity:        errors.StackTrace,
			expectedSeverity: errors.Critical,
			expectedMessage:  fmt.Sprintf("%s\n", expectedTrace),
		}, {
			name:             "fullVerbosity",
			verbosity:        errors.Full,
			expectedSeverity: errors.Critical,
			expectedMessage:  fmt.Sprintf("%s | %s\n", expectedMessage, expectedTrace),
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			/* ARRANGE */
			errors.SetErrorVerbosity(testCase.verbosity)
			fileName := fmt.Sprintf("%s%s", misc.FileName, misc.InvalidYamlExtention)

			// Mock blueprint file
			path, e := misc.CreateTmpYAML(fileName, map[string]interface{}{})
			if e != nil {
				t.Fatalf("[%s] : can't create YAML file, %s", t.Name(), e.Error())
			}
			defer misc.DeleteTmpYAML(path)

			parser := NewParser()

			/* ACT */
			_, err := parser.GetBlueprint(path)

			/* ASSERT */
			// Assert expected error
			if err == nil {
				t.Fatalf("[%s] : expected error -> %v, got -> %v", t.Name(), testCase.expectedMessage, nil)
			}

			// Assert expected error severity
			if err.Severity() != testCase.expectedSeverity {
				t.Fatalf("[%s] : expected severity -> %v, got -> %v", t.Name(), testCase.expectedSeverity, err.Severity())
			}

			// Assert expected error message
			if err.Error() != testCase.expectedMessage {
				t.Fatalf("[%s] : expected message -> %v, got -> %v", t.Name(), testCase.expectedMessage, err.Error())
			}
		})
	}
}

func Test_GetBlueprint_InvalidProjectEntry(t *testing.T) {
	yamlContent := map[string]interface{}{"project": ""}
	yamlError := "yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `` into struct { Name string \"yaml:name\"; Version string \"yaml:version\" }"
	errorPrefix := "[Error] : Failed to unmarshall yaml file"
	expectedMessage := fmt.Sprintf("%s, %s", errorPrefix, yamlError)
	expectedTrace := "[StackTrace] : parseProjectEntry -> yaml|GetBlueprint"

	tests := []struct {
		name             string
		verbosity        errors.VerbosityLevel
		yaml             map[string]interface{}
		expectedSeverity errors.SeverityLevel
		expectedMessage  string
		expectedResult   types.Blueprint
	}{
		{
			name:             "errorVerbosity",
			verbosity:        errors.ErrorMessage,
			yaml:             yamlContent,
			expectedSeverity: errors.Critical,
			expectedMessage:  fmt.Sprintf("%s\n", expectedMessage),
			expectedResult:   types.Blueprint{},
		},
		{
			name:             "stackTraceVerbosity",
			verbosity:        errors.StackTrace,
			yaml:             yamlContent,
			expectedSeverity: errors.Critical,
			expectedMessage:  fmt.Sprintf("%s\n", expectedTrace),
			expectedResult:   types.Blueprint{},
		},
		{
			name:             "fullVerbosity",
			verbosity:        errors.Full,
			yaml:             yamlContent,
			expectedSeverity: errors.Critical,
			expectedMessage:  fmt.Sprintf("%s | %s\n", expectedMessage, expectedTrace),
			expectedResult:   types.Blueprint{},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			/* ARRANGE */
			errors.SetErrorVerbosity(testCase.verbosity)

			// Mock blueprint file
			path, e := misc.CreateTmpYAML(misc.ValidYamlFile, testCase.yaml)
			if e != nil {
				t.Fatalf("[%s] : can't create YAML file, %s", t.Name(), e.Error())
			}
			defer misc.DeleteTmpYAML(path)
			parser := NewParser()

			/* ACT */
			result, err := parser.GetBlueprint(path)

			/* ASSERT */
			// Assert expected error
			if err == nil {
				t.Fatalf("[%s] : expected error -> %v, got -> %v", t.Name(), testCase.expectedMessage, nil)
			}

			// Assert expected error severity
			if err.Severity() != testCase.expectedSeverity {
				t.Fatalf("[%s] : expected severity -> %v, got -> %v", t.Name(), testCase.expectedSeverity, err.Severity())
			}

			// Assert expected error message
			if err.Error() != testCase.expectedMessage {
				t.Fatalf("[%s] : expected message -> %v, got -> %v", t.Name(), testCase.expectedMessage, err.Error())
			}

			// Assert expected result
			if !reflect.DeepEqual(result, testCase.expectedResult) {
				t.Fatalf("[%s] : expected Blueprint -> %v, got -> %v", t.Name(), testCase.expectedResult, result)
			}
		})
	}
}

func Test_GetBlueprint_InvalidRoutinesEntry(t *testing.T) {
	yamlContent := map[string]interface{}{"routines": ""}
	yamlError := "yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `` into []struct { Routine yaml.yamlRoutineEntry \"yaml:routine\" }"
	errorPrefix := "[Error] : Failed to unmarshall yaml file"
	expectedMessage := fmt.Sprintf("%s, %s", errorPrefix, yamlError)
	expectedTrace := "[StackTrace] : parseRoutinesEntry -> yaml|GetBlueprint"

	tests := []struct {
		name             string
		verbosity        errors.VerbosityLevel
		yaml             map[string]interface{}
		expectedSeverity errors.SeverityLevel
		expectedMessage  string
		expectedResult   types.Blueprint
	}{
		{
			name:             "errorVerbosity",
			verbosity:        errors.ErrorMessage,
			yaml:             yamlContent,
			expectedSeverity: errors.Critical,
			expectedMessage:  fmt.Sprintf("%s\n", expectedMessage),
			expectedResult:   types.Blueprint{},
		},
		{
			name:             "stackTraceVerbosity",
			verbosity:        errors.StackTrace,
			yaml:             yamlContent,
			expectedSeverity: errors.Critical,
			expectedMessage:  fmt.Sprintf("%s\n", expectedTrace),
			expectedResult:   types.Blueprint{},
		},
		{
			name:             "fullVerbosity",
			verbosity:        errors.Full,
			yaml:             yamlContent,
			expectedSeverity: errors.Critical,
			expectedMessage:  fmt.Sprintf("%s | %s\n", expectedMessage, expectedTrace),
			expectedResult:   types.Blueprint{},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			/* ARRANGE */
			errors.SetErrorVerbosity(testCase.verbosity)

			// Mock blueprint file
			path, e := misc.CreateTmpYAML(misc.ValidYamlFile, testCase.yaml)
			if e != nil {
				t.Fatalf("[%s] : can't create YAML file, %s", t.Name(), e.Error())
			}
			defer misc.DeleteTmpYAML(path)
			parser := NewParser()

			/* ACT */
			result, err := parser.GetBlueprint(path)

			/* ASSERT */
			// Assert expected error
			if err == nil {
				t.Fatalf("[%s] : expected error -> %v, got -> %v", t.Name(), testCase.expectedMessage, nil)
			}

			// Assert expected error severity
			if err.Severity() != testCase.expectedSeverity {
				t.Fatalf("[%s] : expected severity -> %v, got -> %v", t.Name(), testCase.expectedSeverity, err.Severity())
			}

			// Assert expected error message
			if err.Error() != testCase.expectedMessage {
				t.Fatalf("[%s] : expected message -> %v, got -> %v", t.Name(), testCase.expectedMessage, err.Error())
			}

			// Assert expected result
			if !reflect.DeepEqual(result, testCase.expectedResult) {
				t.Fatalf("[%s] : expected Blueprint -> %v, got -> %v", t.Name(), testCase.expectedResult, result)
			}
		})
	}
}

func Test_GetBlueprint_ValidBlueprintFile(t *testing.T) {
	/* ARRANGE */
	projectEntry := map[string]interface{}{"name": expectedName, "version": expectedVersion}

	firstRoutineFirstStep := map[string]interface{}{
		"name":        "TEST_STEP_0",
		"description": "TEST_STEP_0_DESCRIPTION",
		"attempts":    100,
		"timeout_s":   1000,
		"post": []interface{}{map[string]interface{}{
			"event":   "TEST_STEP_0_POST_EVENT",
			"payload": []string{"PAYLOAD_0", "PAYLOAD_1"},
		}},
	}

	firstRoutineSecondStep := map[string]interface{}{
		"name":        "TEST_STEP_1",
		"description": "TEST_STEP_1_DESCRIPTION",
		"attempts":    200,
		"timeout_s":   2000,
		"post": []interface{}{map[string]interface{}{
			"event":   "TEST_STEP_1_POST_EVENT",
			"payload": []string{"PAYLOAD_0", "PAYLOAD_1", "PAYLOAD_2"},
		}},
	}

	firstRoutineSteps := []interface{}{
		map[string]interface{}{"step": firstRoutineFirstStep},
		map[string]interface{}{"step": firstRoutineSecondStep},
	}

	firstRoutineEntry := map[string]interface{}{
		"name":  "TEST_ROUTINE_0",
		"steps": firstRoutineSteps,
	}

	routinesEntry := []interface{}{
		map[string]interface{}{"routine": firstRoutineEntry},
	}

	yamlContent := map[string]interface{}{"project": projectEntry, "routines": routinesEntry}
	expectedResult := types.Blueprint{
		Project: types.ProjectEntry{Name: expectedName, Version: expectedVersion},
		Routines: []types.RoutineEntry{
			{Name: "TEST_ROUTINE_0", Steps: []types.StepEntry{
				{
					Name:        "TEST_STEP_0",
					Description: "TEST_STEP_0_DESCRIPTION",
					Attempts:    100,
					Timeout_s:   1000,
					Post:        []types.PostEventEntry{{Event: "TEST_STEP_0_POST_EVENT", Payload: []string{"PAYLOAD_0", "PAYLOAD_1"}}},
				},
				{
					Name:        "TEST_STEP_1",
					Description: "TEST_STEP_1_DESCRIPTION",
					Attempts:    200,
					Timeout_s:   2000,
					Post:        []types.PostEventEntry{{Event: "TEST_STEP_1_POST_EVENT", Payload: []string{"PAYLOAD_0", "PAYLOAD_1", "PAYLOAD_2"}}},
				},
			}},
		},
	}

	// Mock blueprint file
	path, e := misc.CreateTmpYAML(misc.ValidYamlFile, yamlContent)
	if e != nil {
		t.Fatalf("[%s] : can't create YAML file, %s", t.Name(), e.Error())
	}
	defer misc.DeleteTmpYAML(path)
	parser := NewParser()

	/* ACT */
	result, err := parser.GetBlueprint(path)

	/* ASSERT */
	// Assert unexpected error
	if err != nil {
		t.Fatalf("[%s] : expected error -> %v, got -> %v", t.Name(), nil, err.Error())
	}

	// Assert expected result
	if !reflect.DeepEqual(result, expectedResult) {
		t.Fatalf("[%s] : expected Blueprint -> %v, got -> %v", t.Name(), expectedResult, result)
	}
}
