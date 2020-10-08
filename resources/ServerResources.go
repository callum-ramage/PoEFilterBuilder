package resources

import (
	"log"
	"os"
	"path/filepath"
)

// ServerRoot returns the home directory of the applicaiton
func ServerRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.FromSlash(dir + "/Content/")
}