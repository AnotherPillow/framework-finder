package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

func main() {
	config, _ := readConfigFile()

	dumpPath := config["dump-path"].(string)

	dumpPath = filepath.Clean(dumpPath) + string(filepath.Separator)

	var manifests []string
	frameworks := make(map[string]int)

	err := filepath.Walk(dumpPath, func(path string, info os.FileInfo, err error) error {
		// Check if the file is a regular file and has the name "manifest.json"
		if err == nil && !info.IsDir() && info.Name() == "manifest.json" {
			manifests = append(manifests, path)
			var manifest, _ = readJSONFile(path)

			_contentpackfor, ex := manifest["ContentPackFor"]
			if !ex {
				return nil
			}
			contentpackfor := _contentpackfor.(map[string]interface{})

			if contentpackfor["UniqueID"] == nil {
				return nil
			}

			_, ex2 := frameworks[contentpackfor["UniqueID"].(string)]
			if ex2 {
				frameworks[contentpackfor["UniqueID"].(string)]++
			} else {
				frameworks[contentpackfor["UniqueID"].(string)] = 1
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	/* https://stackoverflow.com/a/56706305 */
	keys := make([]string, 0, len(frameworks))
	for key := range frameworks {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool { return frameworks[keys[i]] > frameworks[keys[j]] })

	// fmt.Println("Manifest.json files in the directory:")
	// for id, amount := range frameworks {
	// 	fmt.Printf("%s: %d\n", id, amount)
	// }
	fmt.Printf("# Frameworks\n\n")
	fmt.Println("| ID | Count |")
	fmt.Println("| --- | ---- |")
	for _, key := range keys {
		fmt.Printf("| `%s` | %d |\n", key, frameworks[key])
		// fmt.Printf("%s: %d\n", key, frameworks[key])
	}
}
