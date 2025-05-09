package utils

import (
	"fmt"
	"os"
)

func CreateDirs(paths ...string) error {
	for _, path := range paths {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return fmt.Errorf("unable to create directory %s: %w", path, err)
		}
	}
	return nil
}
