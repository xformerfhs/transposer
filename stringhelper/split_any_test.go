package stringhelper_test

import (
	"slices"
	"testing"
	"transposer/stringhelper"
)

func TestSplitAny(t *testing.T) {
	testString := ``
	separators := `:,-.`
	result := stringhelper.SplitAny(testString, separators)
	if !slices.Equal(result, []string{``}) {
		t.Fatalf(`Splitting an empty string got wrong result: %v`, result)
	}

	testString = `No separator at all`
	result = stringhelper.SplitAny(testString, separators)
	if !slices.Equal(result, []string{`No separator at all`}) {
		t.Fatalf(`Splitting a string without a separator got wrong result: %v`, result)
	}

	testString = `No-separator,at all`
	result = stringhelper.SplitAny(testString, separators)
	if !slices.Equal(result, []string{`No`, `separator`, `at all`}) {
		t.Fatalf(`Splitting a string with multiple separators got wrong result: %v`, result)
	}

	testString = `:No-separator,at all`
	result = stringhelper.SplitAny(testString, separators)
	if !slices.Equal(result, []string{``, `No`, `separator`, `at all`}) {
		t.Fatalf(`Splitting a string beginning with a separator got wrong result: %v`, result)
	}

	testString = `No-separator,at all.`
	result = stringhelper.SplitAny(testString, separators)
	if !slices.Equal(result, []string{`No`, `separator`, `at all`, ``}) {
		t.Fatalf(`Splitting a string ending with a separator got wrong result: %v`, result)
	}

	testString = `.No-separator,at all:`
	result = stringhelper.SplitAny(testString, separators)
	if !slices.Equal(result, []string{``, `No`, `separator`, `at all`, ``}) {
		t.Fatalf(`Splitting a string beginning and ending with a separator got wrong result: %v`, result)
	}

	testString = `.Short#`
	result = stringhelper.SplitAny(testString, ``)
	if !slices.Equal(result, []string{`.`, `S`, `h`, `o`, `r`, `t`, `#`}) {
		t.Fatalf(`Splitting a string with an empty separator got wrong result: %v`, result)
	}

	result = stringhelper.SplitAny(``, ``)
	if !slices.Equal(result, []string(nil)) {
		t.Fatalf(`Splitting an empty string with an empty separator got wrong result: %v`, result)
	}
}
