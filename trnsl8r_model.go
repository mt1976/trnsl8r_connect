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
	filters         []Filter
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

// Locale response
type LocaleResponse struct {
	Locales []struct {
		Locale string `json:"locale"`
		Name   string `json:"name"`
	} `json:"locales"`
	Message string `json:"message"`
}

// urlTemplate_Translate is a format string used to construct the URL for the translation service.
// It includes placeholders for the protocol, host, and port.
var urlTemplate_Translate = "%v://%v:%d/trnsl8r/%v/%v"
var urlTemplate_Locales = "%v://%v:%d/locales"

type Filter struct {
	key   string
	value string
}

var LOCALE = Filter{key: "locale"}
var ORIGIN = Filter{key: "origin"}

var filters = []Filter{LOCALE, ORIGIN}
