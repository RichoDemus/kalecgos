//package wow_display_outdated_addons

package main

import (
	"fmt"
	"flag"
)

func main() {
	addonsDirectory := flag.String("addons-directory", "c:\\Games\\World of Warcraft\\Interface\\AddOns\\", "Path to the Addons folder")

	flag.Parse()

	//argsWithoutProg := os.Args[1:]
	fmt.Println("Hello World!")
	fmt.Println("Dir:", *addonsDirectory)
}
