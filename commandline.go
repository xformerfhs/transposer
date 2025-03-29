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

// File commandline contains the functions that parse the command line in the main program.

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"transposer/encodinghelper"
	"transposer/numberformat"
)

// ******** Private variables ********

// -------- Command line variables --------

// fileEncoding contains the file encoding string from the command line.
var fileEncoding string

// passwordList contains the passwords string from the command line.
var passwordList string

var passwordFile string

// conversion contains the convert string from the command line.
var conversion string

// onlyLetters contains the onlyletters value from the command line.
var onlyLetters bool

// -------- Variables derived from the command line --------

// useConversion indicates that characters should be converted.
var useConversion bool

// toLower indicates that characters are converted to lower case.
// If useConversion is true and toLower is false characters are converted to upper case.
var toLower bool

// passwords contains the list of passwords from the command line.
var passwords []string

// inputFiles contains the list of file paths from the command line.
var inputFiles []string

// -------- Flag sets --------

// encryptFlagSet is the flag set used for encryption.
var encryptFlagSet = flag.NewFlagSet(`encrypt`, flag.ExitOnError)

// decryptFlagSet is the flag set used for decryption.
var decryptFlagSet = flag.NewFlagSet(`decrypt`, flag.ExitOnError)

// -------- Errors --------

// ErrMultiplePasswordSources is returned, when multiple password sources where specified.
var ErrMultiplePasswordSources = errors.New(`multiple password sources specified`)

// ErrNoPasswordSources is returned, when no password sources where specified.
var ErrNoPasswordSources = errors.New(`no passwords supplied`)

// ErrPasswordsTooShort is returned, when byte count of the passwords is too short.
var ErrPasswordsTooShort = errors.New(`passwords are too short or file name is not a normal file`)

// ErrPasswordsTooLong is returned, when byte count of the passwords is too long.
var ErrPasswordsTooLong = errors.New(`passwords are too long`)

// ErrNoInputFiles is returned, when no input files were specified.
var ErrNoInputFiles = errors.New(`no input file present`)

// ErrNoConversionType is returned, when the "case" option was given without a value.
var ErrNoConversionType = errors.New(`no conversion type specified`)

// ******** Private functions ********

// defineCommandlineFlags defines the flag sets for command line processing.
func defineCommandlineFlags() {
	encodingName := encodinghelper.PlatformDefaultEncodingName()

	// 1. Encryption.
	encryptFlagSet.StringVar(&fileEncoding, `encoding`, encodingName, `character encoding for input[:output] (separated by ':' or ','`)
	encryptFlagSet.StringVar(&passwordList, `passwords`, ``, `Transposition password(s) (separated by ':' or ',', if there is more than one)`)
	encryptFlagSet.StringVar(&passwordFile, `passwords-file`, ``, `File which contains the transposition password(s) (separated by ':' or ',', if there is more than one)`)
	encryptFlagSet.StringVar(&conversion, `case`, ``, `Characters are converted to 'lower' or 'upper' case`)
	encryptFlagSet.BoolVar(&onlyLetters, `onlyletters`, false, `If set only letters are read and transposed (default: All characters are read)`)

	encryptFlagSet.Usage = printUsageFunction

	// 2. Decryption.
	decryptFlagSet.StringVar(&fileEncoding, `encoding`, encodingName, `character encoding for input[:output]`)
	decryptFlagSet.StringVar(&passwordList, `passwords`, ``, `Transposition passwordList(s) (separated by ':', if there is more than one)`)
	decryptFlagSet.StringVar(&passwordFile, `passwords-file`, ``, `File which contains the transposition password(s) (separated by ':' or ',', if there is more than one)`)

	decryptFlagSet.Usage = printUsageFunction
}

// processCommandlineFlags processes the command line with the correct flag set.
func processCommandlineFlags(doEncrypt bool) error {
	var fs *flag.FlagSet
	if doEncrypt {
		fs = encryptFlagSet
	} else {
		fs = decryptFlagSet
	}

	err := fs.Parse(os.Args[2:])
	if err != nil {
		return err
	}

	inputFiles = fs.Args()

	return nil
}

// checkCommandlineFlags does some basic checks of the command line.
func checkCommandlineFlags(doEncrypt bool) error {
	var err error

	// 1. Check that there are input files.
	if len(inputFiles) == 0 {
		return ErrNoInputFiles
	}

	// 2. Check that passwords is set and not too long.
	passwordList, err = getPasswordList()
	if err != nil {
		return err
	}

	// 3. Convert the passwords string to a list of passwords.
	passwords, err = processPasswordList(passwordList)
	if err != nil {
		return err
	}

	if doEncrypt {
		// 4. Check conversion string for encryption.
		if len(conversion) != 0 {
			err = processConversionFlag()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// processConversionFlag analyzes the conversion flag and sets the corresponding variables.
func processConversionFlag() error {
	conversion = strings.ToLower(strings.TrimSpace(conversion))

	if len(conversion) != 0 {
		switch conversion[0] {
		case 'l':
			useConversion = true
			toLower = true
		case 'u':
			useConversion = true
			toLower = false
		default:
			return fmt.Errorf(`unknown conversion: '%s'`, conversion)
		}
	} else {
		return ErrNoConversionType
	}

	return nil
}

// printUsageFunction is the usage function of the flag sets.
func printUsageFunction() {
	w := flag.CommandLine.Output()

	_, _ = fmt.Fprintf(w,
		"\nUsage:\n   %s encrypt -encoding {input[:output]} {-passwords {transpositionPasswords} | -passwords-file {filename}] -case [upper|lower] -onlyletters {inputFilePath...}",
		myName)
	_, _ = fmt.Fprintf(w,
		"\n   %s decrypt -encoding {input[:output]} {-passwords {transpositionPasswords} | -passwords-file {filename}] {inputFilePath...}\n",
		myName)
	_, _ = fmt.Fprintf(w,
		"\n   %s help\n",
		myName)
	_, _ = fmt.Fprintf(w,
		"\n   %s version\n",
		myName)

	_, _ = fmt.Fprintln(w, "\nEncrypt:")

	encryptFlagSet.PrintDefaults()

	_, _ = fmt.Fprintln(w, "\nDecrypt:")

	decryptFlagSet.PrintDefaults()

	_, _ = fmt.Fprintln(w, "\nHelp:")
	_, _ = fmt.Fprintln(w, "  Show this help")

	_, _ = fmt.Fprintln(w, "\nVersion:")
	_, _ = fmt.Fprintln(w, "  Show version information")

	_, _ = fmt.Fprintln(w)
	_, _ = fmt.Fprintf(w, "Maximum file size is %s bytes.\n", numberformat.FormatInt(maxFileSize))
	_, _ = fmt.Fprintln(w, `Output is written to the {inputFilepath} with '_trans' or '_untrans' added to the file name.`)
	_, _ = fmt.Fprintln(w, `The passwords can either be specified by '-passwords' or 'passwords-file', but not both.`)
	_, _ = fmt.Fprintln(w, `If '-passwords-file' is specified, the passwords are read from specified file.`)
	_, _ = fmt.Fprintln(w, `The passwords are one string separated by ':' or ','.`)
	_, _ = fmt.Fprintln(w, `In a password file they can also be separated by new line.`)
	_, _ = fmt.Fprintf(w, "There can be up to %s passwords with a maximum total length of %s characters, including separators\n",
		numberformat.FormatInt(maxNumPasswords),
		numberformat.FormatInt(maxPasswordsLen))
	_, _ = fmt.Fprintln(w, `The encoding can be specified for the input and output file separately, separated by ':' or ','.`)
	_, _ = fmt.Fprintln(w, `If the output encoding is not specified it is the same as the input encoding.`)
	_, _ = fmt.Fprintln(w, "\nEncoding may be any of the following:")

	for _, encodingName := range encodinghelper.EncodingNames() {
		_, _ = fmt.Fprint(w, `  `)
		_, _ = fmt.Fprintln(w, encodingName)
	}

	_, _ = fmt.Fprintln(w, "\nIf the output encoding begins with 'utf' a BOM (byte order mark) has to be handled.")
	_, _ = fmt.Fprintln(w, `If the encoding has the word 'bom' appended, a BOM is written at the beginning of the output file.`)
	_, _ = fmt.Fprintln(w, `If the encoding has the word 'nobom' appended, no BOM is written at the beginning of the output file.`)
	_, _ = fmt.Fprintln(w, `If the encoding is followed by neither 'bom', nor 'nobom' the output file only has a BOM, if the input file has one.`)
	_, _ = fmt.Fprintln(w, `I.e. 'bom' forces a BOM of the output file, 'nobom' forces no BOM of the output file and neither one copies the BOM setting from the input file.`)
}
