package filters

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"../replacement"
)

// Filters maps reywords to result maps
var Filters map[string]string

// FilterFolder is the folder that filters are saved/loaded from
var FilterFolder string

// OutputLocation is the folder that compiled filters get saved to
var OutputLocation string

// LoadFilters loads filter files from the specified folder location
func LoadFilters(folderLocation string, outputLocation string) {
	FilterFolder = folderLocation
	OutputLocation = outputLocation
	Filters = make(map[string]string)
	err := filepath.Walk(folderLocation, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".fil" {
			dat, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			fileName := filepath.Base(path)
			fileName = fileName[:len(fileName)-4]
			Filters[fileName] = string(dat)
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}

// GetFilterKeys retrieves the keys from the Filters map and returns them as an array
func GetFilterKeys() []string {
	keys := make([]string, 0, len(Filters))
	for k := range Filters {
		keys = append(keys, k)
	}
	return keys
}

// SaveFilter saves the content as a filter file with name {name}.rep
func SaveFilter(name string, content string) {
	if len(name) > 0 {
		d1 := []byte(content)
		err := ioutil.WriteFile(FilterFolder+name+".fil", d1, 0644)
		if err != nil {
			fmt.Println(err)
		}
		Filters[name] = content
	}
}

// DeleteFilter deletes the filter file with this name
func DeleteFilter(name string) {
	err := os.Remove(FilterFolder + name + ".fil")
	if err != nil {
		fmt.Println(err)
	}
	delete(Filters, name)
}

// CompileFilter processes substitutions and builds the result filter
func CompileFilter(name string) {
	if len(name) > 0 {
		if content, ok := Filters[name]; ok {
			content = replacement.ReplaceKeys(content)
			d1 := []byte(content)
			err := ioutil.WriteFile(OutputLocation+name+".filter", d1, 0644)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
