package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
)

var CHANNEL_API_URL = "http://looktvepg.jsa.bcs.ottcn.com:8080/ysten-lvoms-epg/epg/getChannelIndexs.shtml?deviceGroupId=1697"
var BASE_URL = "http://183.207.248.71:80/cntv/live1"

type Channel struct {
	UUID        string `json:"uuid"`
	ChannelName string `json:"channelName"`
	ChannelIcon string `json:"channelIcon"`
}

func (c *Channel) toString() string {
	us := fmt.Sprintf("%s/%s/%s", BASE_URL, url.QueryEscape(c.ChannelName), url.QueryEscape(c.UUID))
	str := fmt.Sprintf("#EXTINF:-1,%s\n%s\n", c.ChannelName, us)
	return str
}

func getJSONContent() map[string]Channel {
	resp, err := http.Get(CHANNEL_API_URL)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("api error: %s", resp.Status)
	}

	var result map[string]Channel
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatalln(err)
	}
	return result
}

func generateChannelList(json *map[string]Channel) []Channel {
	keySort := make(map[int]string, len(*json))
	for key := range *json {
		number := stringArrayToString(regexp.MustCompile("[0-9]").FindAllStringSubmatch(key, -1))
		i, err := strconv.Atoi(number)
		if err != nil {
			log.Fatalf("failed to parse channel list: %v", err)
		}
		keySort[i] = key
	}

	subKeys := make([]int, 0, len(keySort))
	for key := range keySort {
		subKeys = append(subKeys, key)
	}
	sort.Ints(subKeys)

	channels := make([]Channel, 0, len(keySort))
	for _, subKey := range subKeys {
		key := keySort[subKey]
		channels = append(channels, (*json)[key])
	}

	return channels
}

func stringArrayToString(array [][]string) string {
	line := ""
	for _, sa := range array {
		for _, s := range sa {
			line += s
		}
	}
	return line
}

func generateM3U8(channels *[]Channel) {
	content := fmt.Sprintf("#EXTM3U\n\n")
	for _, c := range *channels {
		content += c.toString()
		fmt.Println(c.ChannelName)
	}

	f, err := os.Create("channel.m3u8")
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	_, err = io.WriteString(f, content)
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	fmt.Println("Fetching channel list...")
	j := getJSONContent()
	fmt.Println("Parsing...")
	list := generateChannelList(&j)
	fmt.Println("Generating m3u8...")
	generateM3U8(&list)
	fmt.Println("Done, all saved to channel.m3u8")
}
