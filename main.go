package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/tgeisbacher/vdr-manager/svdrp"
)

const vdrHostname = "vdr.dd:6419"

type apiResponse struct {
	Error string
	Data  interface{}
}

func main() {
	_, err := svdrp.ListAllChannels(vdrHostname)
	if err != nil {
		fmt.Println("ERROR:", err)
	}

	EPGNow, err := svdrp.ListEPGNow(vdrHostname)
	if err != nil {
		fmt.Println("ERROR:", err)
	}
	for _, data := range EPGNow {
		fmt.Println("ChannelID: ", data.ChannelID)
		fmt.Println("Channelname: ", data.ChannelName)
		fmt.Println("Showname: ", data.ShowName)
		fmt.Println("Starttime:", data.StartTime, "Duration:", data.Duration)
		fmt.Println("Subtitle: ", data.Subtitle)
	}
	_, err = svdrp.ListEPGNow(vdrHostname)
	if err != nil {
		fmt.Println("ERROR:", err)
	}

	fs := http.FileServer(http.Dir("html"))
	router := mux.NewRouter().StrictSlash(true)
	router.Methods("GET").Path("/api/channels").HandlerFunc(apiChannelsHandler)
	router.Methods("GET").Path("/api/epg").HandlerFunc(apiEPGHandler)
	router.PathPrefix("/").Handler(fs)
	handler := cors.Default().Handler(router)

	address := ":34973"
	fmt.Printf("server listening on %v ...", address)
	log.Fatal(http.ListenAndServe(address, handler))
}

func apiChannelsHandler(response http.ResponseWriter, r *http.Request) {
	channels, err := svdrp.ListAllChannels(vdrHostname) //get list of all channels
	if err != nil {
		handleAPIError(response, err, "Could not List Channels", "error while marshaling error-json on channel-listing: %v\n")
		return
	}
	marshChannels, err := json.Marshal(apiResponse{Error: "", Data: channels})
	if err != nil {
		handleAPIError(response, err, "could not marshal the channel list", "")
		return
	}
	response.Header().Add("Content-type", "application/json")
	fmt.Fprintln(response, string(marshChannels))
}

func apiEPGHandler(response http.ResponseWriter, r *http.Request) {
	epgData, err := svdrp.ListEPGNow(vdrHostname) //get list of all channels
	if err != nil {
		handleAPIError(response, err, "Could not List EPGData", "error while marshaling error-json on EPGNow-listing: %v\n")
		return
	}
	marshEPG, err := json.Marshal(apiResponse{Error: "", Data: epgData})
	if err != nil {
		handleAPIError(response, err, "could not marshal the EPGData", "")
		return
	}
	response.Header().Add("Content-type", "application/json")
	fmt.Fprintln(response, string(marshEPG))
}

func handleAPIError(response http.ResponseWriter, err error, apiErrMessage, jsonErrMessage string) {
	fmt.Printf(apiErrMessage+": %v\n", err)
	marshErrAPI, jsonErr := json.Marshal(apiResponse{Error: apiErrMessage, Data: nil})
	if jsonErr != nil {
		fmt.Printf(jsonErrMessage, jsonErr)
		response.WriteHeader(http.StatusInternalServerError)
	}
	response.Header().Add("Content-type", "application/json")
	fmt.Fprintln(response, string(marshErrAPI))
}
