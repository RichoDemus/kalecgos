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
	"io"
)

func main() {
	addonsDirectoryPointer := flag.String("addons-directory", "c:\\Games\\World of Warcraft\\Interface\\AddOns\\", "Path to the Addons folder")

	flag.Parse()
	addonsDirectory := *addonsDirectoryPointer

	fmt.Println("Dir:", addonsDirectory)

	var addons = make(map[string]string)

	addonDirectories, _ := ioutil.ReadDir(addonsDirectory)
	for _, addonDirectory := range addonDirectories {
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

	for id, version := range addons {
		fmt.Println("Checking addon " + id)
		url := createAddonUrl(id)
		fmt.Println("Addon url: " + url)
		fmt.Println("Installed version " + version)
		page := getWebpage(url)
		newestVersion := getAddonVersionFromCurseWebpage(page)
		fmt.Println("Latest version " + newestVersion)
	}
}

func getAddonProperties(tocFile string) (string, string) {
	pattern, err := regexp.Compile(`X-Curse-Project-ID: (.*)`)
	if err != nil {
		log.Fatal(err)
	}
	rawId := pattern.FindStringSubmatch(tocFile)
	fixedId :=fixIdString(rawId[1])

	pattern, err = regexp.Compile(`X-Curse-Packaged-Version: (.*)`)
	if err != nil {
		log.Fatal(err)
	}

	version := pattern.FindStringSubmatch(tocFile)

	return fixedId, version[1]
}

func fixIdString(id string) string {
	lastCharacter  := id[len(id)-1:]
	if lastCharacter[0] == 13 {
		// For some reason we sometimes get a stray ascii character 13 at the end
		return id[:len(id)-1]
	}
	return id
}

func createAddonUrl(id string) string {
	return "https://mods.curse.com/addons/wow/" + id
}

func getWebpage(url string) string {
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
		return ""
	} else {
		defer response.Body.Close()
		if response.StatusCode != 200 {
			log.Println("Wrong status code: " + response.Status)
			log.Fatal("Body: " + responseToString(response.Body))
		}
		return string(responseToString(response.Body))
	}
}

func responseToString(body io.ReadCloser) string {
	bs, err := ioutil.ReadAll(body)
	if err != nil {
		log.Fatal(err)
	}
	return string(bs)
}

func getAddonVersionFromCurseWebpage(html string) string {
	pattern, err := regexp.Compile(`<li class="newest-file">Newest File: (.*)</li>`)
	if err != nil {
		log.Fatal(err)
	}
	version := pattern.FindStringSubmatch(html)
	return version[1]
}
