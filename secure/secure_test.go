package secure

import (
	"net/http"
	"testing"
)

func TestIsAuthorizedSuccess(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Error("Unable to create request")
	}

	req.Header["Authorization"] = []string{"Bearer AUTH-TOKEN-1"}

	if isAuthorized(req) {
		t.Log("Request with correct Auth token was correctly processed.")
	} else {
		t.Error("Request with correct Auth token failed.")
	}
}

func TestIsAuthorizedFailToTokenType(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Error("Unable to create request")
	}

	req.Header["Authorization"] = []string{"Token AUTH-TOKEN-1"}

	if isAuthorized(req) {
		t.Error("Request with incorrect Auth token type was successfully processed.")
	} else {
		t.Log("Request with incorrect Auth token type failed as expected.")
	}
}
