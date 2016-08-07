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

	for id, version := range addons {
		fmt.Println("Checking addon " + id)
		url := createAddonUrl(id)
		fmt.Println("Addon url: " + url)
		fmt.Println("Installed version " + version)
		page := getWebpage2(url)
		//if strings.Trim(url, " ") == "https://mods.curse.com/addons/wow/deadly-boss-mods" {
		//	fmt.Println("EQUAL")
		//}
		//page := getWebpage2("https://mods.curse.com/addons/wow/deadly-boss-mods")
		// fmt.Println("Body: " + page)
		newestVersion := getAddonVersionFromCurseWebpage(page)
		fmt.Println("Latest version " + newestVersion)
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

func createAddonUrl(id string) string {
	rawUrl := "https://mods.curse.com/addons/wow/" + id
	// For some reason we get a stray ascii character 13 at the end
	trimmedUrl:=rawUrl[:len(rawUrl)-1]
	return trimmedUrl
}

func getWebpage(url string) string {
	//fmt.Println("Accessing " + url)
	//fmt.Println("url bytes: " + string(len(url)))
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

func getWebpage2(url string) string {
	//fmt.Println("GET " + url)
	//var a [60]byte

	//fmt.Printf("URL [%s]\n", url2)
	// fmt.Println(reflect.TypeOf(url))
	//fmt.Printf("url runes: %d\n", utf8.RuneCountInString(url))
	//fmt.Printf("url bytes: %d\n", len(url))
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Ubuntu Chromium/51.0.2704.79 Chrome/51.0.2704.79 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	return string(body)
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
