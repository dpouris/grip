package command

import "os"

// Validates that a given directory, dir, is valid
func validateDir(dir string) bool {
	_, err := os.ReadDir(dir)

	return err == nil
}
