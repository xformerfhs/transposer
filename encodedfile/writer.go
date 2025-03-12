package encodedfile

import (
	"bufio"
	"golang.org/x/text/transform"
	"os"
	"transposer/encodinghelper"
	"transposer/filehelper"
)

func WriteEncoded(
	path string,
	encodingName string,
	useBom bool,
	useWindowsLineBreaks bool,
	data []rune) error {
	bw, f, err := prepareFileFoWriting(path, encodingName, useBom)
	if err != nil {
		return err
	}
	defer filehelper.CloseWithName(f)

	for _, r := range data {
		if useWindowsLineBreaks && r == lineFeed {
			err = bw.WriteByte(carriageReturn)
			if err != nil {
				return err
			}
		}

		_, err = bw.WriteRune(r)
		if err != nil {
			return err
		}
	}

	return bw.Flush()
}

// ******** Private functions ********

// prepareFileFoWriting prepares the given file path for writing.
func prepareFileFoWriting(path string, encodingName string, useBom bool) (*bufio.Writer, *os.File, error) {
	var f *os.File
	f, err := os.Create(path)
	if err != nil {
		return nil, nil, err
	}

	encoded, _, _ := encodinghelper.TranslateEncoding(encodingName, useBom)

	return bufio.NewWriter(transform.NewWriter(f, encoded.NewEncoder())), f, nil
}
