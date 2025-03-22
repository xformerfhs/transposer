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
	{`Nö sepäratör ät öll`, `:,-.`, -1, []string{`Nö sepäratör ät öll`}},
	{`Nö-sépäràtor,ät öll`, `:,-.`, -1, []string{`Nö`, `sépäràtor`, `ät öll`}},
	{`:Nö-sèpárätör,ät öll`, `:,-.`, -1, []string{``, `Nö`, `sèpárätör`, `ät öll`}},
	{`.Nô-sépàrätór,ät öll:`, `:,-.`, -1, []string{``, `Nô`, `sépàrätór`, `ät öll`, ``}},
	{`.Shört#`, ``, -1, []string{`.`, `S`, `h`, `ö`, `r`, `t`, `#`}},
	{`Whàtéver`, `$%?`, 0, nil},
	{`Nô-sépàrätör,ât äll`, `:,-.`, 2, []string{`Nô`, `sépàrätör,ât äll`}},
	{`%Nó$sèpärátòr?ät_âll`, `ä%?`, 1_234_567, []string{``, `Nó$sèp`, `rátòr`, ``, `t_âll`}},
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
