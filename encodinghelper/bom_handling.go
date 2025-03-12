//
// SPDX-FileCopyrightText: Copyright 2024 Frank Schwab
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
// Version: 3.0.0
//
// Change history:
//    2024-03-10: V1.0.0: Created.
//    2025-01-19: V1.1.0: Correct handling of short files.
//    2025-03-10: V2.0.0: Return encoding name and BOM bytes.
//    2025-03-12: V3.0.0: Adapted to UTF BOM handling functions.
//

package encodinghelper

import (
	"golang.org/x/text/encoding"
	"os"
	"transposer/filehelper"
)

// ******** Public functions ********

// ProbeFile reads the first bytes of a file to check for BOMs.
// If it finds one, it returns the corresponding encoding name
// and the BOM bytes.
func ProbeFile(fileName string) (string, bool, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return ``, false, err
	}
	defer filehelper.CloseWithName(f)

	// 1. Read the first three bytes.
	var readCount int
	miniBuffer := make([]byte, 3)
	readCount, err = f.Read(miniBuffer)
	if err != nil {
		return ``, false, err
	}

	// 3. Check read bytes.

	// File has only 1 byte. There is no BOM.
	if readCount < 2 {
		return ``, false, nil
	}

	if miniBuffer[0] == utf16BeBom[0] &&
		miniBuffer[1] == utf16BeBom[1] {
		return `utf16be`, true, nil
	}

	if miniBuffer[0] == utf16LeBom[0] &&
		miniBuffer[1] == utf16LeBom[1] {
		return `utf16le`, true, nil
	}

	// File has only 2 bytes. There is no UTF-8-BOM.
	if readCount < 3 {
		return ``, false, nil
	}

	if miniBuffer[0] == utf8Bom[0] &&
		miniBuffer[1] == utf8Bom[1] &&
		miniBuffer[2] == utf8Bom[2] {
		return `utf8`, true, nil
	}

	return ``, false, nil
}

// EncodingWithBomHandling returns an encoding with the required BOM handling.
// If the encoding does not use a BOM, just the encoding is returned.
func EncodingWithBomHandling(encodingName string, useBom bool) encoding.Encoding {
	switch encodingName {
	case `utf8`:
		if useBom {
			return utf8BomEncoding
		} else {
			return utf8NoBomEncoding
		}

	case `utf16le`:
		if useBom {
			return utf16leBomEncoding
		} else {
			return utf16leNoBomEncoding
		}

	case `utf16be`:
		if useBom {
			return utf16beBomEncoding
		} else {
			return utf16beNoBomEncoding
		}

	default:
		return textToEncoding[encodingName].encoding
	}
}
