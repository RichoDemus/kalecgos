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
	addonsDirectory := flag.String("addons-directory", "c:\\Games\\World of Warcraft\\Interface\\AddOns\\", "Path to the Addons folder")

	flag.Parse()
	strPointerValue := *addonsDirectory

	fmt.Println("Hello World!")
	fmt.Println("Dir:", strPointerValue)

	addonDirectories, _ := ioutil.ReadDir(strPointerValue)
	for _, addonDirectory := range addonDirectories {
		fmt.Println(strPointerValue +  addonDirectory.Name())
		filesInAddonDirectory, _ := ioutil.ReadDir(strPointerValue +  addonDirectory.Name() + "/")
		for _, filesInAddonDirectory := range filesInAddonDirectory {
			if strings.HasSuffix(filesInAddonDirectory.Name(), ".toc") {
				tocFile, err := ioutil.ReadFile(strPointerValue +  addonDirectory.Name() + "/" + filesInAddonDirectory.Name())
				if err != nil {
				    log.Fatal(err)
				}
				fmt.Println(string(tocFile))

			}

		}
	}
}
