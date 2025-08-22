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
// Version: 1.5.0
//
// Change history:
//    2025-03-12: V1.0.0: Created.
//    2025-04-19: V1.1.0: Restructured password handling.
//    2025-04-20: V1.2.0: Count runes in passwords, not bytes.
//    2025-04-27: V1.3.0: Transposer is an object that no longer needs the passwords.
//    2025-08-13: V1.3.1: Simplified transposition.
//    2025-08-21: V1.4.0: Handle UTF-32, less waste of memory with UTF-16.
//    2025-08-21: V1.5.0: Refactored realMain with the help of JetBrains AI (GPT-5).
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
const maxFileSize = 100_000_000

// minPasswordLen is the minimum length of one password.
const minPasswordLen = 6

// maxPasswordsLen is the maximum length of the string of all passwords.
const maxPasswordsLen = 40_000

// maxNumPasswords is the maximum possible number of passwords.
const maxNumPasswords = 30

// fmtEncodedFileOperation contains the format for encoded file operation messages.
const fmtEncodedFileOperation = `%s file '%s' with %s encoding`

// myVersion contains the version of this program. Update it when anything changes.
const myVersion = `1.5.0`

// mainMsgBase is the base number for program messages.
const mainMsgBase = 10

// myName contains the name of the executable of this program.
var myName = filehelper.RealBaseName(os.Args[0])

// ******** Private structures ********

// appConfig contains the configuration of the program.
// It is built from the command line arguments and the configuration file.
// It is used to build the transposer and to process the input files.
type appConfig struct {
	doEncrypt          bool
	inputEncodingName  string
	outputEncodingName string
	handleBom          BomHandling
}

// ******** Private interfaces ********

// runeTransposer is an interface for transposition of runes.
type runeTransposer interface {
	Transpose([]rune) []rune
	Untranspose([]rune) []rune
}

// ******** Main functions ********

// main is the main program, but it is only a stub that calls the real main program
// in order to have return codes and still obey defers.
func main() {
	os.Exit(realMain())
}

// realMain is the real main program. It has an exit code and obeys defers.
// This is how a main function should be.
func realMain() int {
	// 1. Parse, validate, and build config.
	cfg, rc, done := parseAndValidate()
	if done {
		return rc
	}

	// 2. Prepare transposer once and ensure cleanup.
	transposer := transposition.New[rune](passwords)
	passwords = nil // Passwords is no longer needed after the previous statement.
	defer transposer.Destroy()

	// 3. Process all input files.
	for _, inputFilePath := range inputFiles {
		if rc = processFile(transposer, cfg, inputFilePath); rc != rcOK {
			return rc
		}
	}

	return rcOK
}

// ******** Private functions ********

// parseAndValidate parses the command line arguments, validates them, and builds the config.
func parseAndValidate() (*appConfig, int, bool) {
	defineCommandlineFlags()

	result := &appConfig{}

	if len(os.Args) < 2 {
		return result, printUsageError(mainMsgBase+0, `No command specified`), true
	}

	doEncrypt, rc, done := checkCommand()
	if done {
		return result, rc, true
	}
	result.doEncrypt = doEncrypt

	var err error
	if err = processCommandlineFlags(doEncrypt); err != nil {
		return result, printUsageError(mainMsgBase+1, err.Error()), true
	}

	if err = checkCommandlineFlags(doEncrypt); err != nil {
		return result, printUsageError(mainMsgBase+2, err.Error()), true
	}

	errorList := NormalizeAndCheckPasswords(passwords)
	if len(errorList) != 0 {
		return result, printParameterErrorList(mainMsgBase+3, `Password error`, errorList), true
	}

	result.inputEncodingName,
		result.outputEncodingName,
		result.handleBom,
		err = GetEncodings(paramFileEncoding)
	if err != nil {
		return result, printUsageError(mainMsgBase+4, err.Error()), true
	}

	return result, rcOK, false
}

// processFile processes one input file.
func processFile(t runeTransposer, cfg *appConfig, inputFilePath string) int {
	// 3.1 Probe file, whether it has a BOM.
	actInputEncodingName := cfg.inputEncodingName
	fileBomEncodingName, hasBom, err := encodinghelper.ProbeFile(inputFilePath)
	if err != nil {
		return printProcessingError(mainMsgBase+5, err.Error())
	}

	if hasBom {
		// If it has a BOM and the input encoding is UTF-32, return an error.
		if strings.HasPrefix(fileBomEncodingName, `utf32`) {
			return printProcessingError(mainMsgBase+6, `UTF-32 is not supported`)
		}

		// If it has a BOM and the input encoding is different from it, change the input encoding.
		if fileBomEncodingName != cfg.inputEncodingName {
			actInputEncodingName = fileBomEncodingName
		}
	}

	// 3.2 Read the input with the specified encoding and transformation options and put the
	//     result in a rune slice.
	logger.PrintInfof(mainMsgBase+7, fmtEncodedFileOperation, `Read`, inputFilePath, actInputEncodingName)
	inputContent, useWindowsLineBreak, err := encodedfile.ReadRunes(
		inputFilePath,
		maxFileSize,
		actInputEncodingName,
		useConversion,
		toLower,
		paramOnlyLetters)
	if err != nil {
		return printProcessingError(mainMsgBase+8, err.Error())
	}

	// 3.3 Transpose the input runes.
	var resultContent []rune
	if cfg.doEncrypt {
		resultContent = t.Transpose(inputContent)
	} else {
		resultContent = t.Untranspose(inputContent)
	}

	// 3.4 Write the output file with the correct encoding.
	outputFilePath := BuildOutFilePath(cfg.doEncrypt, inputFilePath)
	logger.PrintInfof(mainMsgBase+9, fmtEncodedFileOperation, `Write`, outputFilePath, cfg.outputEncodingName)
	err = encodedfile.WriteEncoded(
		outputFilePath,
		cfg.outputEncodingName,
		bomUsageForOutput(cfg.handleBom, hasBom),
		useWindowsLineBreak,
		resultContent)
	if err != nil {
		return printProcessingError(mainMsgBase+10, err.Error())
	}

	return rcOK
}

// checkCommand checks the given command and executes it if it is 'help' or 'version'.
// Otherwise, it only signals if the command was 'encrypt', or not.
func checkCommand() (bool, int, bool) {
	cmd := strings.ToLower(strings.TrimSpace(os.Args[1]))[0]

	if cmd != 'd' && cmd != 'e' && cmd != 'h' && cmd != 'v' {
		return false, printUsageErrorf(mainMsgBase+11, `Invalid command: '%s'`, os.Args[1]), true
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
