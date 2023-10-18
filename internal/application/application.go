package application

// IApplication - represents the interface for an application's lifecycle management
// Implementers should provide functionality for starting and gracefully shutting down the application
type IApplication interface {
	// Run - initiates the launch point of the application
	// Implementations should generally block until the application is ready to be shut down
	// Generally expects that it should run as a goroutine
	Run()

	// Shutdown - provides a means to gracefully terminate the application
	// Implementation should clean up resources and finish any outstanding tasks
	Shutdown()
}
