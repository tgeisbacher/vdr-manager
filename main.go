package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {

	fmt.Println("Hello VDR!")

	conn, err := net.Dial("tcp", "10.0.0.198:6419")

	if err != nil {
		fmt.Println("ERROR: ", err)
		os.Exit(1)
	}
	defer conn.Close()

	fromServReader := bufio.NewReader(conn)
	fromServReader.ReadString('\n') // einlesen der Willkommensnachricht

	fmt.Fprintf(conn, "LSTC\n") //Fprintf -> text in beliebigen Writer(hier: conn)

	s := bufio.NewScanner(fromServReader)
	s.Split(bufio.ScanLines)
	for s.Scan() {
		line := s.Text()
		fmt.Printf("%s\n", line)

		if isLastLine(line) == true {
			break
		}

	}
}

func isLastLine(line string) bool {

	if strings.Index(line, "-") == 3 {
		return false
	}

	return true

}
