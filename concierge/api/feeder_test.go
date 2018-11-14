package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetTitleHash(t *testing.T) {
	h1 := getTitleHash("A-Title")
	h2 := getTitleHash("Diff Title")
	hDup := getTitleHash("A-Title")

	for _, tc := range []struct {
		name string
		hashes []string

	}
}