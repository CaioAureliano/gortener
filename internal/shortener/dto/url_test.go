package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValid(t *testing.T) {

	tests := []struct {
		name        string
		gotUrl      string
		expectedRes bool
	}{
		{
			name:        "should be return true with valid URL",
			gotUrl:      "http://www.example.com",
			expectedRes: true,
		},
		{
			name:        "should be return false with simple text (invalid URL)",
			gotUrl:      "example",
			expectedRes: false,
		},
		{
			name:        "should be return false with invalid URL (just dot com)",
			gotUrl:      "example.",
			expectedRes: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := UrlRequest{Url: tt.gotUrl}
			res := req.IsValid()

			assert.Equal(t, tt.expectedRes, res)
		})
	}
}
