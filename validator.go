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

import "errors"

var ErrValueRequired = errors.New("value is required")

// Validator is an interface for validating a value.
// The Validate method is implemented by all validation schemas.
type Validator[T any] interface {
	Validate(value T) error
}

// Schema represents a base type for all validation schemas.
// It implements the Validator interface.
type Schema[T any] struct {
	validators []func(T) error
}

// New creates a new validation schema for type T.
func New[T any]() *Schema[T] {
	return &Schema[T]{
		validators: make([]func(T) error, 0),
	}
}

// Validate runs all validators against the value and returns the first error encountered, if any.
func (s *Schema[T]) Validate(value T) error {
	for _, validator := range s.validators {
		if err := validator(value); err != nil {
			return err
		}
	}
	return nil
}

// Custom adds a custom validation function to the schema and returns the schema for chaining.
func (s *Schema[T]) Custom(fn func(T) error) *Schema[T] {
	s.validators = append(s.validators, fn)
	return s
}
