// Package apache provides functions for reading the Apache configuration files
package apache

import (
	"os"
	"path"
	"io/ioutil"
  	"regexp"
	"log"
	"fmt"
)

// OpenApacheConfigurationFile opens a single apache configuration file and returns a string with the content
func OpenApacheConfigurationFile(confPath string) (res string) {
	currentBuf, err := ioutil.ReadFile(confPath)
	if err != nil {
		log.Fatal(err)
	}
	res = string(currentBuf)
	return
}


// OpenAllApacheConfigurationFiles opens an apache configuration file (and all the included ones) and returns their content as a map of strings:
//    - Key: Path to the Apache configuration file
//    - Value: String with the file content
func OpenAllApacheConfigurationFiles(confPath, apacheRoot string) (resBuffers map[string]string) {
	remainingConfFiles := []string{confPath}
	resBuffers = make(map[string]string);
	for len(remainingConfFiles) > 0 {
		currentConfFile := remainingConfFiles[0]
		remainingConfFiles = remainingConfFiles[1:]
		if _, ok := resBuffers[currentConfFile]; ok {
			continue
		}
		if _, err := os.Stat(currentConfFile); err == nil {
			bufferString := OpenApacheConfigurationFile(currentConfFile)
			resBuffers[currentConfFile] = bufferString
			includeRe, err := regexp.CompilePOSIX("^[[:space:]]*Include [\"]?([^\n\"]+)[\"]?")
			if err != nil {
				log.Fatal(err)
			}
			includeMatches := includeRe.FindAllStringSubmatch(bufferString, -1);
			for _, element := range includeMatches {
				newFile := element[1]
				if !path.IsAbs(newFile) {
					newFile = path.Join(apacheRoot, newFile)
				}
				remainingConfFiles = append(remainingConfFiles, newFile)
			}
		} else {
			fmt.Println("Skipping " + currentConfFile + " as it does not exist in the filesystem")
		}
	}
	return resBuffers
}
