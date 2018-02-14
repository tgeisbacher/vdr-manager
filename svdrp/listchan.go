package svdrp

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

type channel struct {
	position int
	name     string
	group    string
}

func ListAllChannels(addr string) (string, error) {

	response, err := collectDataFromServer(addr, "LSTC")

	if err != nil {
		return "", fmt.Errorf("Input not correct! Could not List All Channels Error: ", err)
	}

	return response, nil
}

func parseListChannelsResponse(response string) []channel {
	readResponse := strings.NewReader(response) // declare a new Reader-> Reads the Response String
	s := bufio.NewScanner(readResponse)         // declare a newScanner from Reader -> Scan the reader
	s.Split(bufio.ScanLines)                    // Scan the resonse from line to line
	vdrChannels := []channel{}

	for s.Scan() {
		line := s.Text()
		splitLine := strings.Split(line, "-")
		lineWithoutStatusCode := strings.Split(splitLine[1], ";")
		partWithPositionAndName := strings.Split(lineWithoutStatusCode[0], " ")
		position, err := strconv.Atoi(partWithPositionAndName[0]) //Position is a string -> here we convert position(string) in position(int)
		if err != nil {
			fmt.Errorf("Invalid Position: %s on line: %s", partWithPositionAndName[0], line)
			continue
		}

		vdrChannel := channel{}
		vdrChannel.position = position

		fullChannelName := ""
		for i := 1; i < len(partWithPositionAndName); i++ {
			fullChannelName = fullChannelName + partWithPositionAndName[i] + " "
		}
		vdrChannel.name = strings.TrimSpace(fullChannelName)

		// TODO: groupName auslesen
		vdrChannels = append(vdrChannels, vdrChannel)

	}

	return vdrChannels
}
