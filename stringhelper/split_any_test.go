package stringhelper_test

import (
	"slices"
	"testing"
	"transposer/stringhelper"
)

type SplitAnyTest struct {
	source    string
	separator string
	n         int
	expected  []string
}

var splitAnyTests = []SplitAnyTest{
	{``, ``, -1, []string{}},
	{``, `:,-.`, -1, []string{``}},
	{`No separator at all`, `:,-.`, -1, []string{`No separator at all`}},
	{`No-separator,at all`, `:,-.`, -1, []string{`No`, `separator`, `at all`}},
	{`:No-separator,at all`, `:,-.`, -1, []string{``, `No`, `separator`, `at all`}},
	{`.No-separator,at all:`, `:,-.`, -1, []string{``, `No`, `separator`, `at all`, ``}},
	{`.Short#`, ``, -1, []string{`.`, `S`, `h`, `o`, `r`, `t`, `#`}},
	{`Whatever`, `$%?`, 0, nil},
	{`No-separator,at all`, `:,-.`, 2, []string{`No`, `separator,at all`}},
	{`%No$separator?at all`, `$%?`, 1_234_567, []string{``, `No`, `separator`, `at all`}},
}

func TestSplitAny(t *testing.T) {
	for _, tt := range splitAnyTests {
		result := stringhelper.SplitAnyN(tt.source, tt.separator, tt.n)
		if !slices.Equal(result, tt.expected) {
			t.Errorf("SplitAny(%q, %q, %d) = %v; want %v", tt.source, tt.separator, tt.n, result, tt.expected)
			continue
		}
	}
}
