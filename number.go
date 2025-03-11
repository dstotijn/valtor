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

import "fmt"

// NumberSchema represents a validation schema for numeric values.
type NumberSchema[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64] struct {
	*Schema[T]
	required bool
}

// Number creates a new validation schema for numeric values.
func Number[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64]() *NumberSchema[T] {
	return &NumberSchema[T]{
		Schema: New[T](),
	}
}

// Required will make a number value required to not be the empty value when validated.
func (s *NumberSchema[T]) Required() *NumberSchema[T] {
	s.required = true
	return s
}

// Validate validates the number against the schema and returns an error if the number is not valid.
func (s *NumberSchema[T]) Validate(value T) error {
	var zero T
	if value == zero && s.required {
		return ErrValueRequired
	}
	return s.Schema.Validate(value)
}

// Min adds a minimum value validator to the schema and returns the schema for chaining.
func (s *NumberSchema[T]) Min(min T) *NumberSchema[T] {
	s.validators = append(s.validators, func(v T) error {
		if v < min {
			return fmt.Errorf("value must be at least %v", min)
		}
		return nil
	})
	return s
}

// Max adds a maximum value validator to the schema and returns the schema for chaining.
func (s *NumberSchema[T]) Max(max T) *NumberSchema[T] {
	s.validators = append(s.validators, func(v T) error {
		if v > max {
			return fmt.Errorf("value must be at most %v", max)
		}
		return nil
	})
	return s
}
