package trnsl8r

import "log"

// Request represents a data source with its connection details and logging configuration.
type Request struct {
	protocol        string      // Protocol used by the source (e.g., HTTP, HTTPS).
	host            string      // Host address of the source.
	port            int         // Port number for the source connection.
	origin          string      // Origin identifier for the source.
	locale          string      // Locale identifier for the source. - Not used
	customLogger    *log.Logger // Logger instance for logging activities.
	isCustomLogger  bool        // Flag indicating if logging is enabled.
	isLoggingActive bool        // Flag indicating if logging is currently active.
}

// Response represents the result of a translation operation.
type Response struct {
	Original    string `json:"original"`
	Translated  string `json:"translated"`
	Information string `json:"information"`
}

// APIResponse represents a generic response message.
type APIResponse struct {
	Message string `json:"message"`
}

// urlTemplate is a format string used to construct the URL for the translation service.
// It includes placeholders for the protocol, host, and port.
var urlTemplate = "%v://%v:%d/trnsl8r/%v/%v"

type Filter struct {
	field string
}

var LOCALE = Filter{field: "locale"}
var ORIGIN = Filter{field: "origin"}

var filters = []Filter{LOCALE, ORIGIN}
