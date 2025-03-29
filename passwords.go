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
//    2025-03-28: V1.0.0: Created.
//

// File passwords contains all functions dealing with passwords.

package main

import (
	"bufio"
	"errors"
	"os"
	"strings"
	"transposer/filehelper"
	"transposer/logger"
	"transposer/mathhelper"
	"transposer/stringhelper"
)

// getPasswordList gets the password list from whatever source is specified.
func getPasswordList() (string, error) {
	listLen := len(passwordList)
	if listLen != 0 {
		if len(passwordFile) != 0 {
			return ``, ErrMultiplePasswordSources
		}

		if listLen > maxPasswordsLen {
			return ``, ErrPasswordsTooLong
		}

		return strings.TrimSpace(passwordList), nil
	} else {
		if len(passwordFile) != 0 {
			return readPasswordsFromFile()
		} else {
			return ``, ErrNoPasswordSources
		}
	}
}

// readPasswordsFromFile reads the passwords from a file.
func readPasswordsFromFile() (string, error) {
	size, err := filehelper.FileSize(passwordFile)
	if err != nil {
		return ``, err
	}

	if size < minPasswordLen {
		return ``, ErrPasswordsTooShort
	}

	if size > maxPasswordsLen+32 {
		return ``, ErrPasswordsTooLong
	}

	var file *os.File
	file, err = os.Open(passwordFile)
	if err != nil {
		return ``, err
	}
	defer filehelper.CloseWithName(file)

	// Create a new scanner to read the file line by line.
	scanner := bufio.NewScanner(file)

	sb := stringhelper.NewBuilder(int(size))
	next := false
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}

		if next {
			_, _ = sb.WriteRune(':')
		} else {
			next = true
		}

		_, _ = sb.WriteString(line)
	}

	// Check for errors during the scan
	if err = scanner.Err(); err != nil {
		return ``, err
	}

	if sb.Len() > maxPasswordsLen {
		return ``, ErrPasswordsTooLong
	}

	return sb.String(), nil
}

// processPasswordList analyzes the passwordList string and returns all given passwords as a slice.
// The passwords are converted to all lower-case. This is done to prevent weakening of the
// passwords by starting them with an upper-case letter and then continuing with lower
// case letters, which would effectively remove one letter from the password strength.
func processPasswordList(password string) ([]string, error) {
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
	// 1. Get lengths of passwords
	lengths, ok := passwordLengths(passwords)
	if !ok {
		return ok
	}

	logger.PrintInfof(mainMsgBase+10, `Password lengths: %v`, lengths)

	// 2. Check password lengths for common divisors.
	return checkPasswordLengths(passwords, lengths)
}

// passwordLengths creates a slice of the password lengths and checks them for valid lengths.
func passwordLengths(passwords []string) ([]int, bool) {
	ok := true
	lengths := make([]int, len(passwords))

	for i, pw := range passwords {
		pwl := len(pw)
		if pwl < minPasswordLen {
			logger.PrintErrorf(mainMsgBase+11, `Password '%s' is too short`, pw)
			ok = false
		}

		lengths[i] = pwl
	}

	return lengths, ok
}

// checkPasswordLengths checks the lengths
func checkPasswordLengths(passwords []string, lengths []int) bool {
	ok := true
	pwLen := len(passwords)

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
				ok = false
			}
		}
	}

	return ok
}
