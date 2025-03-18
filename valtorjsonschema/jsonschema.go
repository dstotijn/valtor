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

package valtorjsonschema

import (
	"errors"
	"fmt"
	"log"
	"math"
	"regexp"
	"slices"

	"github.com/dstotijn/valtor"
	"github.com/invopop/jsonschema"
)

var ErrInvalidType = errors.New("invalid type")

func ParseJSONSchema[T any](schema jsonschema.Schema) (*valtor.Schema[T], error) {
	return parseJSONSchema[T](schema, false)
}

func parseJSONSchema[T any](schema jsonschema.Schema, required bool) (*valtor.Schema[T], error) {
	switch schema.Type {
	case "null":
		nullSchema := valtor.Null()

		return valtor.New[T]().Custom(func(value T) error {
			return nullSchema.Validate(value)
		}), nil
	case "boolean":
		boolSchema := valtor.Bool()

		return valtor.New[T]().Custom(func(value T) error {
			switch v := any(value).(type) {
			case bool:
				return boolSchema.Validate(v)
			case nil:
				if required {
					return valtor.ErrValueRequired
				}
				return nil
			default:
				return fmt.Errorf("expected boolean value, got %T", v)
			}
		}), nil
	case "array":
		if schema.Items == nil {
			arrSchema := valtor.Array[any]()

			if schema.MinItems != nil {
				arrSchema.Min(int(*schema.MinItems))
			}

			if schema.MaxItems != nil {
				arrSchema.Max(int(*schema.MaxItems))
			}

			if schema.UniqueItems {
				arrSchema.UniqueItems()
			}

			return valtor.New[T]().Custom(func(value T) error {
				switch v := any(value).(type) {
				case []any:
					return arrSchema.Validate(v)
				case nil:
					if required && schema.MinItems != nil && *schema.MinItems > 0 {
						return valtor.ErrValueRequired
					}
					return nil
				default:
					return fmt.Errorf("expected array value, got %T", v)
				}
			}), nil
		}

		itemSchema, err := parseJSONSchema[any](*schema.Items, false)
		if err != nil {
			return nil, fmt.Errorf("invalid item schema: %w", err)
		}

		arrSchema := valtor.Array[any]().Items(itemSchema.Validate)

		if schema.MinItems != nil {
			arrSchema.Min(int(*schema.MinItems))
		}

		if schema.MaxItems != nil {
			arrSchema.Max(int(*schema.MaxItems))
		}

		if schema.UniqueItems {
			arrSchema.UniqueItems()
		}

		return valtor.New[T]().Custom(func(value T) error {
			switch v := any(value).(type) {
			case []any:
				return arrSchema.Validate(v)
			case nil:
				if required && schema.MinItems != nil && *schema.MinItems > 0 {
					return valtor.ErrValueRequired
				}
				return nil
			default:
				return fmt.Errorf("expected array value, got %T", v)
			}
		}), nil
	case "string":
		strSchema := valtor.String()

		if schema.MinLength != nil {
			strSchema.Min(int(*schema.MinLength))
		}
		if schema.MaxLength != nil {
			strSchema.Max(int(*schema.MaxLength))
		}
		if schema.Pattern != "" {
			re, err := regexp.Compile(schema.Pattern)
			if err != nil {
				return nil, fmt.Errorf("invalid pattern %q: %w", schema.Pattern, err)
			}
			strSchema.Regexp(re)
		}

		if required {
			strSchema = strSchema.Required()
		}

		return valtor.New[T]().Custom(func(value T) error {
			switch typedValue := any(value).(type) {
			case string:
				return strSchema.Validate(typedValue)
			case nil:
				return strSchema.Validate("")
			default:
				return fmt.Errorf("expected string value, got %T", value)
			}
		}), nil
	case "integer":
		numSchema := valtor.Number[int64]()

		if min := schema.Minimum; min != "" {
			minFloat, err := min.Float64()
			if err != nil {
				return nil, fmt.Errorf("invalid `minimum` value %q", min)
			}
			minInt := int64(math.Ceil(minFloat))
			numSchema.Min(minInt)
		}
		if max := schema.Maximum; max != "" {
			maxFloat, err := max.Float64()
			if err != nil {
				return nil, fmt.Errorf("invalid `maximum` value %q", max)
			}
			maxInt := int64(math.Floor(maxFloat))
			numSchema.Max(maxInt)
		}

		if required {
			numSchema = numSchema.Required()
		}

		return valtor.New[T]().Custom(func(value T) error {
			switch typedValue := any(value).(type) {
			case int64:
				return numSchema.Validate(typedValue)
			case int32:
				return numSchema.Validate(int64(typedValue))
			case int16:
				return numSchema.Validate(int64(typedValue))
			case int8:
				return numSchema.Validate(int64(typedValue))
			case int:
				return numSchema.Validate(int64(typedValue))
			case uint64:
				if typedValue > math.MaxInt64 {
					return fmt.Errorf("uint64 value %d exceeds maximum int64", typedValue)
				}
				return numSchema.Validate(int64(typedValue))
			case uint32:
				return numSchema.Validate(int64(typedValue))
			case uint16:
				return numSchema.Validate(int64(typedValue))
			case uint8:
				return numSchema.Validate(int64(typedValue))
			case uint:
				if uint64(typedValue) > math.MaxInt64 {
					return fmt.Errorf("uint value %d exceeds maximum int64", typedValue)
				}
				return numSchema.Validate(int64(typedValue))
			case float64:
				if typedValue != math.Trunc(typedValue) {
					return fmt.Errorf("expected integer value, got float with fractional part: %v", typedValue)
				}
				if typedValue > math.MaxInt64 || typedValue < math.MinInt64 {
					return fmt.Errorf("float value %v exceeds int64 range", typedValue)
				}
				return numSchema.Validate(int64(typedValue))
			case nil:
				return numSchema.Validate(0)
			default:
				log.Printf("expected integer value, got %T", typedValue)
				return fmt.Errorf("expected integer value, got %T", typedValue)
			}
		}), nil

	case "number":
		numSchema := valtor.Number[float64]()

		if min := schema.Minimum; min != "" {
			minFloat, err := min.Float64()
			if err != nil {
				return nil, fmt.Errorf("invalid `minimum` %q: %w", min, err)
			}
			numSchema.Min(minFloat)
		}
		if max := schema.Maximum; max != "" {
			maxFloat, err := max.Float64()
			if err != nil {
				return nil, fmt.Errorf("invalid `maximum` %q: %w", max, err)
			}
			numSchema.Max(maxFloat)
		}

		if required {
			numSchema = numSchema.Required()
		}

		return valtor.New[T]().Custom(func(value T) error {
			switch typedValue := any(value).(type) {
			case float64:
				return numSchema.Validate(typedValue)
			case int64:
				return numSchema.Validate(float64(typedValue))
			case int32:
				return numSchema.Validate(float64(typedValue))
			case int16:
				return numSchema.Validate(float64(typedValue))
			case int8:
				return numSchema.Validate(float64(typedValue))
			case int:
				return numSchema.Validate(float64(typedValue))
			case uint64:
				return numSchema.Validate(float64(typedValue))
			case uint32:
				return numSchema.Validate(float64(typedValue))
			case uint16:
				return numSchema.Validate(float64(typedValue))
			case uint8:
				return numSchema.Validate(float64(typedValue))
			case uint:
				return numSchema.Validate(float64(typedValue))
			case nil:
				return numSchema.Validate(0)
			default:
				return fmt.Errorf("expected numeric value, got %T", typedValue)
			}
		}), nil
	case "object":
		objSchema := valtor.Object[any]()

		for pair := schema.Properties.Oldest(); pair != nil; pair = pair.Next() {
			if pair.Value == nil {
				continue
			}

			fieldRequired := false
			if slices.Contains(schema.Required, pair.Key) {
				fieldRequired = true
			}

			fieldSchema, err := parseJSONSchema[any](*pair.Value, fieldRequired)
			if err != nil {
				return nil, fmt.Errorf("invalid schema for property %q: %w", pair.Key, err)
			}

			objSchema.Field(pair.Key, fieldSchema.Validate)
		}

		return valtor.New[T]().Custom(func(value T) error {
			return objSchema.Validate(value)
		}), nil
	case "":
		fallthrough
	default:
		return nil, ErrInvalidType
	}
}
