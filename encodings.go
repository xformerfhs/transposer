package main

import (
	"errors"
	"fmt"
	"strings"
	"transposer/encodinghelper"
)

// ******** Private types ********

// BomHandling specifies how BOMs are handled.
type BomHandling byte

// ******** Private constants ********

const (
	// bomMarker means that a BOM must be written.
	bomMarker = `bom`
	// noBomMarker means that a BOM must not be written.
	noBomMarker = `nobom`
)

const (
	// copyBom means that the output file should have the same BOM setting as the input file.
	copyBom BomHandling = iota
	// noBom means that the output file must not have a BOM.
	noBom
	// withBom means that the output file has a BOM.
	withBom
)

// emptyBom is an empty BOM.
var emptyBom = []byte(nil)

// ******** Public functions ********

// GetEncodings analyses the encoding specification and returns the encoding names
// for input and output encoding, the BOM handling to be used and the BOM bytes
// if a BOM should be set on output.
func GetEncodings(fileEncoding string) (string, string, BomHandling, error) {
	elements := strings.Split(fileEncoding, `:`)
	elementsLength := len(elements)

	if elementsLength < 1 {
		return ``, ``, copyBom, errors.New(`no encoding supplied`)
	}
	if elementsLength > 2 {
		return ``, ``, copyBom, errors.New(`more than two encodings supplied`)
	}

	inputEncodingName, handleBom, err := normalizeAndCheckEncodingName(elements[0])
	if err != nil {
		return ``, ``, copyBom, err
	}

	if elementsLength == 1 {
		return inputEncodingName, inputEncodingName, handleBom, nil
	}

	var outputEncodingName string
	outputEncodingName, handleBom, err = normalizeAndCheckEncodingName(elements[1])
	if err != nil {
		return ``, ``, copyBom, err
	}

	return inputEncodingName, outputEncodingName, handleBom, nil
}

// ******** Private functions ********

// normalizeAndCheckEncodingName normalizes the given encoding name, checks, if it is valid
// and returns the BOM information.
func normalizeAndCheckEncodingName(encodingName string) (string, BomHandling, error) {
	normalizedName, handleBom := bomSettings(encodinghelper.NormalizeEncoding(encodingName))
	if !encodinghelper.IsValidEncoding(normalizedName) {
		return ``, copyBom, fmt.Errorf(`invalid encoding: '%s'`, normalizedName)
	}

	return normalizedName, handleBom, nil
}

// bomSettings returns the BOM settings for the given encoding name, cuts off 'bom' and 'nobom'
// suffixes and returns their values as bom handling values.
func bomSettings(encodingName string) (string, BomHandling) {
	handleBom := copyBom

	if strings.HasSuffix(encodingName, noBomMarker) {
		encodingName = encodingName[:len(encodingName)-len(noBomMarker)]
		handleBom = noBom
	} else {
		if strings.HasSuffix(encodingName, bomMarker) {
			encodingName = encodingName[:len(encodingName)-len(bomMarker)]
			handleBom = withBom
		}
	}

	return encodingName, handleBom
}

// bomUsageForOutput returns the correct BOM usage for the output file.
func bomUsageForOutput(handleBom BomHandling, hasBom bool) bool {
	switch handleBom {
	case noBom:
		return false
	case withBom:
		return true
	default:
		return hasBom
	}
}
