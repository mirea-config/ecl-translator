package translate

import (
	parser "ecl-translator/internal/parse"
	"ecl-translator/internal/write"
	"strings"
)

func Translate(src []byte, outputPath string) error {
	lines, err := parser.ParseJsonInput(src)
	if err != nil {
		return err
	}

	text := strings.Join(lines, "\n\n")

	if err = write.ToEclFile(text, outputPath); err != nil {
		return err
	}

	return nil
}
