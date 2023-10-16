package errors

import "fmt"

// SeverityLevel - error severity type and available values
type SeverityLevel int

const (
	Low SeverityLevel = iota
	Medium
	High
	Critical
)

// VerbosityLevel - error verbosity type and available values
type VerbosityLevel int

const (
	ErrorMessage VerbosityLevel = iota
	StackTrace
	Full
)

// Predefined constants for output
const (
	errorMessage = "Error"
	stackTrace   = "StackTrace"
)

// errorVerbosity - private global error verbosity level variable
var errorVerbosity VerbosityLevel = ErrorMessage

// SetErrorVerbosity - sets the global error verbosity level for the application
// This determines the amount of detail or information that will be included in error messages
//
// Parameters:
//  - level: The desired verbosity level to set for error reporting.
func SetErrorVerbosity(level VerbosityLevel) { errorVerbosity = level }

// CustomError - represents a customized error type for the application
// It provides additional context such as severity and stack trace
// to aid in debugging and error handling
type CustomError struct {
	// severity - indicates the severity level of the error,
	// use to extend error handling strategy.
	severity SeverityLevel

	// message - provides a human-readable error description.
	message string

	// stackTrace - captures the call stack at the point where the error was raised.
	stackTrace string
}

// NewCustomError - creates and returns a new instance of CustomError
// It allows setting the severity level, a human-readable message, and the point in the code where the error rized
//
// Parameters:
//  - severity: The severity level of the error
//  - message: A human-readable error message
//  - errorSourcePoint: The point in the code where the error rized
//
// Returns:
//  *CustomError: A new instance of CustomError with the specified parameters
func NewCustomError(severity SeverityLevel, message string, errorSourcePoint string) *CustomError {
	return &CustomError{
		severity:   severity,
		message:    message,
		stackTrace: errorSourcePoint,
	}
}

// Severity - retreives error severity level
//
// Returns:
//  SeverityLevel: A instance of SeverityLevel
func (e *CustomError) Severity() SeverityLevel { return e.severity }

// AppendStackTrace - appends the provided actor to the existing stack trace of the CustomError
// If the current stack trace is not empty, it separates the new actor with an arrow ("->")
//
// Parameters:
//  - actor: The entity or action to append to the current stack trace. It can represent
//  a function name, method, or any other identifier that gives context to the error's origin
func (e *CustomError) AppendStackTrace(actor string) {
	if e.stackTrace != "" {
		e.stackTrace += " -> "
	}
	e.stackTrace += actor
}

// Error returns a string representation of the CustomError based on the set error verbosity level
// Depending on the verbosity, it can return just the error message, just the stack trace,
// or both. The method panics if an invalid verbosity level is encountered
//
// Returns:
//   string: A formatted string representation of the error based on the error verbosity level
func (e *CustomError) Error() string {
	switch errorVerbosity {
	case ErrorMessage:
		return fmt.Sprintf("[%s] : %s\n", errorMessage, e.message)
	case StackTrace:
		return fmt.Sprintf("[%s] : %s\n", stackTrace, e.stackTrace)
	case Full:
		return fmt.Sprintf("[%s] : %s\n[%s] : %s\n", errorMessage, e.message, stackTrace, e.stackTrace)
	default:
		panic("Invalid error verbosity level")
	}
}
