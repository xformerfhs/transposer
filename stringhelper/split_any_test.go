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
	count     int
}

var splitAnyTests = []SplitAnyTest{
	{``, ``, -1, []string{}, 1},
	{``, `:,-.`, -1, []string{``}, 0},
	{`Nö sepäratör ät öll`, `:,-.`, -1, []string{`Nö sepäratör ät öll`}, 0},
	{`Nö-sépäràtor,ät öll`, `:,-.`, -1, []string{`Nö`, `sépäràtor`, `ät öll`}, 2},
	{`:Nö-sèpárätör,ät öll`, `:,-.`, -1, []string{``, `Nö`, `sèpárätör`, `ät öll`}, 3},
	{`.Nô-sépàrätór,ät öll:`, `:,-.`, -1, []string{``, `Nô`, `sépàrätór`, `ät öll`, ``}, 4},
	{`.Shört#`, ``, -1, []string{`.`, `S`, `h`, `ö`, `r`, `t`, `#`}, 8},
	{`Whàtéver`, `$%?`, 0, nil, 0},
	{`Nô-sépàrätör,ât äll`, `:,-.`, 2, []string{`Nô`, `sépàrätör,ât äll`}, 2},
	{`%Nó$sèpärátòr?ät_âll`, `ä%?`, 1_234_567, []string{``, `Nó$sèp`, `rátòr`, ``, `t_âll`}, 4},
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

func TestCountAny(t *testing.T) {
	for _, tt := range splitAnyTests {
		result := stringhelper.CountAny(tt.source, tt.separator)
		if result != tt.count {
			t.Errorf("CountAny(%q, %q) = %v; want %v", tt.source, tt.separator, result, tt.count)
			continue
		}
	}
}
