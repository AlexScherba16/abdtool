package blueprint

import (
	"abdtool/internal/mics"
	"abdtool/utils/errors"
	"fmt"
	"log"
	"os"
	"testing"
)

func TestNewBlueprint_InvalidPath(t *testing.T) {
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
			expectedMessage:  "[Error] : Blueprint file doesn't exist\n",
		}, {
			name:             "stackTraceVerbosity",
			verbosity:        errors.StackTrace,
			expectedSeverity: errors.Critical,
			expectedMessage:  "[StackTrace] : NewBlueprint\n",
		}, {
			name:             "fullVerbosity",
			verbosity:        errors.Full,
			expectedSeverity: errors.Critical,
			expectedMessage:  "[Error] : Blueprint file doesn't exist | [StackTrace] : NewBlueprint\n",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			/* ARRANGE */
			errors.SetErrorVerbosity(testCase.verbosity)

			/* ACT */
			_, err := NewBlueprint("")

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

func TestNewBlueprint_InvalidFileExtention(t *testing.T) {
	expectedMessage := fmt.Sprintf("[Error] : Invalid plueprint file extention, expected .yaml, got %s", mics.InvalidYamlExtention)
	expectedTrace := "[StackTrace] : NewBlueprint"

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
			fileName := fmt.Sprintf("%s%s", mics.FileName, mics.InvalidYamlExtention)

			// Mock blueprint file
			path, e := mics.CreateTmpYAML(fileName, map[string]interface{}{})
			if e != nil {
				t.Fatalf("[%s] : can't create YAML file, %s", t.Name(), e.Error())
			}

			// Delete the temporary YAML file
			defer func() {
				if err := os.Remove(path); err != nil {
					log.Fatalf("[%s] : failed to delete temporary YAML file: %s", t.Name(), e.Error())
				}
			}()

			/* ACT */
			_, err := NewBlueprint(path)

			/* ASSERT */
			// Assert expected error
			if err == nil {
				t.Fatalf("[%s] : expected error -> %v, got -> %v", t.Name(), nil, err.Error())
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

func TestNewBlueprint_ValidYamlFile(t *testing.T) {
	/* ARRANGE */
	t.Fatal("Implement")
}
