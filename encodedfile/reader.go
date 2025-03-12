package encodedfile

import (
	"bufio"
	"errors"
	"fmt"
	"golang.org/x/text/transform"
	"io"
	"os"
	"transposer/encodinghelper"
	"transposer/filehelper"
	"transposer/numberformat"
	"unicode"
)

func ReadAllRunesAsIs(
	path string,
	maxFileSize int64,
	encodingName string) ([]rune, bool, error) {
	result, br, f, err := prepareFileForReading(path, maxFileSize, encodingName)
	if err != nil {
		return nil, false, err
	}
	defer filehelper.CloseWithName(f)

	hadCarriageReturn := false
	hasWindowsLineBreaks := false
	i := 0
	for {
		var r rune
		r, _, err = br.ReadRune()

		if errors.Is(err, io.EOF) {
			return result[:i], hasWindowsLineBreaks, nil
		}

		if err != nil {
			return nil, false, err
		}

		// Do not read carriage return if followed by new line.
		if r != carriageReturn {
			// But read it, if it is not followed by new line.
			if hadCarriageReturn {
				if r != lineFeed {
					result[i] = carriageReturn
					i++
				} else {
					hasWindowsLineBreaks = true
				}
			}

			result[i] = r
			i++
			hadCarriageReturn = false
		} else {
			hadCarriageReturn = true
		}
	}
}

func ReadOnlyLetterRunes(
	path string,
	maxFileSize int64,
	encodingName string,
	changeCase bool,
	useLower bool) ([]rune, bool, error) {
	result, br, f, err := prepareFileForReading(path, maxFileSize, encodingName)
	if err != nil {
		return nil, false, err
	}
	defer filehelper.CloseWithName(f)

	i := 0
	for {
		var r rune
		r, _, err = br.ReadRune()

		if errors.Is(err, io.EOF) {
			return result[:i], false, nil
		}

		if err != nil {
			return nil, false, err
		}

		if unicode.IsLetter(r) {
			if changeCase {
				if useLower {
					r = unicode.ToLower(r)
				} else {
					r = unicode.ToUpper(r)
				}
			}

			result[i] = r
			i++
		}
	}
}

func ReadAllRunes(
	path string,
	maxFileSize int64,
	encodingName string,
	changeCase bool,
	useLower bool) ([]rune, bool, error) {
	result, br, f, err := prepareFileForReading(path, maxFileSize, encodingName)
	if err != nil {
		return nil, false, err
	}
	defer filehelper.CloseWithName(f)

	hadCarriageReturn := false
	hasWindowsLineBreaks := false
	i := 0
	for {
		var r rune
		r, _, err = br.ReadRune()

		if errors.Is(err, io.EOF) {
			return result[:i], hasWindowsLineBreaks, nil
		}

		if err != nil {
			return nil, false, err
		}

		// Do not read carriage return if followed by new line.
		if r != carriageReturn {
			// But read it, if it is not followed by new line.
			if hadCarriageReturn {
				if r != lineFeed {
					result[i] = carriageReturn
					i++
				} else {
					hasWindowsLineBreaks = true
				}
			}

			if changeCase && unicode.IsLetter(r) {
				if useLower {
					r = unicode.ToLower(r)
				} else {
					r = unicode.ToUpper(r)
				}
			}

			result[i] = r
			i++
			hadCarriageReturn = false
		} else {
			hadCarriageReturn = true
		}
	}
}

// ******** Private functions ********

// prepareFileForReading checks the file size and sets up the reader for the given file path.
func prepareFileForReading(path string, maxFileSize int64, encodingName string) ([]rune, *bufio.Reader, *os.File, error) {
	fileSize, err := checkSize(path, maxFileSize)
	if err != nil {
		return nil, nil, nil, err
	}

	var f *os.File
	f, err = os.Open(path)
	if err != nil {
		return nil, nil, nil, err
	}

	result := make([]rune, fileSize)

	// The encoding for reading must always cope with BOMs, when applicable.
	encoded, _, _ := encodinghelper.TranslateEncoding(encodingName, true)

	return result, bufio.NewReader(transform.NewReader(f, encoded.NewDecoder())), f, nil
}

// checkSize checks if the given file is not greater than the maximum file size.
func checkSize(path string, maxFileSize int64) (int64, error) {
	fs, err := filehelper.FileSize(path)
	if err != nil {
		return 0, fmt.Errorf(`error getting size of file '%s': %v`, path, err)
	}

	if fs > maxFileSize {
		return 0, fmt.Errorf(`size of file '%s' is larger than maximum file size %s`, path, numberformat.FormatInt64(maxFileSize))
	}

	return fs, nil
}
