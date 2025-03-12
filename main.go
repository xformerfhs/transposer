package main

import (
	"errors"
	"os"
	"runtime"
	"strings"
	"transposer/encodedfile"
	"transposer/encodinghelper"
	"transposer/filehelper"
	"transposer/logger"
	"transposer/slicehelper"
	"transposer/stringhelper"
	"transposer/transposition"
)

// ******** Private constants ********

const maxFileSize = 10_000_000

const fmtEncodedFileOperation = `%s file '%s' with %s encoding`

const myVersion = `1.0.0`
const mainMsgBase = 10

var myName = filehelper.RealBaseName(os.Args[0])

// ******** Main functions ********

func main() {
	os.Exit(realMain())
}

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
			return printProcessingErrorf(mainMsgBase+4, `Error reading file '%s': %s`, inputFilePath, err.Error())
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
		if doEncrypt {
			if useAllCharacters {
				// Read all characters.
				inputContent, useWindowsLineBreak, err = encodedfile.ReadAllRunes(inputFilePath, maxFileSize, actInputEncodingName, useConversion, toLower)
			} else {
				// Read only letters.
				inputContent, useWindowsLineBreak, err = encodedfile.ReadOnlyLetterRunes(inputFilePath, maxFileSize, actInputEncodingName, useConversion, toLower)
			}
		} else {
			// When decrypting no character transformation or filtering is wanted.
			inputContent, useWindowsLineBreak, err = encodedfile.ReadAllRunesAsIs(inputFilePath, maxFileSize, actInputEncodingName)
		}

		if err != nil {
			return printProcessingErrorf(mainMsgBase+6, `Error reading file '%s': %s`, inputFilePath, err.Error())
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

		slicehelper.ClearInteger(inputContent)

		if err != nil {
			return printProcessingErrorf(mainMsgBase+8, `Error writing file '%s': %s`, inputFilePath, err.Error())
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
		logger.PrintInfof(mainMsgBase+10, `%s V%s (%s)`, myName, myVersion, runtime.Version())
		return false, rcOK, true
	}

	// Tell, if it is encrypt, or not.
	doEncrypt := cmd == 'e'

	return doEncrypt, rcOK, false
}

// getPasswords analyzes the password string and returns all given passwords as a slice.
// The passwords are converted to all lower-case. This is done to prevent weaking of the
// passwords by starting them with an upper-case letter and then continuing with lower
// case letters, which would effectively remove one letter from the password strength.
func getPasswords(password string) ([]string, error) {
	elements := strings.Split(password, `:`)
	result := make([]string, len(elements))
	for i, candidate := range elements {
		if !stringhelper.IsAlphaNumeric(candidate) {
			return nil, errors.New(`passwords must only contain alphanumeric characters`)
		}
		result[i] = strings.ToLower(candidate)
	}

	return result, nil
}
