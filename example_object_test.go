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

func ExampleObjectSchema_Field() {
	type User struct {
		Name  string
		Age   int
		Email string
	}

	schema := valtor.Object[User]().
		Field("name",
			func(u User) error {
				return valtor.String().Min(2).Max(50).Validate(u.Name)
			},
		).
		Field("age",
			func(u User) error {
				return valtor.Number[int]().Min(18).Max(120).Validate(u.Age)
			},
		).
		Field("email",
			func(u User) error {
				return valtor.String().Regexp(regexp.MustCompile(`^.+@.+\..+$`)).Validate(u.Email)
			},
		)

	validUser := User{
		Name:  "John Doe",
		Age:   30,
		Email: "john@example.com",
	}

	invalidUser := User{
		Name:  "J",
		Age:   30,
		Email: "john@example.com",
	}

	err := schema.Validate(validUser)
	fmt.Println(err)
	err = schema.Validate(invalidUser)
	fmt.Println(err)

	// Output:
	// <nil>
	// validation failed for field "name": length must be at least 2
}

func ExampleObjectSchema_Field_validateField() {
	type User struct {
		Name  string
		Age   int
		Email string
	}

	schema := valtor.Object[User]().
		Field("name", valtor.ValidateField(
			func(u User) string { return u.Name },
			valtor.String().Min(2).Max(50),
		)).
		Field("age", valtor.ValidateField(
			func(u User) int { return u.Age },
			valtor.Number[int]().Min(18).Max(120),
		)).
		Field("email", valtor.ValidateField(
			func(u User) string { return u.Email },
			valtor.String().Regexp(regexp.MustCompile(`^.+@.+\..+$`)),
		))

	validUser := User{
		Name:  "John Doe",
		Age:   30,
		Email: "john@example.com",
	}

	invalidUser := User{
		Name:  "J",
		Age:   30,
		Email: "john@example.com",
	}

	err := schema.Validate(validUser)
	fmt.Println(err)
	err = schema.Validate(invalidUser)
	fmt.Println(err)

	// Output:
	// <nil>
	// validation failed for field "name": length must be at least 2
}

func ExampleObjectSchema_Map() {
	type Baz struct {
		Quo string
	}

	type Foo struct {
		Bar string
		Baz Baz
	}

	schema := valtor.Object[Foo]().Map(valtor.FieldValidatorMap[Foo]{
		"bar": func(u Foo) error {
			return valtor.String().Max(5).Validate(u.Bar)
		},
		"baz": func(u Foo) error {
			return valtor.Object[Baz]().Map(valtor.FieldValidatorMap[Baz]{
				"quo": func(a Baz) error {
					return valtor.String().Max(5).Validate(a.Quo)
				},
			}).Validate(u.Baz)
		},
	})

	validUser := Foo{
		Bar: "bar",
		Baz: Baz{
			Quo: "quo",
		},
	}

	invalidUser := Foo{
		Bar: "bar",
		Baz: Baz{
			Quo: "quoquo",
		},
	}

	err := schema.Validate(validUser)
	fmt.Println(err)
	err = schema.Validate(invalidUser)
	fmt.Println(err)

	// Output:
	// <nil>
	// validation failed for field "baz": validation failed for field "quo": length must be at most 5
}
