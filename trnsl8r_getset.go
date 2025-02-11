package trnsl8r

import (
	"fmt"
	"log"
	"slices"

	"github.com/mt1976/frantic-core/id"
)

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

// FromOrigin sets the origin identifier for the Request.
// Parameters:
// - origin: The origin identifier for the source.
// Returns:
// - Request: The updated Request instance.
func (s Request) FromOrigin(origin string) Request {
	var err error
	s.origin, err = id.GetUUIDv2WithPayload(origin)
	if err != nil {
		log.Fatal(err)
	}
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

func (s Request) WithFilter(filter Filter, value string) (Request, error) {

	if !slices.Contains(filters, filter) {
		return s, fmt.Errorf("%v is not a invalid filter, valid filters are %v", filter, filters)
	}

	s.filters = append(s.filters, Filter{key: filter.key, value: value})

	return s, nil
}
