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

// PointerSchema represents a validation schema for pointer values.
type PointerSchema[T any] struct {
	*Schema[*T]
	required bool
}

// Pointer creates a new validation schema for pointer values.
func Pointer[T any]() *PointerSchema[T] {
	return &PointerSchema[T]{
		Schema: New[*T](),
	}
}

// Required will make a pointer value required to not be nil when validated.
func (s *PointerSchema[T]) Required() *PointerSchema[T] {
	s.required = true
	return s
}

// NotNil adds a validation that ensures the pointer is not nil.
// This is an alias for Required() for more explicit validation chains.
func (s *PointerSchema[T]) NotNil() *PointerSchema[T] {
	return s.Required()
}

// Custom adds a custom validation function to the schema and returns the schema for chaining.
func (s *PointerSchema[T]) Custom(fn func(*T) error) *PointerSchema[T] {
	s.Schema.Custom(fn)
	return s
}

// Validate validates the pointer against the schema and returns an error if the pointer is not valid.
func (s *PointerSchema[T]) Validate(value *T) error {
	if value == nil && s.required {
		return ErrValueRequired
	}
	return s.Schema.Validate(value)
}

// Ptr wraps another validator schema to validate the pointed-to value.
func Ptr[T any](schema Validator[T]) *PointerSchema[T] {
	p := Pointer[T]()
	p.Custom(func(value *T) error {
		if value == nil {
			// Skip validation for nil pointers, handled by Required() if needed.
			return nil
		}
		return schema.Validate(*value)
	})
	return p
}
