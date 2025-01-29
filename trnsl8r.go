// Package trnsl8r provides functionality for managing and translating data sources.
package trnsl8r

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

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

func (t Response) String() string {
	return t.Translated
}

// Get sends a request to the translation service to translate the given subject.
// It constructs the URL using the protocol, host, and port defined in the Request struct.
// If any of these fields are missing, it logs an error and returns a Response with the error information.
// It also checks if the subject is empty or contains invalid characters, logging and returning an error if so.
// If the request is successful, it reads the response body, unmarshals the JSON into an APIResponse struct,
// and constructs a Response with the original and translated messages.
// It logs various stages of the process for debugging purposes.
//
// Parameters:
// - subject: The message to be translated.
//
// Returns:
// - Response: A struct containing the original message, translated message, and any additional information.
// - error: An error if any issues occurred during the process.
func (s *Request) Get(subject string) (Response, error) {
	// Check if protocol is defined
	if s.protocol == "" {
		err := fmt.Errorf("No protocol defined")
		s.log(err.Error())
		return Response{Information: err.Error()}, err
	}

	// Check if host is defined
	if s.host == "" {
		err := fmt.Errorf("No host defined")
		s.log(err.Error())
		return Response{Information: err.Error()}, err
	}

	// Check if port is defined
	if s.port == 0 {
		err := fmt.Errorf("No port defined")
		s.log(err.Error())
		return Response{Information: err.Error()}, err
	}

	if s.origin == "" {
		err := fmt.Errorf("No origin defined, and origin identifier is required.")
		s.log(err.Error())
		return Response{Information: err.Error()}, err
	}

	// Check if subject is defined
	if subject == "" {
		err := fmt.Errorf("No message to translate")
		s.log(err.Error())
		return Response{Information: err.Error()}, err
	}

	// Check if subject contains invalid characters
	if strings.Contains(subject, "/") {
		err := fmt.Errorf("Message contains invalid characters")
		s.log(err.Error())
		return Response{Information: err.Error()}, err
	}

	// Construct the full URL
	url := fmt.Sprintf(urlTemplate, s.protocol, s.host, s.port, s.origin, url.QueryEscape(subject))
	s.log(fmt.Sprintf("Request to translate message [%v] by [%v]", subject, url))

	// Send the request via a client
	var client http.Client
	resp, err := client.Get(url)
	if err != nil {
		s.log(err.Error())
		return Response{Information: err.Error()}, err
	}
	defer resp.Body.Close()

	s.log(fmt.Sprintf("Response Status: [%v]", resp.Status))

	// Check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("Response not OK - %s - %d", resp.Status, resp.StatusCode)
		s.log(err.Error())
		return Response{Information: err.Error()}, err
	}

	// Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		s.log(err.Error())
		return Response{Information: err.Error()}, err
	}

	// Unmarshal the JSON byte slice to a predefined struct
	var reponse APIResponse
	err = json.Unmarshal(bodyBytes, &reponse)
	if err != nil {
		s.log(err.Error())
		return Response{Original: subject, Translated: subject, Information: err.Error()}, err
	}

	// Construct the translated response
	var translated Response
	translated.Original = subject
	translated.Translated = reponse.Message
	translated.Information = ""

	// Log the translation result
	msg := fmt.Sprintf("Original:{{%v}} Translation:{{%v}} Information:{{%v}}", translated.Original, translated.Translated, translated.Information)
	s.log(msg)

	return translated, nil
}

// NewRequest creates a new Request instance with default values for logging configuration.
// Returns:
// - Request: A new Request instance with logging disabled.
func NewRequest() Request {
	return Request{isCustomLogger: false, isLoggingActive: true}
}

// WithProtocol sets the protocol for the Request.
// Parameters:
// - protocol: The protocol to be used (e.g., HTTP, HTTPS).
// Returns:
// - Request: The updated Request instance.
func (s Request) WithProtocol(protocol string) Request {
	s.protocol = protocol
	return s
}

// WithHost sets the host for the Request.
// Parameters:
// - host: The host address of the source.
// Returns:
// - Request: The updated Request instance.
func (s Request) WithHost(host string) Request {
	s.host = host
	return s
}

// WithPort sets the port for the Request.
// Parameters:
// - port: The port number for the source connection.
// Returns:
// - Request: The updated Request instance.
func (s Request) WithPort(port int) Request {
	s.port = port
	return s
}

// WithOriginOf sets the origin identifier for the Request.
// Parameters:
// - origin: The origin identifier for the source.
// Returns:
// - Request: The updated Request instance.
func (s Request) WithOriginOf(origin string) Request {
	s.origin = origin
	return s
}

// WithLogger sets the logger for the Request and enables logging.
// Parameters:
// - logger: The logger instance for logging activities.
// Returns:
// - Request: The updated Request instance with logging enabled.
func (s Request) WithLogger(logger *log.Logger) Request {
	s.customLogger = logger
	s.isCustomLogger = true
	return s
}

// EnableLogging enables logging for the Request.
// Returns:
// - Request: The updated Request instance with logging active.
func (s Request) EnableLogging() Request {
	s.isLoggingActive = true
	return s
}

// DisableLogging disables logging for the Request.
// Returns:
// - Request: The updated Request instance with logging inactive.
func (s Request) DisableLogging() Request {
	s.isLoggingActive = false
	return s
}

// String constructs and returns the URL string for the Request.
// Returns:
// - string: The constructed URL string.
func (s Request) String() string {
	return fmt.Sprintf(urlTemplate, s.protocol, s.host, s.port)
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
func (s Request) Validate() error {
	if s.protocol == "" {
		return fmt.Errorf("protocol is required")
	}
	if s.host == "" {
		return fmt.Errorf("host is required")
	}
	if s.port == 0 {
		return fmt.Errorf("port is required")
	}
	return nil
}

// Spew outputs the contents of the Request struct to the log.
func (s Request) Spew() {
	message := fmt.Sprintf(
		"Request struct contents:\nProtocol: %s\nHost: %s\nPort: %d\nOrigin: %s\nLocale: %s\nLogger: %+v\nIsCustomLogger: %t\nIsLoggingActive: %t",
		s.protocol, s.host, s.port, s.origin, s.locale, s.customLogger, s.isCustomLogger, s.isLoggingActive,
	)
	s.log(message)
}
