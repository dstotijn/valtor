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

func ExampleBool() {
	// Create a boolean schema
	schema := valtor.Bool()

	// Validate true and false against the schema
	fmt.Println("Validating true:", schema.Validate(true))
	fmt.Println("Validating false:", schema.Validate(false))

	// Output:
	// Validating true: <nil>
	// Validating false: <nil>
}

func ExampleBool_mustBeTrue() {
	// Create a boolean schema that requires the value to be true
	schema := valtor.Bool().MustBeTrue()

	// Validate true and false against the schema
	fmt.Println("Validating true:", schema.Validate(true))
	fmt.Println("Validating false:", schema.Validate(false))

	// Output:
	// Validating true: <nil>
	// Validating false: bool value must be true
}

func ExampleBool_mustBeFalse() {
	// Create a boolean schema that requires the value to be false
	schema := valtor.Bool().MustBeFalse()

	// Validate true and false against the schema
	fmt.Println("Validating true:", schema.Validate(true))
	fmt.Println("Validating false:", schema.Validate(false))

	// Output:
	// Validating true: bool value must be false
	// Validating false: <nil>
}
