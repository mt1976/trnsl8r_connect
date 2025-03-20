package trnsl8r

import (
	"fmt"
	"log"

	"github.com/mt1976/frantic-core/commonErrors"
)

// String constructs and returns the URL string for the Request.
// Returns:
// - string: The constructed URL string.
func (s Request) String() string {
	return fmt.Sprintf(urlTemplate_Translate, s.protocol, s.host, s.port)
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

// log logs a message using the Request's logger if logging is enabled, otherwise logs to the default logger.
// Parameters:
// - message: The message to be logged.
func (s Request) logError(message error) {
	if s.isLoggingActive {
		if s.isCustomLogger {
			s.customLogger.Panicln(message.Error())
		} else {
			log.Println("[ERROR]" + message.Error())
		}
		panic(message)
	}
}

// Validate checks if the required fields of the Request are set.
// Returns:
// - error: An error if any required fields are missing.
func (s Request) Validate(message string) error {
	if s.protocol == "" {
		return commonErrors.ErrProtocolIsRequired
	}
	if s.host == "" {
		return commonErrors.ErrHostIsRequired
	}
	if s.port == 0 {
		return commonErrors.ErrPortIsRequired
	}
	if s.origin == "" {
		return commonErrors.ErrOriginIsRequired
	}
	// Check if subject is defined
	if message == "" {
		return commonErrors.ErrNoMessageToTranslate
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
