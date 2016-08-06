package main

import "testing"

func TestParseTocFile(t *testing.T) {
	id, version := getAddonProperties(`## X-Website: http://www.deadlybossmods.com
## X-Curse-Packaged-Version: 7.0.1
## X-Curse-Project-Name: Deadly Boss Mods
## X-Curse-Project-ID: deadly-boss-mods`)

	if id != "deadly-boss-mods" {
		t.Log("Wrong id: " + id)
		t.Fail()
	}
	if version != "7.0.1" {
		t.Log("Wrong version: " + version)
		t.Fail()
	}
}

func TestParseVersionFromPage(t *testing.T) {
	version := getAddonVersionFromCurseWebpage(`<li class="newest-file">Newest File: 7.0.3.7</li>`)

	if version != "7.0.3.7" {
		t.Log("Wrong version: " + version)
		t.Fail()
	}
}