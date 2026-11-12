package httpapi

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestValidJSONDecoding(t *testing.T) {
	data := map[string]interface{}{
		"name": "John",
		"age":  30,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Error marshaling JSON: %v", err)
	}

	var decodedData map[string]interface{}
	err = json.Unmarshal(jsonData, &decodedData)
	if err != nil {
		t.Fatalf("Error unmarshaling JSON: %v", err)
	}

	if decodedData["name"] != "John" {
		t.Errorf("Expected name to be 'John', got %v", decodedData["name"])
	}

	if decodedData["age"] != float64(30) {
		t.Errorf("Expected age to be 30, got %v", decodedData["age"])
	}
}

func TestMalformedJSONDecoding(t *testing.T) {
	malformedJSON := `{"name": "John", "age": 30` // Missing closing brace
	var decodedData map[string]interface{}
	err := json.Unmarshal([]byte(malformedJSON), &decodedData)
	if err == nil {
		t.Errorf("Expected error when unmarshaling malformed JSON, but got none")
	}
}

func TestEmptyRequestBody(t *testing.T) {
	_, err := http.NewRequest("POST", "/test", nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
}


func TestUnkownJson(t *testing.T) {
	unknownJSON := `{"unknown_field": "value"}`
	var decodedData map[string]interface{}
	err := json.Unmarshal([]byte(unknownJSON), &decodedData)
	if err != nil {
		t.Fatalf("Error unmarshaling unknown JSON: %v", err)
	}
}

func TestNestedJSON(t *testing.T) {
	nestedJSON := `{"user": {"name": "John", "age": 30}}`
	var decodedData map[string]map[string]interface{}
	err := json.Unmarshal([]byte(nestedJSON), &decodedData)
	if err != nil {
		t.Fatalf("Error unmarshaling nested JSON: %v", err)
	}
}


func TestMultipleJSONValues(t *testing.T) {
	multipleJSON := `{"name": "John"}{"age": 30}`
	var decodedData map[string]interface{}
	err := json.Unmarshal([]byte(multipleJSON), &decodedData)
	if err == nil {
		t.Fatal("Expected an error for multiple JSON values")
	}
}

func TestWrongContentType(t *testing.T) {
	// Simulate a request with a wrong content type
	req, err := http.NewRequest("POST", "/test", nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "text/plain")
}

func TestJsonReponseFormat(t *testing.T) {
	// Simulate a JSON response
	response := map[string]interface{}{
		"name": "John",
		"age":  30,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Error marshaling JSON response: %v", err)
	}

	var decodedResponse map[string]interface{}
	err = json.Unmarshal(jsonResponse, &decodedResponse)
	if err != nil {
		t.Fatalf("Error unmarshaling JSON response: %v", err)
	}

	if decodedResponse["name"] != "John" {
		t.Errorf("Expected name to be 'John', got %v", decodedResponse["name"])
	}

	if decodedResponse["age"] != float64(30) {
		t.Errorf("Expected age to be 30, got %v", decodedResponse["age"])
	}
}


func TestErrorResponseFormat(t *testing.T) {
	// Simulate an error response
	errorResponse := map[string]interface{}{
		"error": "Invalid request",
	}

	jsonErrorResponse, err := json.Marshal(errorResponse)
	if err != nil {
		t.Fatalf("Error marshaling JSON error response: %v", err)
	}

	var decodedErrorResponse map[string]interface{}
	err = json.Unmarshal(jsonErrorResponse, &decodedErrorResponse)
	if err != nil {
		t.Fatalf("Error unmarshaling JSON error response: %v", err)
	}

	if decodedErrorResponse["error"] != "Invalid request" {
		t.Errorf("Expected error to be 'Invalid request', got %v", decodedErrorResponse["error"])
	}
}