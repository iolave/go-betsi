package anyutils

import (
	"reflect"
	"testing"
)

type First struct {
	Data   string
	Second *Second
}

type Second struct {
	First *First
}

func TestDestroyCircular(t *testing.T) {
	t.Run("should return the same value using simple types", func(t *testing.T) {
		type testStruct struct {
			v any
			w any
		}

		tests := []testStruct{
			{v: 1, w: 1},
			{v: "str", w: "str"},
		}

		for _, tt := range tests {
			if got := DestroyCircular(tt.v); tt.v != got {
				t.Errorf("DestroyCircular() = %v, want %v", got, tt.w)
			}
		}
	})

	t.Run("should return the same value using slices", func(t *testing.T) {
		type testStruct struct {
			v any
			w any
		}

		tests := []testStruct{
			{v: []int{1, 2, 3}, w: []int{1, 2, 3}},
			{v: []any{"str1", 1, false}, w: []any{"str1", 1, false}},
		}

		for _, tt := range tests {
			if got := DestroyCircular(tt.v); !reflect.DeepEqual(got, tt.w) {
				t.Errorf("DestroyCircular() = %v, want %v", got, tt.w)
			}
		}
	})

	t.Run("should zero out circular slices items", func(t *testing.T) {
		type testStruct struct {
			v any
			w any
		}

		first := &First{Data: "str"}
		second := &Second{}
		first.Second = second
		second.First = first

		test := testStruct{
			v: []any{first, "string"},
			w: []any{nil, "string"},
		}

		if got := DestroyCircular(test.v); !reflect.DeepEqual(got, test.w) {
			t.Errorf("DestroyCircular() = %v, want %v", got, test.w)
		}
	})

	t.Run("should zero out circular map items", func(t *testing.T) {
		type testStruct struct {
			v any
			w any
		}

		first := &First{Data: "str"}
		second := &Second{}
		first.Second = second
		second.First = first

		test := testStruct{
			v: map[string]any{"circular": first, "string": "string"},
			w: map[string]any{"string": "string"},
		}

		if got := DestroyCircular(test.v); !reflect.DeepEqual(got, test.w) {
			t.Errorf("DestroyCircular() = %v, want %v", got, test.w)
		}
	})

	t.Run("should return original pointer when it's not a circular", func(t *testing.T) {
		type testStruct struct {
			v any
			w any
		}

		first := &First{Data: "str"}

		test := testStruct{
			v: first,
			w: first,
		}

		if got := DestroyCircular(test.v); !reflect.DeepEqual(got, test.w) {
			t.Errorf("DestroyCircular() = %v, want %v", got, test.w)
		}
	})

	t.Run("should nil out circular structs", func(t *testing.T) {
		type testStruct struct {
			v any
			w any
		}

		first := &First{Data: "str"}
		second := &Second{}
		first.Second = second
		second.First = first

		test := testStruct{
			v: first,
			w: nil,
		}

		if got := DestroyCircular(test.v); got != test.w {
			t.Errorf("DestroyCircular() = %v, want %v", got, test.w)
		}
	})

	t.Run("should return the same value when type is not identified", func(t *testing.T) {
		fn := func() {}
		type testStruct struct {
			v any
			w any
		}

		test := testStruct{
			v: fn,
			w: fn,
		}

		if got := DestroyCircular(test.v); reflect.ValueOf(got).Pointer() != reflect.ValueOf(test.w).Pointer() {
			t.Errorf("DestroyCircular() = %v, want %v", got, test.w)
		}
	})
}
