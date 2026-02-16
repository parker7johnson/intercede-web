package webhandlers

// Handlers holds dependencies for HTTP handlers
type WebHandlers struct {
	// Add dependencies here (e.g., database, logger, services)
}

// New creates a new Handlers instance with dependencies
func NewWeb() *WebHandlers {
	return &WebHandlers{
		// Initialize dependencies here
	}
}
