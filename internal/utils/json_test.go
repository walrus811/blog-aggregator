package utils

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func TestRespondWithJSON(t *testing.T) {
	w := httptest.NewRecorder()

	type tempResponse struct {
		Message string `json:"message"`
	}

	payload := tempResponse{Message: "ok"}
	code := 200

	RespondWithJSON(w, code, payload)
	if w.Code != code {
		t.Errorf("Status code does not match")
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Content-Type header does not match")
	}

	decoder := json.NewDecoder(w.Body)
	resObj := tempResponse{}
	resDecodeErr := decoder.Decode(&resObj)
	if resDecodeErr != nil {
		t.Errorf("Error decoding response")
	}

	if resObj.Message != "ok" {
		t.Errorf("Body does not match")
	}
}

func TestRespondWithJSONWithErr(t *testing.T) {
	w := httptest.NewRecorder()

	payload := complex(1, 2)
	code := 200

	RespondWithJSON(w, code, payload)

	if w.Code != 500 {
		t.Errorf("Status code does not match")
	}
}

func TestRespondWithError(t *testing.T) {
	w := httptest.NewRecorder()

	errMessage := "Error somewhat"
	code := 500

	RespondWithError(w, code, errMessage)
	if w.Code != code {
		t.Errorf("Status code does not match")
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Content-Type header does not match")
	}

	decoder := json.NewDecoder(w.Body)
	resObj := ErrorResponse{}
	resDecodeErr := decoder.Decode(&resObj)
	if resDecodeErr != nil {
		t.Errorf("Error decoding response")
	}

	if resObj.Error != errMessage {
		t.Errorf("Body does not match")
	}
}
