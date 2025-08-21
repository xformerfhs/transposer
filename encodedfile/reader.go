//
// SPDX-FileCopyrightText: Copyright 2025 Frank Schwab
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileType: SOURCE
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
//
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Author: Frank Schwab
//
// Version: 3.1.0
//
// Change history:
//    2025-03-12: V1.0.0: Created.
//    2025-03-14: V2.0.0: Replaced different Read functions with just one with filter and converter.
//    2025-03-14: V3.0.0: Removed filter and converter and use flags, as it is faster.
//    2025-08-21: V3.1.0: Do not waste memory with UTF-16 encodings.
//

package encodedfile

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"transposer/encodinghelper"
	"transposer/filehelper"
	"transposer/numberformat"
	"unicode"

	"golang.org/x/text/transform"
)

// ReadRunes reads runes from the given file and filters and converts the runes.
// The function returns the filtered and converted runes and a flag indicating
// if the file contained Windows line breaks.
// If the file size is larger than the maximum file size, the function returns an error.
// If the file cannot be opened or read, the function returns an error.
// The encoding of the file is given.
// If doConversion is true, the runes are converted to upper case or lower case, depending on the value of toLower.
// If onlyLetters is true, only letters are read.
func ReadRunes(
	path string,
	maxFileSize int64,
	encodingName string,
	doConversion bool,
	toLower bool,
	onlyLetters bool) ([]rune, bool, error) {
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

		// Do not read carriage return if followed by a new line.
		if r != carriageReturn {
			// But read it if it is not followed by a new line.
			if hadCarriageReturn {
				if r != lineFeed {
					result[i] = carriageReturn
					i++
				} else {
					hasWindowsLineBreaks = true
				}
			}

			hadCarriageReturn = false

			if onlyLetters && !unicode.IsLetter(r) {
				continue
			}

			if doConversion {
				if toLower {
					r = unicode.ToLower(r)
				} else {
					r = unicode.ToUpper(r)
				}
			}

			result[i] = r
			i++
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

	resultSize := int(fileSize)
	if strings.HasPrefix(encodingName, `utf16`) {
		resultSize >>= 1
	}
	result := make([]rune, resultSize)

	// The encoding for reading must always cope with BOMs, when applicable.
	encoder, _, _ := encodinghelper.EncodingForName(encodingName, true)

	return result, bufio.NewReader(transform.NewReader(f, encoder.NewDecoder())), f, nil
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
