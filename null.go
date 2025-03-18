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

// NullSchema represents a validation schema for null values.
type NullSchema struct {
	*Schema[any]
}

// Null creates a new validation schema for null values.
func Null() *NullSchema {
	return &NullSchema{
		Schema: New[any](),
	}
}

// Validate validates that the value is null.
func (s *NullSchema) Validate(value any) error {
	if value != nil {
		return fmt.Errorf("expected null value, got %T", value)
	}
	return s.Schema.Validate(value)
}
