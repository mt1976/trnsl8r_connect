// Package trnsl8r provides functionality for managing and translating data sources.
package trnsl8r

import (
	"log"
	"reflect"
	"testing"
)

func TestRequest_Get(t *testing.T) {
	type fields struct {
		protocol        string
		host            string
		port            int
		origin          string
		locale          string
		customLogger    *log.Logger
		isCustomLogger  bool
		isLoggingActive bool
	}
	type args struct {
		subject string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Response
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Test Case 1",
			fields: fields{
				protocol:        "HTTP",
				host:            "localhost",
				port:            8080,
				origin:          "testorigin",
				locale:          "en",
				customLogger:    nil,
				isCustomLogger:  false,
				isLoggingActive: false,
			},
			args: args{
				subject: "Hello",
			},
			want: Response{
				Original:    "Hello",
				Translated:  "Hello",
				Information: "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Request{
				protocol:        tt.fields.protocol,
				host:            tt.fields.host,
				port:            tt.fields.port,
				origin:          tt.fields.origin,
				locale:          tt.fields.locale,
				customLogger:    tt.fields.customLogger,
				isCustomLogger:  tt.fields.isCustomLogger,
				isLoggingActive: tt.fields.isLoggingActive,
			}
			got, err := s.Get(tt.args.subject)
			if (err != nil) != tt.wantErr {
				t.Errorf("Request.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Request.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
