package svdrp

import "fmt"

func ListAllChannels(addr string) (string, error) {

	response, err := collectDataFromServer(addr, "LSTC")

	if err != nil {
		return "", fmt.Errorf("Input not correct! Could not List All Channels Error: ", err)
	}

	return response, nil
}
