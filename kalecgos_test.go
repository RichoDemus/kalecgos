package main

import (
	"testing"
)

func TestParseTocFile(t *testing.T) {
	id, version := getAddonProperties("dbm", `## X-Website: http://www.deadlybossmods.com
## X-Curse-Packaged-Version: 7.0.1
## X-Curse-Project-Name: Deadly Boss Mods
## X-Curse-Project-ID: deadly-boss-mods
## X-Curse-Repository-ID: wow/deadly-boss-mods/mainline`)

	if id != "deadly-boss-mods" {
		t.Log("Wrong id: " + id)
		t.Fail()
	}
	if version != "7.0.1" {
		t.Log("Wrong version: " + version)
		t.Fail()
	}
}

func TestCreateAddonURL(t *testing.T) {
	result := createAddonUrl("deadly-boss-mods")

	if result != "https://mods.curse.com/addons/wow/deadly-boss-mods" {
		t.Log("Wrong url: " + result)
		t.Fail()
	}
}

func TestParseVersionFromPage(t *testing.T) {
	version := getAddonVersionFromCurseWebpage("DBM", `<li class="newest-file">Newest File: 7.0.3.7</li>`)

	if version != "7.0.3.7" {
		t.Log("Wrong version: " + version)
		t.Fail()
	}
}

func TestParseVersionFromTocFileWithXCurseVersion(t *testing.T) {
	result := parseVersion("dbm", `## X-Website: http://www.deadlybossmods.com
	## X-Curse-Packaged-Version: 7.0.1
	## X-Curse-Project-Name: Deadly Boss Mods
	## X-Curse-Project-ID: deadly-boss-mods
	## X-Curse-Repository-ID: wow/deadly-boss-mods/mainline`)

	if result != "7.0.1" {
		t.Log("Wrong version: " + result)
		t.Fail()
	}
}

func TestParseVersionFromTocFileNormalVersion(t *testing.T) {
	result := parseVersion("dbm", `## X-Website: http://www.deadlybossmods.com
	## Version: 2.8.3
	## X-Curse-Project-Name: Deadly Boss Mods
	## X-Curse-Project-ID: deadly-boss-mods
	## X-Curse-Repository-ID: wow/deadly-boss-mods/mainline`)

	if result != "2.8.3" {
		t.Log("Wrong version: " + result)
		t.Fail()
	}
}

func TestParseIDFromPage(t *testing.T) {
	result := getAddonIdFromCurseWebpage("Dominos", `<tr class="wow">
	    <td>
	        <dl>
	            <dt><a href="/addons/wow/dominos">Dominos</a></dt>

	                <dd>Project Manager: <a href="/members/Tuller">Tuller</a></dd>

	        </dl>
	    </td>
	`)

	if result != "dominos" {
		t.Log("Wrong id: " + result)
		t.Fail()
	}
}

func TestCreateSearchUrl(t *testing.T) {
	result := createSeatchUrl("Deadly Boss Mods")

	if result != "https://mods.curse.com/search?search=Deadly+Boss+Mods" {
		t.Log("Wrong url: " + result)
		t.Fail()
	}
}
