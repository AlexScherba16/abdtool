package types

import (
	"abdtool/internal/mics"
	"abdtool/utils/errors"
	"fmt"
	"os"
	"testing"
)

func TestNewProjectEntry(t *testing.T) {
	expectedName := "test"
	expectedVersion := "2.2.8"

	tests := []struct {
		name           string
		yaml           map[string]interface{}
		expectedError  *errors.CustomError
		expectedResult ProjectEntry
	}{
		{
			name:           "emptyProjectEntry",
			yaml:           map[string]interface{}{},
			expectedError:  errors.NewCustomError(errors.Critical, fmt.Sprintf("%q must be provided", "project.name"), "NewProjectEntry"),
			expectedResult: ProjectEntry{},
		},
		{
			name: "emptyProjectName",
			yaml: map[string]interface{}{
				"project": map[string]interface{}{},
			},
			expectedError:  errors.NewCustomError(errors.Critical, fmt.Sprintf("%q must be provided", "project.name"), "NewProjectEntry"),
			expectedResult: ProjectEntry{},
		},
		{
			name: "emptyProjectVersion",
			yaml: map[string]interface{}{
				"project": map[string]interface{}{
					"name": expectedName,
				},
			},
			expectedError:  errors.NewCustomError(errors.Critical, fmt.Sprintf("%q must be provided", "project.version"), "NewProjectEntry"),
			expectedResult: ProjectEntry{},
		},
		{
			name: "validProjectEntry",
			yaml: map[string]interface{}{
				"project": map[string]interface{}{
					"name":    expectedName,
					"version": expectedVersion,
				},
			},
			expectedError:  nil,
			expectedResult: ProjectEntry{Name: expectedName, Version: expectedVersion},
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
			project, err := NewProjectEntry(data)

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

			// Assert required Name field
			if project.Name != test.expectedResult.Name {
				t.Fatalf("[%s] : expected name -> %v, got -> %v", t.Name(), test.expectedResult.Name, project.Name)
			}

			// Assert required Version field
			if project.Version != test.expectedResult.Version {
				t.Fatalf("[%s] : expected version -> %v, got -> %v", t.Name(), test.expectedResult.Version, project.Version)
			}
		})
	}
}
