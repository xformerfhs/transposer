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

// This is the main program file.

package main

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"transposer/encodedfile"
	"transposer/encodinghelper"
	"transposer/filehelper"
	"transposer/logger"
	"transposer/mathhelper"
	"transposer/stringhelper"
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
const maxNumPasswords = 100

// fmtEncodedFileOperation contains the format for encoded file operation messages.
const fmtEncodedFileOperation = `%s file '%s' with %s encoding`

// myVersion contains the version of this program. Update it, when anything changes.
const myVersion = `1.0.0`

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

// realMain is the real main program.
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

	// 2. Process command line arguments.
	if !checkPasswords(passwords) {
		return rcParameterError
	}

	var inputEncodingName string
	var outputEncodingName string
	var handleBom BomHandling

	inputEncodingName, outputEncodingName, handleBom, err = GetEncodings(fileEncoding)
	if err != nil {
		return printUsageError(mainMsgBase+3, err.Error())
	}

	// 3. Loop through input files.
	for _, inputFilePath := range inputFiles {
		// 3.1 Probe file, whether it has a BOM.
		actInputEncodingName := inputEncodingName
		var fileBomEncodingName string
		var hasBom bool
		fileBomEncodingName, hasBom, err = encodinghelper.ProbeFile(inputFilePath)
		if err != nil {
			return printProcessingError(mainMsgBase+4, err.Error())
		}

		// If it has a BOM and the input encoding is different from it, change the input encoding.
		if len(fileBomEncodingName) != 0 && fileBomEncodingName != inputEncodingName {
			actInputEncodingName = fileBomEncodingName
		}

		// 3.2 Read the input with the specified encoding and transformation options and put the
		//     result in a rune slice.
		logger.PrintInfof(mainMsgBase+5, fmtEncodedFileOperation, `Read`, inputFilePath, actInputEncodingName)
		var inputContent []rune
		var useWindowsLineBreak bool
		inputContent, useWindowsLineBreak, err = encodedfile.ReadRunes(
			inputFilePath,
			maxFileSize,
			actInputEncodingName,
			useConversion,
			toLower,
			onlyLetters)
		if err != nil {
			return printProcessingError(mainMsgBase+6, err.Error())
		}

		// 3.3 Transpose the input runes.
		var resultContent []rune
		if doEncrypt {
			resultContent = transposition.TransposeRuneArray(inputContent, passwords)
		} else {
			resultContent = transposition.UnTransposeRuneArray(inputContent, passwords)
		}

		// 3.4 Write the output file with the correct encoding.
		outputFilePath := BuildOutFilePath(doEncrypt, inputFilePath)
		logger.PrintInfof(mainMsgBase+7, fmtEncodedFileOperation, `Write`, outputFilePath, outputEncodingName)
		err = encodedfile.WriteEncoded(
			outputFilePath,
			outputEncodingName,
			bomUsageForOutput(handleBom, hasBom),
			useWindowsLineBreak,
			resultContent)

		if err != nil {
			return printProcessingError(mainMsgBase+8, err.Error())
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
		return false, printUsageErrorf(mainMsgBase+9, `Invalid command: '%s'`, os.Args[1]), true
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

// getPasswords analyzes the passwordSpec string and returns all given passwords as a slice.
// The passwords are converted to all lower-case. This is done to prevent weakening of the
// passwords by starting them with an upper-case letter and then continuing with lower
// case letters, which would effectively remove one letter from the password strength.
func getPasswords(password string) ([]string, error) {
	elements := stringhelper.SplitAnyN(password, `:,`, maxNumPasswords+1)
	if len(elements) > maxNumPasswords {
		return nil, errors.New(`too many passwords`)
	}

	result := make([]string, len(elements))
	for i, candidate := range elements {
		if !stringhelper.IsAlphaNumeric(candidate) {
			return nil, errors.New(`passwords must only contain alphanumeric characters`)
		}
		result[i] = strings.ToLower(candidate)
	}

	return result, nil
}

// checkPasswords checks the lengths of the passwords for a common divisor
// which weakens the encryption.
func checkPasswords(passwords []string) bool {
	result := true

	// 1. Get lengths of passwords
	pwLen := len(passwords)
	lengths := make([]int, pwLen)
	for i, pw := range passwords {
		pwl := len(pw)
		if pwl < minPasswordLen {
			logger.PrintErrorf(mainMsgBase+10, `Password '%s' is too short`, pw)
			return false
		}

		lengths[i] = pwl
	}

	logger.PrintInfof(mainMsgBase+11, `Password lengths: %v`, lengths)

	// 2. Check password lengths for common divisors.
	for i := 0; i < pwLen-1; i++ {
		li := lengths[i]
		for j := i + 1; j < pwLen; j++ {
			lj := lengths[j]
			gcd := mathhelper.Gcd(li, lj)
			if gcd != 1 {
				logger.PrintErrorf(mainMsgBase+12,
					`The lengths of the passwords '%s' (%d) and '%s' (%d) share a common divisor: %d`,
					passwords[i],
					li,
					passwords[j],
					lj,
					gcd)
				result = false
			}
		}
	}

	return result
}
