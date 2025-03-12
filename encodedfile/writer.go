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
// Version: 1.0.0
//
// Change history:
//    2025-03-12: V1.0.0: Created.
//

package encodedfile

import (
	"bufio"
	"golang.org/x/text/transform"
	"os"
	"transposer/encodinghelper"
	"transposer/filehelper"
)

func WriteEncoded(
	path string,
	encodingName string,
	useBom bool,
	useWindowsLineBreaks bool,
	data []rune) error {
	bw, f, err := prepareFileFoWriting(path, encodingName, useBom)
	if err != nil {
		return err
	}
	defer filehelper.CloseWithName(f)

	for _, r := range data {
		if useWindowsLineBreaks && r == lineFeed {
			err = bw.WriteByte(carriageReturn)
			if err != nil {
				return err
			}
		}

		_, err = bw.WriteRune(r)
		if err != nil {
			return err
		}
	}

	return bw.Flush()
}

// ******** Private functions ********

// prepareFileFoWriting prepares the given file path for writing.
func prepareFileFoWriting(path string, encodingName string, useBom bool) (*bufio.Writer, *os.File, error) {
	var f *os.File
	f, err := os.Create(path)
	if err != nil {
		return nil, nil, err
	}

	encoded, _, _ := encodinghelper.TranslateEncoding(encodingName, useBom)

	return bufio.NewWriter(transform.NewWriter(f, encoded.NewEncoder())), f, nil
}
