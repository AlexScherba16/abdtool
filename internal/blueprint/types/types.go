package types

// Blueprint - encapsulates the configuration blueprint of an application
type Blueprint struct {
	// Project - represents the essential information about the project
	Project ProjectEntry

	// Routines - represents slice of configurations for routine runners
	Routines []RoutineEntry
}

// ProjectEntry - represents essential details about project. It captures
// the project's name and its version
type ProjectEntry struct {
	// Name - contains specified title of the project
	Name string

	// Version - denotes the version string of the project
	Version string
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

// PostEventEntry - represents a specific event that should be triggered
// during the execution of a certain step. It provides details about the event name
// and the associated data (payload) that might be needed when the event is posted
type PostEventEntry struct {
	// Event - is the name or identifier of the event to be triggered
	Event string

	// Payload contains the list of fields or data associated with the event
	Payload []string
}

type StepEntry struct {
	// Name - contains specified name of step
	Name string

	// Description - contains more detailed info about step
	Description string

	// Attempts - specifies the number of retry attempts to execute the step in case of failure
	Attempts int

	// Timeout_s - time limit to execute step
	Timeout_s int

	// Post - is a slice of PostEventEntry. Each entry represents an event
	// that should be posted during the execution of a specific step,
	// this allows the system to be notified or updated based on the results
	// or progress of the step being executed
	Post []PostEventEntry

	// Subscribe - is a slice of event names to which Step is subscribed
	// IMPORTANT: Step will be blocked and wait until one of these events is emitted
	// if the expected event is not emitted, there's a risk that the Step could hang indefinitely
	// if all events occurs, the Step will continue its operation
	Subscribe []string
}
