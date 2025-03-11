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
	"regexp"
)

// StringSchema represents a validation schema for string values.
type StringSchema struct {
	*Schema[string]
	required bool
}

// String creates a new validation schema for string values.
func String() *StringSchema {
	return &StringSchema{
		Schema: New[string](),
	}
}

// Required will make a string value required to be not empty when validated.
func (s *StringSchema) Required() *StringSchema {
	s.required = true
	return s
}

// Min adds a minimum length validator to the schema and returns the schema for chaining.
func (s *StringSchema) Min(min int) *StringSchema {
	s.validators = append(s.validators, func(v string) error {
		if len(v) < min {
			return fmt.Errorf("length must be at least %d", min)
		}
		return nil
	})
	return s
}

// Max adds a maximum length validator to the schema and returns the schema for chaining.
func (s *StringSchema) Max(max int) *StringSchema {
	s.validators = append(s.validators, func(v string) error {
		if len(v) > max {
			return fmt.Errorf("length must be at most %d", max)
		}
		return nil
	})
	return s
}

// Length adds a length validator to the schema and returns the schema for chaining.
func (s *StringSchema) Length(length int) *StringSchema {
	s.validators = append(s.validators, func(v string) error {
		if len(v) != length {
			return fmt.Errorf("length must be exactly %d", length)
		}
		return nil
	})
	return s
}

// Regexp adds a regular expression pattern validator to the schema and returns the schema for chaining.
func (s *StringSchema) Regexp(re *regexp.Regexp) *StringSchema {
	s.validators = append(s.validators, func(v string) error {
		if !re.MatchString(v) {
			return fmt.Errorf("string must match pattern %q", re.String())
		}
		return nil
	})
	return s
}

// Validate validates the string against the schema and returns an error if the string is not valid.
func (s *StringSchema) Validate(value string) error {
	if value == "" && s.required {
		return ErrValueRequired
	}
	return s.Schema.Validate(value)
}
