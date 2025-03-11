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
	"regexp"

	"github.com/dstotijn/valtor"
)

func ExampleString() {
	schema := valtor.String().
		Max(8).
		Regexp(regexp.MustCompile(`^[a-zA-Z]*$`))

	err := schema.Validate("hello")
	fmt.Println(err)
	err = schema.Validate("")
	fmt.Println(err)
	err = schema.Validate("hellohello")
	fmt.Println(err)
	err = schema.Validate("hello123")
	fmt.Println(err)

	// Output:
	// <nil>
	// <nil>
	// length must be at most 8
	// string must match pattern "^[a-zA-Z]*$"
}

func ExampleStringSchema_Required() {
	schema := valtor.String().Required()

	err := schema.Validate("")
	fmt.Println(err)

	// Output:
	// value is required
}

func ExampleStringSchema_Min() {
	schema := valtor.String().Min(3)

	err := schema.Validate("hello")
	fmt.Println(err)
	err = schema.Validate("hi")
	fmt.Println(err)

	// Output:
	// <nil>
	// length must be at least 3
}

func ExampleStringSchema_Max() {
	schema := valtor.String().Max(5)

	err := schema.Validate("hello")
	fmt.Println(err)
	err = schema.Validate("too long")
	fmt.Println(err)

	// Output:
	// <nil>
	// length must be at most 5
}

func ExampleStringSchema_Length() {
	schema := valtor.String().Length(5)

	err := schema.Validate("hello")
	fmt.Println(err)
	err = schema.Validate("too long")
	fmt.Println(err)

	// Output:
	// <nil>
	// length must be exactly 5
}

func ExampleStringSchema_Regexp() {
	schema := valtor.String().Regexp(regexp.MustCompile(`^[a-z]+$`))

	err := schema.Validate("hello")
	fmt.Println(err)
	err = schema.Validate("Hello123")
	fmt.Println(err)

	// Output:
	// <nil>
	// string must match pattern "^[a-z]+$"
}

func ExampleStringSchema_Custom() {
	schema := valtor.String().Custom(func(s string) error {
		if s == "hello" {
			return nil
		}
		return fmt.Errorf("invalid string")
	})

	err := schema.Validate("hello")
	fmt.Println(err)
	err = schema.Validate("world")
	fmt.Println(err)

	// Output:
	// <nil>
	// invalid string
}
