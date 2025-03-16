//
// SPDX-FileCopyrightText: Copyright 2024-2025 Frank Schwab
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
//    2024-03-10: V1.0.0: Created.
//

package encodinghelper

import (
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/unicode"
)

// ******** Private types ********

// encodingInfo holds the name and the encoding.Encoding of an encoding.
type encodingInfo struct {
	name     string
	encoding encoding.Encoding
}

// ******** Private constants ********

// utf16BeBom contains the bytes of a UTF-16BE BOM.
var utf16BeBom = []byte{0xfe, 0xff}

// utf16LeBom contains the bytes of a UTF-16LE BOM.
var utf16LeBom = []byte{0xff, 0xfe}

// utf8Bom contains the bytes of a UTF-8 BOM.
var utf8Bom = []byte{0xef, 0xbb, 0xbf}

// utf16beBomEncoding is the encoding for UTF-16BE with BOM read and write.
var utf16beBomEncoding = unicode.UTF16(unicode.BigEndian, unicode.UseBOM)

// utf16beNoBomEncoding is the encoding for UTF-16BE without BOM read and write.
var utf16beNoBomEncoding = unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)

// utf16leBomEncoding is the encoding for UTF-16LE with BOM read and write.
var utf16leBomEncoding = unicode.UTF16(unicode.LittleEndian, unicode.UseBOM)

// utf16leNoBomEncoding is the encoding for UTF-16LE without BOM read and write.
var utf16leNoBomEncoding = unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)

// utf8BomEncoding is the encoding for UTF-8 with BOM read and write.
var utf8BomEncoding = unicode.UTF8BOM

// utf8NoBomEncoding is the encoding for UTF-8 without BOM read and write.
var utf8NoBomEncoding = unicode.UTF8

// textToEncoding maps an encoding specification to the corresponding encoding information.
var textToEncoding = make(map[string]encodingInfo, 50)

// ******** Init function *********

// init initializes encoding variables.
func init() {
	fillEncodingMap()
}

// ******** Private functions *********

// fillEncodingMap fills the encoding map.
func fillEncodingMap() {
	textToEncoding[`utf8`] = encodingInfo{name: `UTF-8`, encoding: nil}
	textToEncoding[`utf16be`] = encodingInfo{name: `UTF-16BE`, encoding: nil}
	textToEncoding[`utf16le`] = encodingInfo{name: `UTF-16LE`, encoding: nil}

	for _, enc := range charmap.All {
		cm, isCm := enc.(*charmap.Charmap)
		if isCm {
			charMapName := cm.String()
			textToEncoding[NormalizeEncoding(charMapName)] = encodingInfo{name: charMapName, encoding: enc}
		}
	}
}
