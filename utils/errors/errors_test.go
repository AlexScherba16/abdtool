package errors

import (
	"testing"
)

func TestNewCustomError_ErrorMessageVerbosity(t *testing.T) {
	tests := []struct {
		name             string
		severity         SeverityLevel
		message          string
		source           string
		expectedSeverity SeverityLevel
		expectedMessage  string
	}{
		{
			name:             "lowSeverity",
			severity:         Low,
			message:          "test Low severity message",
			source:           "function_a",
			expectedSeverity: Low,
			expectedMessage:  "[Error] : test Low severity message\n",
		},
		{
			name:             "mediumSeverity",
			severity:         Medium,
			message:          "test Medium severity message",
			source:           "function_a",
			expectedSeverity: Medium,
			expectedMessage:  "[Error] : test Medium severity message\n",
		},
		{
			name:             "highSeverity",
			severity:         High,
			message:          "test High severity message",
			source:           "function_a",
			expectedSeverity: High,
			expectedMessage:  "[Error] : test High severity message\n",
		},
		{
			name:             "criticalSeverity",
			severity:         Critical,
			message:          "test Critical severity message",
			source:           "function_a",
			expectedSeverity: Critical,
			expectedMessage:  "[Error] : test Critical severity message\n",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			/* ARRANGE */
			globalErrorVerbosity = ErrorMessage

			/* ACT */
			err := NewCustomError(testCase.severity, testCase.message, testCase.source)

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

func TestNewCustomError_StackTraceVerbosity(t *testing.T) {
	tests := []struct {
		name             string
		severity         SeverityLevel
		message          string
		source           string
		expectedSeverity SeverityLevel
		expectedMessage  string
	}{
		{
			name:             "lowSeverity",
			severity:         Low,
			message:          "test Low severity message",
			source:           "lowSeverityFunction",
			expectedSeverity: Low,
			expectedMessage:  "[StackTrace] : lowSeverityFunction\n",
		},
		{
			name:             "mediumSeverity",
			severity:         Medium,
			message:          "test Medium severity message",
			source:           "mediumSeverityFunction",
			expectedSeverity: Medium,
			expectedMessage:  "[StackTrace] : mediumSeverityFunction\n",
		},
		{
			name:             "highSeverity",
			severity:         High,
			message:          "test High severity message",
			source:           "highSeverityFunction",
			expectedSeverity: High,
			expectedMessage:  "[StackTrace] : highSeverityFunction\n",
		},
		{
			name:             "criticalSeverity",
			severity:         Critical,
			message:          "test Critical severity message",
			source:           "criticalSeverityFunction",
			expectedSeverity: Critical,
			expectedMessage:  "[StackTrace] : criticalSeverityFunction\n",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			/* ARRANGE */
			globalErrorVerbosity = StackTrace

			/* ACT */
			err := NewCustomError(testCase.severity, testCase.message, testCase.source)

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

func TestNewCustomError_FullVerbosity(t *testing.T) {
	tests := []struct {
		name             string
		severity         SeverityLevel
		message          string
		source           string
		expectedSeverity SeverityLevel
		expectedMessage  string
	}{
		{
			name:             "lowSeverity",
			severity:         Low,
			message:          "test Low severity message",
			source:           "lowSeverityFunction",
			expectedSeverity: Low,
			expectedMessage:  "[Error] : test Low severity message | [StackTrace] : lowSeverityFunction\n",
		},
		{
			name:             "mediumSeverity",
			severity:         Medium,
			message:          "test Medium severity message",
			source:           "mediumSeverityFunction",
			expectedSeverity: Medium,
			expectedMessage:  "[Error] : test Medium severity message | [StackTrace] : mediumSeverityFunction\n",
		},
		{
			name:             "highSeverity",
			severity:         High,
			message:          "test High severity message",
			source:           "highSeverityFunction",
			expectedSeverity: High,
			expectedMessage:  "[Error] : test High severity message | [StackTrace] : highSeverityFunction\n",
		},
		{
			name:             "criticalSeverity",
			severity:         Critical,
			message:          "test Critical severity message",
			source:           "criticalSeverityFunction",
			expectedSeverity: Critical,
			expectedMessage:  "[Error] : test Critical severity message | [StackTrace] : criticalSeverityFunction\n",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			/* ARRANGE */
			globalErrorVerbosity = Full

			/* ACT */
			err := NewCustomError(testCase.severity, testCase.message, testCase.source)

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

func TestNewCustomError_AppendStackTraceWithVerbosityLevel(t *testing.T) {
	trace := []string{"foo", "bar", "main"}

	tests := []struct {
		name             string
		verbosity        VerbosityLevel
		severity         SeverityLevel
		message          string
		source           string
		stacktrace       []string
		expectedSeverity SeverityLevel
		expectedMessage  string
	}{
		{
			name:             "errorMessageVerbosity",
			verbosity:        ErrorMessage,
			severity:         Critical,
			message:          "assert error verbosity",
			source:           "sourceFunc",
			stacktrace:       trace,
			expectedSeverity: Critical,
			expectedMessage:  "[Error] : assert error verbosity\n",
		},
		{
			name:             "stackTraceVerbosity",
			verbosity:        StackTrace,
			severity:         Critical,
			message:          "assert stacktrace verbosity",
			source:           "sourceFunc",
			stacktrace:       trace,
			expectedSeverity: Critical,
			expectedMessage:  "[StackTrace] : sourceFunc -> foo -> bar -> main\n",
		},
		{
			name:             "fullVerbosity",
			verbosity:        Full,
			severity:         Critical,
			message:          "assert full verbosity",
			source:           "sourceFunc",
			stacktrace:       trace,
			expectedSeverity: Critical,
			expectedMessage:  "[Error] : assert full verbosity | [StackTrace] : sourceFunc -> foo -> bar -> main\n",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			/* ARRANGE */
			SetErrorVerbosity(testCase.verbosity) // same as globalErrorVerbosity = testCase.verbosity

			/* ACT */
			err := NewCustomError(testCase.severity, testCase.message, testCase.source)

			// simulate error passing throug the functions
			for _, tracePoint := range testCase.stacktrace {
				err.AppendStackTrace(tracePoint)
			}

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

func TestNewCustomError_InvalidErrorVerbosity(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	err := NewCustomError(Medium, "Test error", "TestFunc")
	globalErrorVerbosity = VerbosityLevel(1337) // set invalid verbosity level
	_ = err.Error()                             // this should panic
}
