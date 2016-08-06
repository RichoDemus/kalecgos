//package wow_display_outdated_addons

package main

import (
	"fmt"
	"flag"
	"io/ioutil"
	"strings"
	"log"
)

func main() {
	addonsDirectoryPointer := flag.String("addons-directory", "c:\\Games\\World of Warcraft\\Interface\\AddOns\\", "Path to the Addons folder")

	flag.Parse()
	addonsDirectory := *addonsDirectoryPointer

	fmt.Println("Hello World!")
	fmt.Println("Dir:", addonsDirectory)

	addonDirectories, _ := ioutil.ReadDir(addonsDirectory)
	for _, addonDirectory := range addonDirectories {
		fmt.Println("Addon: " + addonsDirectory +  addonDirectory.Name())
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
				fmt.Println("Tocfile: " + string(tocFile))

			}

		}
	}
}
