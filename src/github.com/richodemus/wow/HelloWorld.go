//package wow_display_outdated_addons

package main

import (
	"fmt"
	"flag"
	"io/ioutil"
	"strings"
	"log"
	"regexp"
	"net/http"
	"os"
)

func main() {
	addonsDirectoryPointer := flag.String("addons-directory", "c:\\Games\\World of Warcraft\\Interface\\AddOns\\", "Path to the Addons folder")

	flag.Parse()
	addonsDirectory := *addonsDirectoryPointer

	fmt.Println("Dir:", addonsDirectory)

	var addons = make(map[string]string)

	addonDirectories, _ := ioutil.ReadDir(addonsDirectory)
	for _, addonDirectory := range addonDirectories {
		// fmt.Println("Addon: " + addonsDirectory +  addonDirectory.Name())
		filesInAddonDirectory, err := ioutil.ReadDir(addonsDirectory + "/" + addonDirectory.Name() + "/")
		if err != nil {
			log.Fatal(err)
		}
		for _, addonFile := range filesInAddonDirectory {
			if strings.HasSuffix(addonFile.Name(), ".toc") {
				tocFile, err := ioutil.ReadFile(addonsDirectory + "/" + addonDirectory.Name() + "/" + addonFile.Name())
				if err != nil {
					log.Fatal(err)
				}
				id, version := getAddonProperties(string(tocFile))
				addons[id] = version
			}
		}
	}

	for key, value := range addons {
		fmt.Println("Key:", key)
		fmt.Println("Value:", value)
	}
}

func getAddonProperties(tocFile string) (string, string) {
	pattern, err := regexp.Compile(`X-Curse-Project-ID: (.*)`)
	if err != nil {
		log.Fatal(err)
	}
	id := pattern.FindStringSubmatch(tocFile)

	pattern, err = regexp.Compile(`X-Curse-Packaged-Version: (.*)`)
	if err != nil {
		log.Fatal(err)
	}

	version := pattern.FindStringSubmatch(tocFile)

	return id[1], version[1]
}

func getWebpage(url string) string {
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
		return ""
	} else {
		defer response.Body.Close()
		bs, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		return string(bs)
	}
}

func getAddonVersionFromCurseWebpage(html string) string {
	pattern, err := regexp.Compile(`<li class="newest-file">Newest File: (.*)</li>`)
	if err != nil {
		log.Fatal(err)
	}
	version := pattern.FindStringSubmatch(html)
	return version[1]
}
