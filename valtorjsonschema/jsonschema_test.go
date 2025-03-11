// Copyright 2025 David Stotijn
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package valtorjsonschema

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/invopop/jsonschema"
)

func TestParseJSONSchema(t *testing.T) {
	validData := map[string]any{
		"name":   "John Doe",
		"age":    int64(30),
		"height": 1.75,
		"email":  "john@example.com",
	}
	invalidData := map[string]any{
		"name":   "John123", // contains numbers
		"age":    int64(0),  // not greater than 0
		"height": 4.0,       // too tall
		"email":  "invalid-email",
	}
	missingRequiredData := map[string]any{
		"height": 1.75,
		"email":  "john@example.com",
		// missing required fields: name and age
	}

	schemaBytes, err := os.ReadFile("testdata/basic.json")
	if err != nil {
		t.Fatalf("failed to read schema file: %v", err)
	}

	var jsonSchema jsonschema.Schema
	err = json.Unmarshal(schemaBytes, &jsonSchema)
	if err != nil {
		t.Fatalf("failed to unmarshal schema: %v", err)
	}

	// Parse the JSON schema.
	valtorSchema, err := ParseJSONSchema[any](jsonSchema)
	if err != nil {
		t.Fatalf("failed to parse schema: %v", err)
	}

	// Test valid data against parsed schema.
	err = valtorSchema.Validate(validData)
	if err != nil {
		t.Errorf("expected valid data to pass validation, got error: %v", err)
	}

	// Test invalid data against parsed schema.
	err = valtorSchema.Validate(invalidData)
	if err == nil {
		t.Error("expected invalid data to fail validation, got no error")
	}

	// Test missing required fields.
	err = valtorSchema.Validate(missingRequiredData)
	if err == nil {
		t.Error("expected missing required fields to fail validation, got no error")
	}
}

func TestParseJSONSchemaErrors(t *testing.T) {
	tests := []struct {
		name          string
		schema        jsonschema.Schema
		expectedError string
	}{
		{
			name: "missing type",
			schema: jsonschema.Schema{
				Type: "",
			},
			expectedError: ErrInvalidType.Error(),
		},
		{
			name: "invalid type",
			schema: jsonschema.Schema{
				Type: "foobar",
			},
			expectedError: ErrInvalidType.Error(),
		},
		{
			name: "invalid exclusive minimum for integer",
			schema: jsonschema.Schema{
				Type:    "integer",
				Minimum: json.Number("invalid"),
			},
			expectedError: "invalid `minimum` value \"invalid\"",
		},
		{
			name: "invalid exclusive maximum for integer",
			schema: jsonschema.Schema{
				Type:    "integer",
				Maximum: json.Number("invalid"),
			},
			expectedError: "invalid `maximum` value \"invalid\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseJSONSchema[any](tt.schema)
			if tt.expectedError == "" && err != nil {
				t.Errorf("expected no error, got %q", err)
			}
			if tt.expectedError != "" && err == nil {
				t.Errorf("expected error %q, got no error", tt.expectedError)
			}
			if tt.expectedError != "" && err != nil && err.Error() != tt.expectedError {
				t.Errorf("expected error %q, got %q", tt.expectedError, err)
			}
		})
	}
}
