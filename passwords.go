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
//    2025-03-28: V1.0.0: Created.
//    2025-04-19: V2.0.0: Restructured.
//

// File passwords contains all functions dealing with passwords.

package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"transposer/filehelper"
	"transposer/mathhelper"
	"transposer/stringhelper"
)

// ******** Private constants ********

// -------- Errors --------

// ErrMultiplePasswordSources is returned when multiple password sources where specified.
var ErrMultiplePasswordSources = errors.New(`multiple password sources specified`)

// ErrNoPasswordSources is returned when no password sources where specified.
var ErrNoPasswordSources = errors.New(`no passwords supplied`)

// ErrPasswordsTooLong is returned when the byte count of the passwords is too long.
var ErrPasswordsTooLong = errors.New(`passwords are too long`)

// ErrTooManyPasswords is returned when there are more passwords than allowed.
var ErrTooManyPasswords = errors.New(`too many passwords`)

// ErrInvalidPasswordFile is returned when a password file is not valid.
var ErrInvalidPasswordFile = errors.New(`password file is empty or not a normal file`)

// ******** Public functions ********

// GetPasswords gets the passwords from either the passwords text or the passwords file.
func GetPasswords() ([]string, error) {
	listLen := len(paramPasswordsText)
	if listLen != 0 {
		if len(paramPasswordsFile) != 0 {
			return nil, ErrMultiplePasswordSources
		}

		if listLen > maxPasswordsLen {
			return nil, ErrPasswordsTooLong
		}

		return readPasswordsFromText(paramPasswordsText)
	} else {
		if len(paramPasswordsFile) != 0 {
			return readPasswordsFromFile(paramPasswordsFile)
		} else {
			return nil, ErrNoPasswordSources
		}
	}
}

// NormalizeAndCheckPasswords normalizes and validates a list of passwords,
// returning errors for invalid passwords.
func NormalizeAndCheckPasswords(passwords []string) []error {
	lengths, errorList := normalizePasswords(passwords)
	if len(errorList) != 0 {
		return errorList
	}

	return checkPasswordLengths(passwords, lengths)
}

// ******** Private functions ********

// readPasswordsFromText converts a password text into a slice of passwords, ensuring a maximum number of them.
// Returns the passwords as a slice of strings or an error if the constraint is violated.
func readPasswordsFromText(passwordsText string) ([]string, error) {
	elements := stringhelper.SplitAnyN(passwordsText, `:,`, maxNumPasswords+1)
	if len(elements) > maxNumPasswords {
		return nil, ErrTooManyPasswords
	}

	return elements, nil
}

// readPasswordsFromFile reads the passwords from a file.
func readPasswordsFromFile(passwordsFileName string) ([]string, error) {
	size, err := filehelper.FileSize(passwordsFileName)
	if err != nil {
		return nil, err
	}

	if size == 0 {
		return nil, ErrInvalidPasswordFile
	}

	if size > maxPasswordsLen+32 {
		return nil, ErrPasswordsTooLong
	}

	var file *os.File
	file, err = os.Open(passwordsFileName)
	if err != nil {
		return nil, err
	}
	defer filehelper.CloseWithName(file)

	// Create a new scanner to read the file line by line.
	scanner := bufio.NewScanner(file)

	result := make([]string, 0, maxNumPasswords)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}

		result = append(result, line)
	}

	// Check for errors during the scan
	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// normalizePasswords checks the given list of passwords for invalid characters and normalizes them.
// It also checks the lengths of the passwords.
// Returns a slice of errors for all passwords that are invalid.
// The passwords are normalized by converting them to lower case. This is done to prevent weakening
// of the passwords by starting them with an upper-case letter and then continuing with lower
// case letters, which would effectively remove one letter from the password strength.
func normalizePasswords(passwords []string) ([]int, []error) {
	errorList := make([]error, 0, len(passwords))
	lengths := make([]int, len(passwords))
	for i, password := range passwords {
		pwLen := len(password)
		lengths[i] = pwLen
		if pwLen < minPasswordLen {
			errorList = append(errorList, fmt.Errorf(`password '%s' is too short`, password))
			continue
		}

		if !stringhelper.IsAlphaNumeric(password) {
			errorList = append(errorList, fmt.Errorf(`password '%s' is not alphanumeric`, password))
		}

		passwords[i] = strings.ToLower(password)
	}

	return lengths, errorList
}

// checkPasswordLengths checks the lengths of the passwords for a common divisor
// which weakens the encryption.
func checkPasswordLengths(passwords []string, lengths []int) []error {
	result := make([]error, 0, len(passwords))
	pwLen := len(passwords)

	for i := 0; i < pwLen-1; i++ {
		li := lengths[i]
		for j := i + 1; j < pwLen; j++ {
			lj := lengths[j]
			gcd := mathhelper.Gcd(li, lj)
			if gcd != 1 {
				result = append(result,
					fmt.Errorf(
						`the lengths of the passwords '%s' (%d) and '%s' (%d) share a common divisor: %d`,
						passwords[i],
						li,
						passwords[j],
						lj,
						gcd))
			}
		}
	}

	return result
}
