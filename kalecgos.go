package main

import (
	"fmt"
	"flag"
	"io/ioutil"
	"strings"
	"log"
	"regexp"
	"net/http"
	"net/url"
	"os"
	"io"
	"github.com/skratchdot/open-golang/open"
)

type addon struct {
	id            string
	version       string
	newVersion    string
	hasNewVersion bool
	url           string
	successful    bool
}

func main() {
	addonsDirectoryPointer := flag.String("addons-directory", "Interface/AddOns/", "Path to the Addons folder")

	flag.Parse()
	addonsDirectory := *addonsDirectoryPointer

	fmt.Println("Dir:", addonsDirectory)

	addons := getAddons(addonsDirectory)

	addons = addVersionDataToAddons(addons)

	f, err := os.Create("addons.html")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	f.WriteString("<html><body><h1>Addons:</h1><h2>Outdated:</h2><ul>\n")

	for _, addon := range addons {
		if addon.hasNewVersion && addon.successful {
			fmt.Println("Found newer version of", addon.id, "(", addon.version, "->", addon.newVersion, "): ", addon.url)
			f.WriteString("<li>Newer version of " + addon.id + " ( " + addon.version + " -> " + addon.newVersion + " ): <a href=\"" + addon.url + "\">Curse link</a></li>\n")
		}
	}

	f.WriteString("</ul><h2>Up do date:</h2><ul>\n")

	for _, addon := range addons {
		if !addon.hasNewVersion && addon.successful {
			fmt.Println("Addon", addon.id, "(", addon.version, ") is at the latest version")
			f.WriteString("<li>Addon " + addon.id + " ( " + addon.version + " ) is at the latest version</li>\n")
		}
	}

	f.WriteString("</ul><h2>Unable to determine if a new version is available:</h2><ul>\n")

	for _, addon := range addons {
		if !addon.successful {
			fmt.Println("Failed to scan", addon.id)
			f.WriteString("<li>Failed to scan addon " + addon.id + "</li>\n")
		}
	}

	f.WriteString("</ul></body></html>")
	f.Sync()
	open.Run("addons.html")
}

// Takes the path to the addons directory and returns a slice of addons
func getAddons(addonsDirectory string) []addon {
	addons := make([]addon, 0)
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
				id, version := getAddonProperties(addonDirectory.Name(), string(tocFile))
				if version == "" {
					addon := addon{id: id, successful: false}
					if id != "" && !contains(addons, addon) {
						addons = append(addons, addon)
					}
				} else {
					addon := addon{id: id, version: version, successful: true}
					if id != "" && !contains(addons, addon) {
						addons = append(addons, addon)
					}
				}

			}
		}
	}
	return addons
}

func addVersionDataToAddons(addons []addon) []addon {
	result := make([]addon, 0)
	for _, oldAddon := range addons {
		if !oldAddon.successful {
			result = append(result, oldAddon)
			continue
		}
		url := createAddonUrl(oldAddon.id)
		page := getWebpage(url)
		newestVersion := getAddonVersionFromCurseWebpage(oldAddon.id, page)
		id := oldAddon.id
		version := oldAddon.version
		if oldAddon.version != newestVersion {
			newAddon := addon{id: id, version: version, newVersion: newestVersion, hasNewVersion: true, successful:true, url: url}
			result = append(result, newAddon)
		} else {
			newAddon := addon{id: id, version: version, hasNewVersion: false, successful:true, url:url}
			result = append(result, newAddon)
		}
	}
	return result
}

func getAddonProperties(addon string, tocFile string) (string, string) {
	xId, title := parseAddonId(tocFile)
	if(len(xId) == 0) {
		xId = tryToFindAddonOnCurseSite(title)
	}
	version := parseVersion(addon, tocFile)

	return xId, version
}

func parseAddonId(tocFile string) (string, string) {
	pattern, err := regexp.Compile(`X-Curse-Project-ID: (.*)`)
	if err != nil {
		log.Fatal(err)
	}
	rawId := pattern.FindStringSubmatch(tocFile)
	if len(rawId) == 0 {
		pattern, err = regexp.Compile(`Title: (.*)`)
		if err != nil {
			log.Fatal(err)
		}
		rawId = pattern.FindStringSubmatch(tocFile)
		return "", fixParsedString(rawId[1])
	}
	return fixParsedString(rawId[1]), ""
}

func parseVersion(addon string, tocFile string) string {
	pattern, err := regexp.Compile(`Version: (.*)`)
	if err != nil {
		log.Fatal(err)
	}

	rawVersion := pattern.FindStringSubmatch(tocFile)
	if(len(rawVersion) == 0) {
		log.Println("Failed to parse version for addon: " + addon)
		return ""
	}

	fixedVersion := fixParsedString(rawVersion[1])

	return fixedVersion
}

func fixParsedString(str string) string {
	lastCharacter := str[len(str) - 1:]
	if lastCharacter[0] == 13 {
		// For some reason we sometimes get a stray ascii character 13 at the end
		return str[:len(str) - 1]
	}
	return str
}

func tryToFindAddonOnCurseSite(title string) string {
	url := createSeatchUrl(title)
	page := getWebpage(url)
	if(len(page) == 0) {
		return ""
	}
	id := getAddonIdFromCurseWebpage(title, page)
	return id
}

func getAddonIdFromCurseWebpage(title string, html string) string {
	pattern, err := regexp.Compile(`<dt><a href="/addons/wow/(.*)">(.*)</a></dt>`)
	if err != nil {
		log.Fatal(err)
	}
	version := pattern.FindStringSubmatch(html)
	if len(version) == 0 {
		log.Println("Didn't find addon id when searhing curse page for addon ",title)
		return title
	}
	return version[1]
}

func createAddonUrl(id string) string {
	return "https://mods.curse.com/addons/wow/" + id
}

func createSeatchUrl(title string) string {
	searchUrl,_ := url.Parse("https://mods.curse.com/search?search=" + url.QueryEscape(title))
	return searchUrl.String()
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
			log.Println("Failed to fetch page: " + url)
			log.Println("Wrong status code: " + response.Status)
			log.Println("Body: " + responseToString(response.Body))
			return ""
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

func getAddonVersionFromCurseWebpage(addon string, html string) string {
	pattern, err := regexp.Compile(`<li class="newest-file">Newest File: (.*)</li>`)
	if err != nil {
		log.Fatal(err)
	}
	version := pattern.FindStringSubmatch(html)
	if len(version) == 0 {
		log.Println("Didn't find addon version for addon ",addon)
		return ""
	}
	return version[1]
}

func contains(s []addon, e addon) bool {
	for _, a := range s {
		if a.id == e.id {
			return true
		}
	}
	return false
}
