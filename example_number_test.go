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

func ExampleNumber() {
	schema := valtor.Number[int]().
		Min(18).
		Max(120)

	err := schema.Validate(25)
	fmt.Println(err)
	err = schema.Validate(15)
	fmt.Println(err)
	err = schema.Validate(150)
	fmt.Println(err)

	// Output:
	// <nil>
	// value must be at least 18
	// value must be at most 120
}

func ExampleNumberSchema_Required() {
	schema := valtor.Number[int]().Required()

	err := schema.Validate(0)
	fmt.Println(err)

	// Output:
	// value is required
}

func ExampleNumberSchema_Min() {
	schema := valtor.Number[float64]().Min(0.5)

	err := schema.Validate(1.0)
	fmt.Println(err)
	err = schema.Validate(0.1)
	fmt.Println(err)

	// Output:
	// <nil>
	// value must be at least 0.5
}

func ExampleNumberSchema_Max() {
	schema := valtor.Number[uint]().Max(100)

	err := schema.Validate(50)
	fmt.Println(err)
	err = schema.Validate(200)
	fmt.Println(err)

	// Output:
	// <nil>
	// value must be at most 100
}

func ExampleNumberSchema_Custom() {
	schema := valtor.Number[int]().Custom(func(n int) error {
		if n < 0 {
			return fmt.Errorf("value must be positive")
		}
		return nil
	})

	err := schema.Validate(10)
	fmt.Println(err)
	err = schema.Validate(-1)
	fmt.Println(err)

	// Output:
	// <nil>
	// value must be positive
}
