package replacement

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// Replacements maps reywords to result maps
var Replacements map[string]string

// ReplacementFolder is the folder when replacements are saved/loaded from
var ReplacementFolder string

// LoadReplacements loads replacement files from the specified folder location
func LoadReplacements(folderLocation string) {
	ReplacementFolder = folderLocation
	Replacements = make(map[string]string)
	err := filepath.Walk(folderLocation, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".rep" {
			dat, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			fileName := filepath.Base(path)
			fileName = fileName[:len(fileName)-4]
			Replacements[fileName] = string(dat)
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}

// GetReplacementKeys retrieves the keys from the Replacements map and returns them as an array
func GetReplacementKeys() []string {
	keys := make([]string, 0, len(Replacements))
	for k := range Replacements {
		keys = append(keys, k)
	}
	return keys
}

// SaveReplacement saves the content as a replacement file with name {name}.rep
func SaveReplacement(name string, content string) {
	if len(name) > 0 {
		d1 := []byte(content)
		err := ioutil.WriteFile(ReplacementFolder+name+".rep", d1, 0644)
		if err != nil {
			fmt.Println(err)
		}
		Replacements[name] = content
	}
}

// DeleteReplacement deletes the replacement file with this name
func DeleteReplacement(name string) {
	err := os.Remove(ReplacementFolder + name + ".rep")
	if err != nil {
		fmt.Println(err)
	}
	delete(Replacements, name)
}

// findReplacementKeys finds the keys of replacement expressions in a filter
func findReplacementKeys(filterContent string) [][]byte {
	re := regexp.MustCompile(`\{[^\{\}]+\}`)
	return re.FindAll([]byte(filterContent), -1)
}

func splitParameters(filterKey string) (string, []string) {
	filterKey = strings.Trim(filterKey, "{}")
	components := strings.Split(filterKey, ",")
	return components[0], components[1:]
}

func replaceParameters(filterContent string, replacementParams []string) string {
	for i, val := range replacementParams {
		filterContent = strings.Replace(filterContent, "["+strconv.Itoa(i+1)+"]", val, -1)
	}
	return filterContent
}

func replaceContent(filterContent string, replacementKey string, replacementName string, replacementParams []string) string {
	if val, ok := Replacements[replacementName]; ok {
		resultReplacement := ReplaceKeys(replaceParameters(val, replacementParams))
		filterContent = strings.Replace(filterContent, replacementKey, resultReplacement, -1)
	}
	return filterContent
}

// ReplaceKeys replaces the replacement keys found within a filter with their value
func ReplaceKeys(filterContent string) string {
	keys := findReplacementKeys(filterContent)
	for _, key := range keys {
		keyString := string(key)
		keyName, parameters := splitParameters(keyString)
		filterContent = replaceContent(filterContent, keyString, keyName, parameters)
	}
	return filterContent
}
