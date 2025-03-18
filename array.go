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
	"encoding/json"
	"fmt"
)

// ArraySchema represents a validation schema for array values.
type ArraySchema[T any] struct {
	*Schema[[]T]
	itemValidator func(T) error
}

// Array creates a new validation schema for array values.
func Array[T any]() *ArraySchema[T] {
	return &ArraySchema[T]{
		Schema: New[[]T](),
	}
}

// Items adds a validator for each item in the array.
func (s *ArraySchema[T]) Items(validator func(T) error) *ArraySchema[T] {
	s.itemValidator = validator
	s.validators = append(s.validators, func(arr []T) error {
		for i, item := range arr {
			if err := validator(item); err != nil {
				return fmt.Errorf("invalid item at index %d: %w", i, err)
			}
		}
		return nil
	})
	return s
}

// Min adds a minimum length validator to the schema.
func (s *ArraySchema[T]) Min(min int) *ArraySchema[T] {
	s.validators = append(s.validators, func(arr []T) error {
		if len(arr) < min {
			return fmt.Errorf("array length must be at least %d", min)
		}
		return nil
	})
	return s
}

// Max adds a maximum length validator to the schema.
func (s *ArraySchema[T]) Max(max int) *ArraySchema[T] {
	s.validators = append(s.validators, func(arr []T) error {
		if len(arr) > max {
			return fmt.Errorf("array length must be at most %d", max)
		}
		return nil
	})
	return s
}

// Length adds a validator that checks if the array has exactly the specified length.
func (s *ArraySchema[T]) Length(length int) *ArraySchema[T] {
	s.validators = append(s.validators, func(arr []T) error {
		if len(arr) != length {
			return fmt.Errorf("array length must be exactly %d", length)
		}
		return nil
	})
	return s
}

// UniqueItems adds a validator that checks if all items in the array are unique.
func (s *ArraySchema[T]) UniqueItems() *ArraySchema[T] {
	s.validators = append(s.validators, func(arr []T) error {
		seen := make(map[string]struct{})
		for i, item := range arr {
			// Use JSON marshaling to get a string representation for comparison
			key, err := json.Marshal(item)
			if err != nil {
				return fmt.Errorf("failed to marshal array item for uniqueness check at index %d: %w", i, err)
			}
			keyStr := string(key)
			if _, exists := seen[keyStr]; exists {
				return fmt.Errorf("array items must be unique (duplicate found at index %d)", i)
			}
			seen[keyStr] = struct{}{}
		}
		return nil
	})
	return s
}

// Validate validates the array against the schema and returns an error if the array is not valid.
func (s *ArraySchema[T]) Validate(value []T) error {
	if value == nil {
		// Check if Min validator exists and requires a non-empty array
		for _, validator := range s.validators {
			if err := validator([]T{}); err != nil {
				return err
			}
		}
		return nil
	}
	return s.Schema.Validate(value)
}
