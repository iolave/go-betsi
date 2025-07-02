package anyutils

import (
	"reflect"

	"github.com/theothertomelliott/acyclic"
)

// DestroyCircular masks any circular references in the given value
// with it's zero value.
//
//   - If the value is a known type, it will be returned as-is.
//   - CAUTION: this function will modify the value in place and if
//     modifying the value is not intended, pass a copy of v.
//
// TODO: make a recursive valudation to not just check the first level
// and thus, loose as few as possible.
func DestroyCircular(v any) any {
	switch v := v.(type) {
	case string, *string,
		int, *int,
		int8, *int8,
		int16, *int16,
		int32, *int32,
		int64, *int64,
		uint, *uint,
		uint8, *uint8,
		uint16, *uint16,
		uint32, *uint32,
		uint64, *uint64,
		float32, *float32,
		float64, *float64,
		bool, *bool,
		[]string, []*string,
		[]int, []*int,
		[]int8, []*int8,
		[]int16, []*int16,
		[]int32, []*int32,
		[]int64, []*int64,
		[]uint, []*uint,
		[]uint8, []*uint8,
		[]uint16, []*uint16,
		[]uint32, []*uint32,
		[]uint64, []*uint64,
		[]float32, []*float32,
		[]float64, []*float64,
		[]bool, []*bool:
		return v
	default:
		vtype := reflect.TypeOf(v)
		vvalue := reflect.ValueOf(v)

		switch vtype.Kind() {
		case reflect.Slice:
			for i := 0; i < vvalue.Len(); i++ {
				vvalue := vvalue.Index(i)
				v := vvalue.Interface()

				if err := acyclic.Check(v); err == nil {
					continue
				}

				if vvalue.CanSet() {
					vvalue.SetZero()
				}
			}

			return v
		case reflect.Map:
			iter := vvalue.MapRange()
			for iter.Next() {
				key := iter.Key()
				vvalueKey := iter.Value()
				vvalueKeyVal := vvalueKey.Interface()

				if err := acyclic.Check(vvalueKeyVal); err == nil {
					continue
				}

				vvalue.SetMapIndex(key, reflect.Value{})
			}

			return v
		case reflect.Ptr, reflect.Struct:
			if err := acyclic.Check(v); err == nil {
				return v
			}

			return nil
		default:
			return v
		}

	}
}
