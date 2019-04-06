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

type apiResponse struct {
	Error string
	Data  interface{}
}

func main() {
	_, err := svdrp.ListAllChannels("vdr.dd:6419")
	if err != nil {
		fmt.Println("ERROR:", err)
	}

	EPGNow, err := svdrp.ListEPGNow("vdr.dd:6419")
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
	_, err = svdrp.ListEPGNow("vdr.dd:6419")
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
	channels, err := svdrp.ListAllChannels("vdr.dd:6419") //get list of all channels
	if err != nil {
		errResponse := handleAPIError(response, "Could not List Channels", "error while marshaling error-json on channel-listing: %v\n")
		if errResponse {
			fmt.Println("return error")
			return
		}
	}
	marshChannels, err := json.Marshal(channels)
	if err != nil {
		fmt.Fprintln(response, "could not marshal the channel list")
	}
	response.Header().Add("Content-type", "application/json")
	fmt.Fprintln(response, string(marshChannels))
}

func apiEPGHandler(response http.ResponseWriter, r *http.Request) {
	epgData, err := svdrp.ListEPGNow("vdr.dd:6419") //get list of all channels
	if err != nil {
		errResponse := handleAPIError(response, "Could not List EPGData", "error while marshaling error-json on EPGNow-listing: %v\n")
		if errResponse {
			fmt.Println("return error")
			return
		}
	}
	marshEPG, err := json.Marshal(epgData)
	if err != nil {
		fmt.Fprintln(response, "could not marshal the EPGData")
	}
	response.Header().Add("Content-type", "application/json")
	fmt.Fprintln(response, string(marshEPG))
}

func handleAPIError(response http.ResponseWriter, apiErrMassage, jsonErrMassage string) bool {
	fmt.Println(apiErrMassage)
	marshErrAPI, jsonErr := json.Marshal(apiResponse{apiErrMassage, nil})
	if jsonErr != nil {
		fmt.Printf(apiErrMassage, jsonErr)
		response.WriteHeader(http.StatusInternalServerError)
		return true
	}
	response.Header().Add("Content-type", "application/json")
	fmt.Fprintln(response, string(marshErrAPI))
	return true
}
