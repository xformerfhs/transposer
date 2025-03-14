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
// Version: 2.0.0
//
// Change history:
//    2025-03-12: V1.0.0: Created.
//    2025-03-14: V2.0.0: Replaced different Read functions with just one with filter and converter.
//

package encodedfile

import (
	"bufio"
	"errors"
	"fmt"
	"golang.org/x/text/transform"
	"io"
	"os"
	"transposer/converters"
	"transposer/encodinghelper"
	"transposer/filehelper"
	"transposer/filters"
	"transposer/numberformat"
)

// ReadRunes reads runes from the given file and filters and converts the runes.
func ReadRunes(
	path string,
	maxFileSize int64,
	encodingName string,
	filter filters.RuneFilter,
	converter converters.RuneConverter) ([]rune, bool, error) {
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

			hadCarriageReturn = false

			if filter(r) {
				result[i] = converter(r)
				i++
			}
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
