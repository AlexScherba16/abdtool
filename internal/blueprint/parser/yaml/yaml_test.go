package yaml

import (
	"abdtool/internal/blueprint/types"
	"abdtool/internal/mics"
	"abdtool/utils/errors"
	"os"
	"reflect"
	"testing"
)

func Test_parseProjectEntry(t *testing.T) {
	expectedName := "test"
	expectedVersion := "2.2.8"

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
			path, e := mics.CreateTmpYAML(mics.ValidYamlFile, test.yaml)
			if e != nil {
				t.Fatalf("[%s] : can't create YAML file, %s", t.Name(), e.Error())
			}
			defer mics.DeleteTmpYAML(path)

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
