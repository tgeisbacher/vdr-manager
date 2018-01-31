package svdrp

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func collectDataFromServer(addr, command string) (string, error) {

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return "", fmt.Errorf("error while connecting to server '%s', error: %v", addr, err)
	}
	defer conn.Close()

	// reading welcome-message from server
	fromServReader := bufio.NewReader(conn)
	_, err = fromServReader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("could not read welcome-message from server, error: %v", err)
	}

	// sending command
	fmt.Fprint(conn, command+"\n")

	// receiving answer from server
	response := ""
	s := bufio.NewScanner(fromServReader)
	s.Split(bufio.ScanLines)
	for s.Scan() {
		line := s.Text()
		response = response + line + "\n"
		if isLastLine(line) == true {
			break
		}
	}

	return response, nil
}

func isLastLine(line string) bool {
	if strings.Index(line, "-") == 3 {
		return false
	}

	return true
}
