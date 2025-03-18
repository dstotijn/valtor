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

package valtor_test

import (
	"fmt"

	"github.com/dstotijn/valtor"
)

func ExampleNull() {
	// Create a null schema
	schema := valtor.Null()

	// Validate nil against the schema
	fmt.Println("Validating nil:", schema.Validate(nil))

	// Validate non-nil values against the schema
	fmt.Println("Validating string:", schema.Validate("hello"))
	fmt.Println("Validating int:", schema.Validate(42))
	fmt.Println("Validating bool:", schema.Validate(true))

	// Output:
	// Validating nil: <nil>
	// Validating string: expected null value, got string
	// Validating int: expected null value, got int
	// Validating bool: expected null value, got bool
}

func ExampleNull_custom() {
	// Create a null schema with a custom validator
	schema := valtor.Null()

	// Add a custom validator (this is just for demonstration, as it doesn't
	// make much practical sense for null values)
	schema.Custom(func(v any) error {
		// Additional validation logic could be added here if needed
		return nil
	})

	// Validate nil against the schema
	fmt.Println("Validating nil:", schema.Validate(nil))

	// Validate non-nil value against the schema
	fmt.Println("Validating string:", schema.Validate("hello"))

	// Output:
	// Validating nil: <nil>
	// Validating string: expected null value, got string
}
