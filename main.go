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
// Version: 1.2.0
//
// Change history:
//    2025-03-12: V1.0.0: Created.
//    2025-04-19: V1.1.0: Restructured password handling.
//    2025-04-20: V1.2.0: Count runes in passwords, not bytes.
//

// This is the main program file.

package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"transposer/encodedfile"
	"transposer/encodinghelper"
	"transposer/filehelper"
	"transposer/logger"
	"transposer/transposition"
)

// ******** Private constants ********

// maxFileSize is the maximum size of a text file.
const maxFileSize = 10_000_000

// minPasswordLen is the minimum length of one password.
const minPasswordLen = 6

// maxPasswordsLen is the maximum length of the string of all passwords.
const maxPasswordsLen = 40_000

// maxNumPasswords is the maximum possible number of passwords.
const maxNumPasswords = 30

// fmtEncodedFileOperation contains the format for encoded file operation messages.
const fmtEncodedFileOperation = `%s file '%s' with %s encoding`

// myVersion contains the version of this program. Update it when anything changes.
const myVersion = `1.2.0`

// mainMsgBase is the base number for program messages.
const mainMsgBase = 10

// myName contains the name of the executable of this program.
var myName = filehelper.RealBaseName(os.Args[0])

// ******** Main functions ********

// main is the main program, but it is only a stub that calls the real main program
// in order to have return codes and still obey defers.
func main() {
	os.Exit(realMain())
}

// realMain is the real main program. It has an exit code and obeys defers.
// This is how a main function should be.
func realMain() int {
	// 1. Process commands.
	defineCommandlineFlags()

	if len(os.Args) < 2 {
		return printUsageError(mainMsgBase+0, `No command specified`)
	}

	doEncrypt, rc, done := checkCommand()
	if done {
		return rc
	}

	err := processCommandlineFlags(doEncrypt)
	if err != nil {
		return printUsageError(mainMsgBase+1, err.Error())
	}

	err = checkCommandlineFlags(doEncrypt)
	if err != nil {
		return printUsageError(mainMsgBase+2, err.Error())
	}

	errorList := NormalizeAndCheckPasswords(passwords)
	if len(errorList) != 0 {
		return printParameterErrorList(mainMsgBase+3, `Password error`, errorList)
	}

	// 2. Process command line arguments.
	var inputEncodingName string
	var outputEncodingName string
	var handleBom BomHandling

	inputEncodingName, outputEncodingName, handleBom, err = GetEncodings(paramFileEncoding)
	if err != nil {
		return printUsageError(mainMsgBase+4, err.Error())
	}

	// 3. Loop through input files.
	for _, inputFilePath := range inputFiles {
		// 3.1 Probe file, whether it has a BOM.
		actInputEncodingName := inputEncodingName
		var fileBomEncodingName string
		var hasBom bool
		fileBomEncodingName, hasBom, err = encodinghelper.ProbeFile(inputFilePath)
		if err != nil {
			return printProcessingError(mainMsgBase+5, err.Error())
		}

		// If it has a BOM and the input encoding is different from it, change the input encoding.
		if len(fileBomEncodingName) != 0 && fileBomEncodingName != inputEncodingName {
			actInputEncodingName = fileBomEncodingName
		}

		// 3.2 Read the input with the specified encoding and transformation options and put the
		//     result in a rune slice.
		logger.PrintInfof(mainMsgBase+6, fmtEncodedFileOperation, `Read`, inputFilePath, actInputEncodingName)
		var inputContent []rune
		var useWindowsLineBreak bool
		inputContent, useWindowsLineBreak, err = encodedfile.ReadRunes(
			inputFilePath,
			maxFileSize,
			actInputEncodingName,
			useConversion,
			toLower,
			paramOnlyLetters)
		if err != nil {
			return printProcessingError(mainMsgBase+7, err.Error())
		}

		// 3.3 Transpose the input runes.
		var resultContent []rune
		if doEncrypt {
			resultContent = transposition.TransposeMultiplePasswords(inputContent, passwords)
		} else {
			resultContent = transposition.UntransposeMultiplePasswords(inputContent, passwords)
		}

		// 3.4 Write the output file with the correct encoding.
		outputFilePath := BuildOutFilePath(doEncrypt, inputFilePath)
		logger.PrintInfof(mainMsgBase+8, fmtEncodedFileOperation, `Write`, outputFilePath, outputEncodingName)
		err = encodedfile.WriteEncoded(
			outputFilePath,
			outputEncodingName,
			bomUsageForOutput(handleBom, hasBom),
			useWindowsLineBreak,
			resultContent)

		if err != nil {
			return printProcessingError(mainMsgBase+9, err.Error())
		}
	}

	return rcOK
}

// ******** Private functions ********

// checkCommand checks the given command and executes it, if it is 'help' or 'version'.
// Otherwise, it only signals if the command was 'encrypt', or not.
func checkCommand() (bool, int, bool) {
	cmd := strings.ToLower(strings.TrimSpace(os.Args[1]))[0]

	if cmd != 'd' && cmd != 'e' && cmd != 'h' && cmd != 'v' {
		return false, printUsageErrorf(mainMsgBase+10, `Invalid command: '%s'`, os.Args[1]), true
	}

	// Print help
	if cmd == 'h' {
		printUsage()
		return false, rcOK, true
	}

	// Print version
	if cmd == 'v' {
		fmt.Printf(`%s V%s (%s, %d cpus)`, myName, myVersion, runtime.Version(), runtime.NumCPU())
		return false, rcOK, true
	}

	// Tell, if it is encrypt, or not.
	doEncrypt := cmd == 'e'

	return doEncrypt, rcOK, false
}
