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

package valtor

import (
	"fmt"
)

// ObjectSchema represents a validation schema for object values.
type ObjectSchema[T any] struct {
	*Schema[T]
	fieldValidators map[string]func(any) error
}

// FieldValidatorMap is a type alias for a map of field names to validator functions.
type FieldValidatorMap[T any] map[string]func(T) error

// Object creates a new validation schema for object values.
func Object[T any]() *ObjectSchema[T] {
	return &ObjectSchema[T]{
		Schema:          New[T](),
		fieldValidators: make(map[string]func(any) error),
	}
}

// Field adds a field validator to the schema and returns the schema for chaining.
func (s *ObjectSchema[T]) Field(fieldName string, validateFn func(T) error) *ObjectSchema[T] {
	s.fieldValidators[fieldName] = func(value any) error {
		// Test whether the value is of type T, else use its zero value (which
		// could be nil, and should be handled by the validator).
		typedValue, _ := value.(T)

		if err := validateFn(typedValue); err != nil {
			return fmt.Errorf("validation failed for field %q: %w", fieldName, err)
		}
		return nil
	}
	return s
}

// ValidateField is a helper function to create a field validator.
func ValidateField[T any, F any](getter func(T) F, schema Validator[F]) func(T) error {
	return func(value T) error {
		return schema.Validate(getter(value))
	}
}

// Map adds multiple field validators to the schema at once using a map.
func (s *ObjectSchema[T]) Map(fieldValidators FieldValidatorMap[T]) *ObjectSchema[T] {
	for fieldName, validateFn := range fieldValidators {
		s.Field(fieldName, validateFn)
	}
	return s
}

// Validate validates a value against the schema.
func (s *ObjectSchema[T]) Validate(value T) error {
	mapValue, ok := any(value).(map[string]any)
	if ok {
		return s.ValidateMap(mapValue)
	}
	for _, validator := range s.fieldValidators {
		if err := validator(value); err != nil {
			return err
		}
	}
	return nil
}

// ValidateMap validates a map (keyed by field name) of values against the schema.
func (s *ObjectSchema[T]) ValidateMap(values map[string]any) error {
	for fieldName, validateFn := range s.fieldValidators {
		value := values[fieldName]
		if err := validateFn(value); err != nil {
			return err
		}
	}
	return nil
}
