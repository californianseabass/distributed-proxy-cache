package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestReloadRing(t *testing.T) {
	var stringToHash map[string]string = map[string]string{
		"Google.com":          "",
		"Youtube.com":         "",
		"Facebook.com":        "",
		"Wikipedia.org":       "",
		"Yahoo.com":           "",
		"Amazon.com":          "",
		"Pages.tmall.com":     "",
		"Reddit.com":          "",
		"Live.com":            "",
		"Netflix.com":         "",
		"Blogspot.com":        "",
		"Office.com":          "",
		"Instagram.com":       "",
		"Yahoo.co.jp":         "",
		"Bing.com":            "",
		"Microsoft.com":       "",
		"Google.com.hk":       "",
		"Stackoverflow.com":   "",
		"Babytree.com":        "",
		"Twitter.com":         "",
		"Ebay.com":            "",
		"Amazon.co.jp":        "",
		"Twitch.tv":           "",
		"Apple.com":           "",
		"Google.co.in":        "",
		"Microsoftonline.com": "",
		"Msn.com":             "",
		"Wordpress.com":       "",
	}

	hashToServerFile, err := ioutil.TempFile(os.TempDir(), "ketama.ring")

	if err != nil {
		t.Error(err)
		return
	}

	servers := getServers()
	ketama := createRingFromServers(servers)
	segmentId := persist(hashToServerFile.Name(), ketama)
	for siteName := range stringToHash {
		stringToHash[siteName] = ketama.getServerForString(siteName)
	}
	ketama = recreateKetama(hashToServerFile.Name(), segmentId)
	for siteName, resultOne := range stringToHash {
		resultTwo := ketama.getServerForString(siteName)
		if resultOne != resultTwo {
			t.Errorf("recreated Ring not returning expected server, expected %s, received %s", resultOne, resultTwo)
		}
	}
}
