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

func ExampleArray() {
	// Create an array schema
	schema := valtor.Array[int]()

	// Validate arrays against the schema
	fmt.Println("Validating [1, 2, 3]:", schema.Validate([]int{1, 2, 3}))
	fmt.Println("Validating []:", schema.Validate([]int{}))
	fmt.Println("Validating nil:", schema.Validate(nil))

	// Output:
	// Validating [1, 2, 3]: <nil>
	// Validating []: <nil>
	// Validating nil: <nil>
}

func ExampleArraySchema_Min() {
	// Create an array schema with a minimum length
	schema := valtor.Array[string]().Min(2)

	// Validate arrays of different lengths
	fmt.Println("Validating [\"a\", \"b\", \"c\"]:", schema.Validate([]string{"a", "b", "c"}))
	fmt.Println("Validating [\"a\", \"b\"]:", schema.Validate([]string{"a", "b"}))
	fmt.Println("Validating [\"a\"]:", schema.Validate([]string{"a"}))
	fmt.Println("Validating []:", schema.Validate([]string{}))

	// Output:
	// Validating ["a", "b", "c"]: <nil>
	// Validating ["a", "b"]: <nil>
	// Validating ["a"]: array length must be at least 2
	// Validating []: array length must be at least 2
}

func ExampleArraySchema_Max() {
	// Create an array schema with a maximum length
	schema := valtor.Array[int]().Max(2)

	// Validate arrays of different lengths
	fmt.Println("Validating [1]:", schema.Validate([]int{1}))
	fmt.Println("Validating [1, 2]:", schema.Validate([]int{1, 2}))
	fmt.Println("Validating [1, 2, 3]:", schema.Validate([]int{1, 2, 3}))

	// Output:
	// Validating [1]: <nil>
	// Validating [1, 2]: <nil>
	// Validating [1, 2, 3]: array length must be at most 2
}

func ExampleArraySchema_Length() {
	// Create an array schema with an exact length requirement
	schema := valtor.Array[int]().Length(2)

	// Validate arrays of different lengths
	fmt.Println("Validating [1]:", schema.Validate([]int{1}))
	fmt.Println("Validating [1, 2]:", schema.Validate([]int{1, 2}))
	fmt.Println("Validating [1, 2, 3]:", schema.Validate([]int{1, 2, 3}))

	// Output:
	// Validating [1]: array length must be exactly 2
	// Validating [1, 2]: <nil>
	// Validating [1, 2, 3]: array length must be exactly 2
}

func ExampleArraySchema_UniqueItems() {
	// Create an array schema that requires unique items
	schema := valtor.Array[string]().UniqueItems()

	// Validate arrays with and without duplicate items
	fmt.Println("Validating [\"a\", \"b\", \"c\"]:", schema.Validate([]string{"a", "b", "c"}))
	fmt.Println("Validating [\"a\", \"b\", \"a\"]:", schema.Validate([]string{"a", "b", "a"}))

	// Output:
	// Validating ["a", "b", "c"]: <nil>
	// Validating ["a", "b", "a"]: array items must be unique (duplicate found at index 2)
}

func ExampleArraySchema_Items() {
	// Create an array schema with item validation
	schema := valtor.Array[int]().Items(func(item int) error {
		if item < 0 {
			return fmt.Errorf("item must be non-negative")
		}
		return nil
	})

	// Validate arrays with valid and invalid items
	fmt.Println("Validating [1, 2, 3]:", schema.Validate([]int{1, 2, 3}))
	fmt.Println("Validating [1, -2, 3]:", schema.Validate([]int{1, -2, 3}))

	// Output:
	// Validating [1, 2, 3]: <nil>
	// Validating [1, -2, 3]: invalid item at index 1: item must be non-negative
}

func ExampleArraySchema_multiple_validators() {
	// Create an array schema with multiple validators
	schema := valtor.Array[int]().
		Min(2).
		Max(4).
		Items(func(item int) error {
			if item <= 0 {
				return fmt.Errorf("item must be positive")
			}
			return nil
		}).
		UniqueItems()

	// Validate arrays against multiple validators
	fmt.Println("Valid array:", schema.Validate([]int{1, 2, 3}))
	fmt.Println("Too short:", schema.Validate([]int{1}))
	fmt.Println("Too long:", schema.Validate([]int{1, 2, 3, 4, 5}))
	fmt.Println("Invalid item:", schema.Validate([]int{1, 0, 3}))
	fmt.Println("Duplicate items:", schema.Validate([]int{1, 2, 1}))

	// Output:
	// Valid array: <nil>
	// Too short: array length must be at least 2
	// Too long: array length must be at most 4
	// Invalid item: invalid item at index 1: item must be positive
	// Duplicate items: array items must be unique (duplicate found at index 2)
}
