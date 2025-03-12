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

var fileEncoding string
var password string
var conversion string
var useAllCharacters bool

var useConversion bool
var toLower bool
var passwords []string
var inputFiles []string

var encryptFlagSet = flag.NewFlagSet(`encrypt`, flag.ExitOnError)
var decryptFlagSet = flag.NewFlagSet(`decrypt`, flag.ExitOnError)

func defineCommandlineFlags() {
	encodingName := encodinghelper.PlatformDefaultEncodingName()

	encryptFlagSet.StringVar(&fileEncoding, `encoding`, encodingName, `character encoding for input[:output]`)
	encryptFlagSet.StringVar(&password, `passwords`, ``, `Transposition password(s) (separated by ':', if there is more than one)`)
	encryptFlagSet.StringVar(&conversion, `convert`, ``, `Characters are converted to 'lower' or 'upper' case`)
	encryptFlagSet.BoolVar(&useAllCharacters, `allchars`, false, `If set all characters are transposed (default: Only letters are transposed)`)

	encryptFlagSet.Usage = printUsageFunction

	decryptFlagSet.StringVar(&fileEncoding, `encoding`, encodingName, `character encoding for input[:output]`)
	decryptFlagSet.StringVar(&password, `passwords`, ``, `Transposition password(s) (separated by ':', if there is more than one)`)

	decryptFlagSet.Usage = printUsageFunction
}

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

func checkCommandlineFlags(doEncrypt bool) error {
	if len(password) == 0 {
		return errors.New(`transposition passwords must not be empty`)
	}

	var err error
	passwords, err = getPasswords(password)
	if err != nil {
		return err
	}

	if len(inputFiles) == 0 {
		return errors.New(`input file(s) must not be be empty`)
	}

	if doEncrypt {
		if len(conversion) != 0 {
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
				return errors.New(`no conversion type specified`)
			}
		}
	}

	return nil
}

func printUsageFunction() {
	w := flag.CommandLine.Output()

	_, _ = fmt.Fprintf(w,
		"\nUsage:\n   %s encrypt -encoding {input[:output]} -passwords {transpositionPasswords} -convert [upper|lower] -allchars {inputFilePath...}",
		myName)
	_, _ = fmt.Fprintf(w,
		"\n   %s decrypt -encoding {input[:output]} -passwords {transpositionPasswords} {inputFilePath...}\n",
		myName)

	_, _ = fmt.Fprintln(w, "\nEncrypt:")

	encryptFlagSet.PrintDefaults()

	_, _ = fmt.Fprintln(w, "\nDecrypt:")

	decryptFlagSet.PrintDefaults()

	_, _ = fmt.Fprintln(w)
	_, _ = fmt.Fprintf(w, "Maximum file size is %s bytes.\n", numberformat.FormatInt(maxFileSize))
	_, _ = fmt.Fprintln(w, "Output is written to the {inputFilepath} with '_trans' or '_untrans' added to the file name.")
	_, _ = fmt.Fprintln(w, "The encoding can be specified for the input and output file separately, separated by ':'.")
	_, _ = fmt.Fprintln(w, "If the output encoding is not specified it is the same as the input encoding.")
	_, _ = fmt.Fprintln(w, "\nEncoding may be any of the following:")

	for _, encodingName := range encodinghelper.EncodingNames() {
		_, _ = fmt.Fprint(w, `   `)
		_, _ = fmt.Fprintln(w, encodingName)
	}

	_, _ = fmt.Fprintln(w, "\nIf the output encoding begins with 'utf' a BOM (byte order mark) has to be handled.")
	_, _ = fmt.Fprintln(w, `If the encoding has the word 'bom' appended, a BOM is written at the beginning of the output file.`)
	_, _ = fmt.Fprintln(w, `If the encoding has the word 'nobom' appended, no BOM is written at the beginning of the output file.`)
	_, _ = fmt.Fprintln(w, `If the encoding is followed by neither 'bom', nor 'nobom' the output file only has a BOM, if the input file has one.`)
	_, _ = fmt.Fprintln(w, `I.e. 'bom' forces a BOM of the output file, 'nobom' forces no BOM of the output file and neither one copies the BOM setting from the input file.`)
}
