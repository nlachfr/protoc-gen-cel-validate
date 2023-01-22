package validate

import (
	"regexp"
	"strings"
	"testing"

	"github.com/nlachfr/protocel/testdata/validate"
)

func TestGenerateResourceTypePatternMapping(t *testing.T) {
	mapping := GenerateResourceTypePatternMapping(validate.File_testdata_validate_test_proto)
	regexpMapping := map[string]*regexp.Regexp{}
	for k, v := range mapping {
		r, err := regexp.Compile(strings.ReplaceAll(v, `\\`, `\`))
		if err != nil {
			t.Errorf("regexp compile error: %v", err)
		}
		regexpMapping[k] = r
	}
	tests := []struct {
		Name    string
		Type    string
		Value   string
		WantErr bool
	}{
		{
			Name:    "Ref (OK)",
			Type:    "testdata/Ref",
			Value:   "refs/myRef",
			WantErr: false,
		},
		{
			Name:    "Ref (NOK)",
			Type:    "testdata/Ref",
			Value:   "ref/myRef",
			WantErr: true,
		},
		{
			Name:    "RefMultiple Refs (OK)",
			Type:    "testdata/RefMultiple",
			Value:   "refs/myRef",
			WantErr: false,
		},
		{
			Name:    "RefMultiple Multiple (OK)",
			Type:    "testdata/RefMultiple",
			Value:   "multiples/my/refs/ref",
			WantErr: false,
		},
		{
			Name:    "RefMultiple Others (OK)",
			Type:    "testdata/RefMultiple",
			Value:   "others/refs/myRef",
			WantErr: false,
		},
		{
			Name:    "RefMultiple (NOK)",
			Type:    "testdata/RefMultiple",
			Value:   "nok/refs",
			WantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			r, ok := regexpMapping[tt.Type]
			if !ok {
				t.Errorf("unknown type: %v", tt.Type)
			} else {
				got := r.MatchString(tt.Value)
				if got == tt.WantErr {
					t.Errorf("want %v, got %v", tt.WantErr, got)
				}
			}
		})
	}
}
