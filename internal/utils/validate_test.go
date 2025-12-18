package utils

import "testing"

func TestValidateRecursively(t *testing.T) {
	type TestCase struct {
		Name    string
		In      any
		WantErr bool
	}

	testCases := []TestCase{}
	testCases = append(testCases, TestCase{
		Name:    "validation should pass for nil",
		In:      nil,
		WantErr: false,
	})
	testCases = append(testCases, TestCase{
		Name: "validation should pass for struct",
		In: struct {
			Test string `validate:"required"`
		}{"test"},
		WantErr: false,
	})
	testCases = append(testCases, TestCase{
		Name: "validation should fail for struct",
		In: struct {
			Test string `validate:"required"`
		}{""},
		WantErr: true,
	})
	testCases = append(testCases, TestCase{
		Name: "validation should pass for *struct",
		In: &struct {
			Test string `validate:"required"`
		}{"test"},
		WantErr: false,
	})
	testCases = append(testCases, TestCase{
		Name: "validation should fail for *struct",
		In: &struct {
			Test string `validate:"required"`
		}{""},
		WantErr: true,
	})
	testCases = append(testCases, TestCase{
		Name: "validation should pass for slice",
		In: []struct {
			Test string `validate:"required"`
		}{
			{"test"},
		},
		WantErr: false,
	})
	testCases = append(testCases, TestCase{
		Name: "validation should fail for slice",
		In: []struct {
			Test string `validate:"required"`
		}{
			{""},
		},
		WantErr: true,
	})
	testCases = append(testCases, TestCase{
		Name: "validation should pass for *slice",
		In: &[]struct {
			Test string `validate:"required"`
		}{
			{"test"},
		},
		WantErr: false,
	})
	testCases = append(testCases, TestCase{
		Name: "validation should fail for *slice",
		In: &[]struct {
			Test string `validate:"required"`
		}{
			{""},
		},
		WantErr: true,
	})
	testCases = append(testCases, TestCase{
		Name: "validation should pass for map",
		In: map[string]struct {
			Test string `validate:"required"`
		}{
			"test": {Test: "test"},
		},
		WantErr: false,
	})
	testCases = append(testCases, TestCase{
		Name: "validation should fail for map",
		In: map[string]struct {
			Test string `validate:"required"`
		}{
			"test": {Test: ""},
		},
		WantErr: true,
	})
	testCases = append(testCases, TestCase{
		Name: "validation should pass for *map",
		In: &map[string]struct {
			Test string `validate:"required"`
		}{
			"test": {Test: "test"},
		},
		WantErr: false,
	})
	testCases = append(testCases, TestCase{
		Name: "validation should fail for *map",
		In: &map[string]struct {
			Test string `validate:"required"`
		}{
			"test": {Test: ""},
		},
		WantErr: true,
	})
	testCases = append(testCases, TestCase{
		Name: "validation should pass for struct with pointer",
		In: struct {
			Test *struct {
				Test string `validate:"required"`
			} `validate:"required"`
		}{
			Test: &struct {
				Test string `validate:"required"`
			}{Test: "test"},
		},
		WantErr: false,
	})
	testCases = append(testCases, TestCase{
		Name: "validation should fail for struct with pointer",
		In: struct {
			Test *struct {
				Test string `validate:"required"`
			} `validate:"required"`
		}{
			Test: &struct {
				Test string `validate:"required"`
			}{Test: ""},
		},
		WantErr: true,
	})
	testCases = append(testCases, TestCase{
		Name: "validation should pass for *struct with pointer",
		In: &struct {
			Test *struct {
				Test string `validate:"required"`
			} `validate:"required"`
		}{
			Test: &struct {
				Test string `validate:"required"`
			}{Test: "test"},
		},
		WantErr: false,
	})
	testCases = append(testCases, TestCase{
		Name: "validation should fail for *struct with pointer",
		In: &struct {
			Test *struct {
				Test string `validate:"required"`
			} `validate:"required"`
		}{
			Test: &struct {
				Test string `validate:"required"`
			}{Test: ""},
		},
		WantErr: true,
	})

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			err := ValidateRecursively(tt.In)
			if err != nil && !tt.WantErr {
				t.Errorf("got error %v, wantErr: %v", err, tt.WantErr)
			}
		})
	}
}
