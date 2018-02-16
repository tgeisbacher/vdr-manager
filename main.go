package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/tgeisbacher/vdr-manager/svdrp"
)

func main() {
	_, err := svdrp.ListAllChannels("vdr.dd:6419")

	if err != nil {
		fmt.Println("ERROR:", err)
	}

	http.HandleFunc("/", startPageHandler)
	err = http.ListenAndServe(":34973", nil) // listen to all who select the right port
	if err != nil {
		fmt.Println("ERROR:", err)
	}
}

func startPageHandler(response http.ResponseWriter, r *http.Request) {
	content := "<h1>VDR Manager</h1>"
	channels, err := svdrp.ListAllChannels("vdr.dd:6419")
	if err != nil {
		fmt.Println("Error:", err)
	}
	content = content + `<div class="alert alert-danger" role="alert">
	This asdfasdfasdfasdfasdf
  </div>`
	content = content + "<ul>"
	for _, chann := range channels {
		content = content + fmt.Sprintf("<li>%d - %s</li>", chann.Position, chann.Name)
	}
	content = content + "</ul>"

	renderHTML(response, content)
}

func renderHTML(response http.ResponseWriter, content string) {
	htmlHeader, err := ioutil.ReadFile("html/header.html")
	if err != nil {
		fmt.Fprintln(response, "Could not load header")
		return
	}
	htmlFooter, err := ioutil.ReadFile("html/footer.html")
	if err != nil {
		fmt.Fprintln(response, "Could not load footer")
		return
	}

	fmt.Fprintln(response, string(htmlHeader)+content+string(htmlFooter))

}
