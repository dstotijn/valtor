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

// BoolSchema represents a validation schema for boolean values.
type BoolSchema struct {
	*Schema[bool]
}

// Bool creates a new validation schema for boolean values.
func Bool() *BoolSchema {
	return &BoolSchema{
		Schema: New[bool](),
	}
}

// Validate validates the boolean against the schema and returns an error if the boolean is not valid.
func (s *BoolSchema) Validate(value bool) error {
	return s.Schema.Validate(value)
}

// MustBeTrue adds a validator that checks if the boolean value is true.
func (s *BoolSchema) MustBeTrue() *BoolSchema {
	s.validators = append(s.validators, func(v bool) error {
		if !v {
			return fmt.Errorf("bool value must be true")
		}
		return nil
	})
	return s
}

// MustBeFalse adds a validator that checks if the boolean value is false.
func (s *BoolSchema) MustBeFalse() *BoolSchema {
	s.validators = append(s.validators, func(v bool) error {
		if v {
			return fmt.Errorf("bool value must be false")
		}
		return nil
	})
	return s
}
