package main

import (
	"encoding/json"
	"flag"
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

var (
	help       = flag.Bool("h", false, "This help.")
	verbose    = flag.Bool("v", false, "Verbose mode.")
	baseURL    = flag.String("base", "http://183.207.248.71:80/cntv/live1", "Base URL for stream.")
	apiURL     = flag.String("api", "http://looktvepg.jsa.bcs.ottcn.com:8080/ysten-lvoms-epg/epg/getChannelIndexs.shtml?deviceGroupId=1697", "API URL to fetch channel list.")
	outputFile = flag.String("o", "channel.m3u8", "Output file path.")
)

type Channel struct {
	UUID        string `json:"uuid"`
	ChannelName string `json:"channelName"`
	ChannelIcon string `json:"channelIcon"`
}

func (c *Channel) toString() string {
	u, _ := url.Parse(*baseURL + "/" + c.ChannelName + "/" + c.UUID)
	return fmt.Sprintf("#EXTINF:-1,%s\n%s\n", c.ChannelName, u.String())
}

func getJSONContent() map[string]Channel {
	fmt.Println("Fetching channel list...")

	resp, err := http.Get(*apiURL)
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
	fmt.Println("Parsing...")

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
	fmt.Println("Generating m3u8...")

	content := fmt.Sprintf("#EXTM3U\n\n")
	for _, c := range *channels {
		content += c.toString()
		if *verbose {
			fmt.Println(c.ChannelName)
		}
	}

	f, err := os.Create(*outputFile)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	_, err = io.WriteString(f, content)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Done, all saved to %s\n", *outputFile)
}

func main() {
	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	j := getJSONContent()
	l := generateChannelList(&j)
	generateM3U8(&l)
}
