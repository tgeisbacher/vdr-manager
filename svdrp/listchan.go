package svdrp

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

type channel struct {
	Position int
	Name     string
	Group    string
}

func ListAllChannels(addr string) ([]channel, error) {

	response, err := collectDataFromServer(addr, "LSTC")
	if err != nil {
		return nil, fmt.Errorf("Could not List All Channels Error: %v", err)
	}

	vDRChannel := parseListChannelsResponse(response)

	return vDRChannel, nil
}

func parseListChannelsResponse(response string) []channel {
	readResponse := strings.NewReader(response) // declare a new Reader-> Reads the Response String
	s := bufio.NewScanner(readResponse)         // declare a newScanner from Reader -> Scan the reader
	s.Split(bufio.ScanLines)                    // Scan the resonse from line to line
	vdrChannels := []channel{}

	for s.Scan() {
		line := s.Text()
		splitLine := []string{}

		//last line of response has no "-" between statuscode and programposition; detect last line
		if isLastLine(line) == true {
			splitLine = strings.SplitN(line, " ", 2)
		} else {
			splitLine = strings.Split(line, "-")
		}
		lineWithoutStatusCode := strings.Split(splitLine[1], ";")
		partWithPositionAndName := strings.Split(lineWithoutStatusCode[0], " ")
		position, err := strconv.Atoi(partWithPositionAndName[0]) //Position is a string -> here we convert position(string) in position(int)
		if err != nil {
			fmt.Errorf("Invalid Position: %s on line: %s", partWithPositionAndName[0], line)
			continue
		}

		vdrChannel := channel{}
		vdrChannel.Position = position
		//TODO: hr - fernsehen, n-tv
		fullChannelName := ""
		for i := 1; i < len(partWithPositionAndName); i++ {
			fullChannelName = fullChannelName + partWithPositionAndName[i] + " "
		}
		vdrChannel.Name = strings.TrimSpace(fullChannelName)

		// TODO: groupName auslesen
		vdrChannels = append(vdrChannels, vdrChannel)

	}

	return vdrChannels
}
