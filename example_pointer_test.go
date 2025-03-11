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

func ExamplePointer() {
	schema := valtor.Pointer[string]().Required()

	var nilPtr *string
	str := "hello"

	err := schema.Validate(&str)
	fmt.Println(err)
	err = schema.Validate(nilPtr)
	fmt.Println(err)

	// Output:
	// <nil>
	// value is required
}

func ExamplePointerSchema_Custom() {
	schema := valtor.Pointer[int]().Custom(func(n *int) error {
		if n != nil && *n < 0 {
			return fmt.Errorf("value must be positive")
		}
		return nil
	})

	positive := 10
	negative := -1

	err := schema.Validate(&positive)
	fmt.Println(err)
	err = schema.Validate(&negative)
	fmt.Println(err)

	// Output:
	// <nil>
	// value must be positive
}

func ExamplePtr() {
	// Create a string validator that requires length >= 3.
	stringSchema := valtor.String().Min(3)

	// Wrap it in a pointer validator.
	// Note: We're *not* chaining it with `Required()` here.
	schema := valtor.Ptr(stringSchema)

	longStr := "hello"
	shortStr := "hi"
	var nilPtr *string

	err := schema.Validate(&longStr)
	fmt.Println(err)
	err = schema.Validate(&shortStr)
	fmt.Println(err)
	err = schema.Validate(nilPtr)
	fmt.Println(err)

	// Make it required (non-nil)
	schema = schema.Required()
	err = schema.Validate(nilPtr)
	fmt.Println(err)

	// Output:
	// <nil>
	// length must be at least 3
	// <nil>
	// value is required
}
