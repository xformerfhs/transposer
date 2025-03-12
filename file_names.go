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
//    2025-01-02: V1.0.0: Created.
//

// File file_names contains the functions that handle file names in the main program.

package main

import (
	"path/filepath"
	"strings"
	"transposer/filehelper"
)

// ******** Private constants ********

// decryptedMarker is the last part of the default file name for an decrypted file.
const decryptedMarker = `_untrans`

// encryptMarker is the last part of the default file name for an encrypted file.
const encryptMarker = `_trans`

// ******** Public functions ********

func BuildOutFilePath(doEncrypt bool, filePath string) string {
	if doEncrypt {
		return BuildEncryptOutFilePath(filePath)
	} else {
		return BuildDecryptOutFilePath(filePath)
	}
}

// BuildDecryptOutFilePath builds the file path of the decrypted output file.
func BuildDecryptOutFilePath(filePath string) string {
	return buildFilePathWithMarker(filePath, decryptedMarker)
}

// BuildEncryptOutFilePath builds the file path of the encrypted output file.
func BuildEncryptOutFilePath(filePath string) string {
	return buildFilePathWithMarker(filePath, encryptMarker)
}

// ******** Private functions ********

// cleanPathComponents returns the directory, base name and extension
// with the program's markers removed from base name.
func cleanPathComponents(filePath string) (string, string, string) {
	dir, base, ext := filehelper.PathComponents(filePath)

	base = strings.TrimSuffix(base, encryptMarker)
	base = strings.TrimSuffix(base, decryptedMarker)

	return dir, base, ext
}

// buildFilePathWithMarker builds a file path with a marker that is separated by '_' after the base name.
func buildFilePathWithMarker(filePath string, marker string) string {
	dir, base, ext := cleanPathComponents(filePath)

	return filepath.Join(dir, base+marker+ext)
}
