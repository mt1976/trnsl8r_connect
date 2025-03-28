// Package trnsl8r provides functionality for managing and translating data sources.
package trnsl8r

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/mt1976/frantic-core/commonErrors"
	"github.com/mt1976/frantic-core/htmlHelpers"
)

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

	// Validate the request
	err := s.Validate(subject)
	if err != nil {
		s.log(err.Error())
		return Response{Information: err.Error()}, err
	}

	origSubject := subject
	subject, _ = htmlHelpers.ToPathSafe(subject)
	// Construct the full URL
	base := fmt.Sprintf(urlTemplate_Translate, s.protocol, s.host, s.port, s.origin, subject)

	xx, err := url.Parse(base)
	if err != nil {
		s.log(err.Error())
		return Response{Information: err.Error()}, err
	}
	//base = xx.String()

	//s.log(fmt.Sprintf("Request to translate message [%v] by [%v]", origSubject, base))

	q := xx.Query()

	//fmt.Printf("s.filters: %v\n", s.filters)
	// Add filters to the URL
	for _, filter := range s.filters {
		//fmt.Printf("filter: %v %v\n", filter.key, filter.value)
		yy, err := htmlHelpers.ToPathSafe(filter.value)
		if err != nil {
			s.log(err.Error())
			return Response{Information: err.Error()}, err
		}
		q.Add(filter.key, yy)
	}
	xx.RawQuery = q.Encode()

	//fmt.Printf("xx.String(): %v\n", xx.String())
	s.log(fmt.Sprintf("Request to translate message [%v] by [%v]", origSubject, xx.String()))

	//os.Exit(0)

	// Send the request via a client
	var client http.Client
	resp, err := client.Get(xx.String())
	if err != nil {
		s.log(err.Error())
		return Response{Information: err.Error()}, err
	}
	defer resp.Body.Close()

	//s.log(fmt.Sprintf("Response Status: [%v]", resp.Status))

	// Check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			s.log(err.Error())
			return Response{Information: err.Error()}, err
		}
		var reponse APIResponse
		err = json.Unmarshal(bodyBytes, &reponse)
		if err != nil {
			s.log(err.Error())
			return Response{Original: subject, Translated: subject, Information: err.Error()}, err
		}

		//err = commonErrors.WrapError(fmt.Errorf("[ERROR!] - Status=[%s] Reason=[%v]", resp.Status, reponse.Message))
		err = commonErrors.WrapInvalidHttpReturnStatusWithMessageError(resp.Status, reponse.Message)
		s.log(err.Error())
		return Response{Information: reponse.Message}, err
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
	msg := fmt.Sprintf("Original:[%v] Request:[%v] Translation:[%v] Information:[%v]", origSubject, translated.Original, translated.Translated, translated.Information)
	s.log(msg)

	return translated, nil
}

// NewRequest creates a new Request instance with default values for logging configuration.
// Returns:
// - Request: A new Request instance with logging disabled.
func NewRequest() Request {
	return Request{isCustomLogger: false, isLoggingActive: true}
}

func (s Request) GetLocales() (LocaleResponse, error) {

	base := fmt.Sprintf(urlTemplate_Locales, s.protocol, s.host, s.port)

	xx, err := url.Parse(base)
	if err != nil {
		s.log(err.Error())
		return LocaleResponse{Message: err.Error()}, err
	}

	s.log(fmt.Sprintf("Request to fetch valid locales [%v]", xx.String()))

	// Send the request via a client
	var client http.Client
	resp, err := client.Get(xx.String())
	if err != nil {
		s.log(err.Error())
		return LocaleResponse{Message: err.Error()}, err
	}
	defer resp.Body.Close()

	// Check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			s.log(err.Error())
			return LocaleResponse{Message: err.Error()}, err
		}
		var reponse LocaleResponse
		var returnData LocaleResponse
		err = json.Unmarshal(bodyBytes, &reponse)
		if err != nil {
			s.log(err.Error())

			returnData.Message = err.Error()

			return returnData, err
		}

		err = commonErrors.WrapInvalidHttpReturnStatusWithMessageError(resp.Status, reponse.Message)
		s.log(err.Error())
		return LocaleResponse{Message: err.Error()}, err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		s.log(err.Error())
		return LocaleResponse{Message: err.Error()}, err
	}
	var reponse LocaleResponse
	err = json.Unmarshal(bodyBytes, &reponse)
	if err != nil {
		s.log(err.Error())
		reponse.Message = err.Error()
		return reponse, err
	}
	var returnData LocaleResponse
	for _, v := range reponse.Locales {
		//s.log(fmt.Sprintf("Locale: [%v] Name: [%v]", v.Locale, v.Name))
		v.Locale, _ = htmlHelpers.FromPathSafe(v.Locale)
		v.Name, _ = htmlHelpers.FromPathSafe(v.Name)
		returnData.Locales = append(returnData.Locales, v)
	}
	returnData.Message, _ = htmlHelpers.FromPathSafe(reponse.Message)
	return returnData, nil
}

func (s *Request) Localise(subject, locale string) (Response, error) {

	LCL := LOCALE
	LCL.value = locale
	// Forced localisation
	s.filters = append(s.filters, LCL)
	return s.Get(subject)
}
