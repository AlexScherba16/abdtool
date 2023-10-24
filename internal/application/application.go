package application

// IApplication - represents the interface for an application's lifecycle management
// implementers should provide functionality for starting and gracefully shutting down the application
type IApplication interface {
	// Run - initiates the launch point of the application
	// implementations should generally block until the application is ready to be shut down
	// generally expects that it should run as a goroutine
	Run()

	// Shutdown - provides a means to gracefully terminate the application
	// implementation should clean up resources and finish any outstanding tasks
	Shutdown()
}
