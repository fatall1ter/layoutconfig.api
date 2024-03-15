package infra

import (
	"testing"

	"github.com/labstack/echo/v4"
)

func TestServer_responserMIME(t *testing.T) {

	type args struct {
		c       echo.Context
		code    int
		payload interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{}
			if err := s.responserMIME(tt.args.c, tt.args.code, tt.args.payload); (err != nil) != tt.wantErr {
				t.Errorf("Server.responserMIME() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
