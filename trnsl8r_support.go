package trnsl8r

import (
	"fmt"
	"log"
)

// String constructs and returns the URL string for the Request.
// Returns:
// - string: The constructed URL string.
func (s Request) String() string {
	return fmt.Sprintf(urlTemplate, s.protocol, s.host, s.port)
}

func (t Response) String() string {
	return t.Translated
}

// log logs a message using the Request's logger if logging is enabled, otherwise logs to the default logger.
// Parameters:
// - message: The message to be logged.
func (s Request) log(message string) {
	if s.isLoggingActive {
		if s.isCustomLogger {
			s.customLogger.Println(message)
		} else {
			log.Println(message)
		}
	}
}

// Validate checks if the required fields of the Request are set.
// Returns:
// - error: An error if any required fields are missing.
func (s Request) Validate(subject string) error {
	if s.protocol == "" {
		return fmt.Errorf("protocol is required")
	}
	if s.host == "" {
		return fmt.Errorf("host is required")
	}
	if s.port == 0 {
		return fmt.Errorf("port is required")
	}
	if s.origin == "" {
		return fmt.Errorf("no origin defined, and origin identifier is required.")
	}
	// Check if subject is defined
	if subject == "" {
		return fmt.Errorf("no message to translate")
	}
	return nil
}

// Spew outputs the contents of the Request struct to the log.
func (s Request) Spew() {
	message := fmt.Sprintf(
		"Request struct contents:\nProtocol: %s\nHost: %s\nPort: %d\nOrigin: %s\nLocale: %s\nLogger: %+v\nIsCustomLogger: %t\nIsLoggingActive: %t \nFilters: %v",
		s.protocol, s.host, s.port, s.origin, s.locale, s.customLogger, s.isCustomLogger, s.isLoggingActive, s.filters,
	)
	s.log(message)
}
