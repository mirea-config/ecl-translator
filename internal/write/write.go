package write

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ToEclFile(text string, outputPath string) error {
	switch ext := filepath.Ext(outputPath); ext {
	case "":
		if strings.HasSuffix(outputPath, ".") {
			outputPath += "ecl"
		} else {
			outputPath += ".ecl"
		}

	case ".ecl":
		break

	default:
		return fmt.Errorf("'%s' is not a valid filename extension", ext)
	}

	fmt.Println(outputPath)

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	if _, err = outputFile.Write([]byte(text)); err != nil {
		return err
	}

	return nil
}
