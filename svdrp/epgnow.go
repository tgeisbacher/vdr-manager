package svdrp

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type EPGEntry struct {
	ChannelID   string
	ChannelName string
	StartTime   time.Time
	Duration    time.Duration
	ShowName    string
	Subtitle    string
}

// ListEPGNow list all programs which runs currently on TV
func ListEPGNow(addr string) ([]EPGEntry, error) {
	response, err := collectDataFromServer(addr, "LSTE now")
	if err != nil {
		return nil, fmt.Errorf("ERROR:Could not list EPG now from VDR %v", err)
	}
	epgChannelBlocks := splitEPGNowByChannel(response)
	fullEPGEntries, err := parseEPGResponseBlocks(epgChannelBlocks)
	if err != nil {
		return nil, fmt.Errorf("could not parse response to Channelblocks:", err)
	}
	return fullEPGEntries, nil

}

func splitEPGNowByChannel(response string) []string {
	responseReader := strings.NewReader(response)
	s := bufio.NewScanner(responseReader)
	s.Split(bufio.ScanLines)
	fullChannelBlocks := []string{}
	fullChannelBlockString := ""
	for s.Scan() {
		line := s.Text()
		if strings.Index(line, "c") == 4 {
			fullChannelBlocks = append(fullChannelBlocks, strings.TrimSpace(fullChannelBlockString))
			fullChannelBlockString = ""
		} else {
			fullChannelBlockString = fullChannelBlockString + line + "\n"
		}
	}
	return fullChannelBlocks
}

func parseEPGResponseBlocks(blocks []string) ([]EPGEntry, error) {
	parsedChannelBlock := []EPGEntry{}
	for i := 0; i < len(blocks); i++ {
		responseEPGBlock, err := parseEPGResponseBlock(blocks[i])
		if err != nil {
			return nil, fmt.Errorf("ERROR:", err)
		}
		parsedChannelBlock = append(parsedChannelBlock, responseEPGBlock)
	}
	return parsedChannelBlock, nil
}

func parseEPGResponseBlock(responseEPG string) (EPGEntry, error) {
	readerResponse := strings.NewReader(responseEPG)
	s := bufio.NewScanner(readerResponse)
	s.Split(bufio.ScanLines)
	epgData := EPGEntry{}
	for s.Scan() {
		line := s.Text()
		if strings.Index(line, "C") == 4 {
			channelLineWithoutSpaceBetweenWord := strings.Split(line, " ")
			fullChannelName := ""
			for i := 2; i < len(channelLineWithoutSpaceBetweenWord); i++ {
				fullChannelName = fullChannelName + channelLineWithoutSpaceBetweenWord[i] + " "
			}
			epgData.ChannelName = strings.TrimSpace(fullChannelName)
			epgData.ChannelID = channelLineWithoutSpaceBetweenWord[1]

		}

		if strings.Index(line, "T") == 4 {
			titleLineWithoutDash := strings.Split(line, " ")
			fullTitle := ""
			for i := 1; i < len(titleLineWithoutDash); i++ {
				fullTitle = fullTitle + titleLineWithoutDash[i] + " "
			}

			epgData.ShowName = strings.TrimSpace(fullTitle)
		}

		if strings.Index(line, "E") == 4 {
			epgLineWithoutDash := strings.Split(line, " ")
			epgStartTimeString := strings.TrimSpace(epgLineWithoutDash[2])
			epgDurationString := strings.TrimSpace(epgLineWithoutDash[3])
			epgStartTime, err := strconv.ParseInt(epgStartTimeString, 0, 64)
			if err != nil {
				return EPGEntry{}, fmt.Errorf("Can not convert starttime string to integer")
			}
			epgStartTimeUnix := time.Unix(epgStartTime, 0)
			epgDurationUnix, err := time.ParseDuration(epgDurationString + "s")
			if err != nil {
				return EPGEntry{}, fmt.Errorf("Could not convert Duration string to time.Duration")
			}
			epgData.StartTime = epgStartTimeUnix
			epgData.Duration = epgDurationUnix
		}

		if strings.Index(line, "S") == 4 {
			epgSubtitleOfShow := strings.Split(line, " ")
			fullSubtitleName := ""
			for i := 1; i < len(epgSubtitleOfShow); i++ {
				fullSubtitleName = fullSubtitleName + epgSubtitleOfShow[i] + " "
			}
			epgData.Subtitle = strings.TrimSpace(fullSubtitleName)
		}

	}

	return epgData, nil
}
