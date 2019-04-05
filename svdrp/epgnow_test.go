package svdrp

import (
	"testing"
	"time"
)

var eurosport string = `215-C S28.2E-1-1091-31200 Eurosport Deutschland OBSOLETE
215-E 8882 1519111824 7800 4E 3
215-T Eishockey: Männer, Play-offs
215-S Slowenien - Norwegen
215-D Bei den 23. Olympischen Winterspielen in Pyeongchang (KOR) werden im Eishockey im Gangneung Hockey Centre zwei Wettbewerbe ausgetragen: das Turnier der Männer mit zwölf Mannschaften und das Turnier der Frauen mit acht Mannschaften. Das Teilnehmerfeld ergibt sich aus der IIHF-Weltrangliste sowie durch Qualifikationsturniere, wodurch sich auch die deutsche Männermannschaft das Ticket nach Südkorea sichern konnte. Bei den Spielen 2014 in Sotschi konnte Kanada sich beide Goldmedaillen sichern.
215-e`
var testResponse string = `215-C S19.2E-1-1082-20003 Kabel 1 Schweiz
215-c
215-C S19.2E-1-1082-20005 SAT.1 A
215-c
` + eurosport + `
215-c`

func TestParseEPGEntryResponse(t *testing.T) {
	epgEntry := parseEPGEntryResponse(eurosport)

	if epgEntry.ChannelID != "S28.2E-1-1091-31200" {
		t.Errorf("IDOfChannelFromEPG should be 'S28.2E-1-1091-31200' but is '%s'", epgEntry.ChannelID)
	}
	if epgEntry.ChannelName != "Eurosport Deutschland OBSOLETE" {
		t.Errorf("nameOfChannelFromEPG should be 'Eurosport Deutschland' but is '%s'", epgEntry.ChannelName)
	}
	//fmt.Println("TIME IS:", time.Unix(1519111824, 0))
	if !time.Unix(1519111824, 0).Equal(epgEntry.StartTime) {
		t.Errorf("Starttime should be '2018-02-20 08:30' but is '%v'", epgEntry.StartTime)
	}
	expectedDuration, err := time.ParseDuration("7800s")
	if err != nil {
		t.Error("Could not parse '7800s'")
	}
	if epgEntry.Duration.Nanoseconds() != expectedDuration.Nanoseconds() {
		t.Errorf("Duration should be 130min but is '%v'", epgEntry.Duration)
	}
	if epgEntry.ShowName != "Eishockey: Männer, Play-offs" {
		t.Errorf("nameOfShowFromEPG should be 'Eishockey: Männer, Play-offs' but is '%s'", epgEntry.ShowName)
	}
	if epgEntry.Subtitle != "Slowenien - Norwegen" {
		t.Errorf("Subtitle should be 'Slowenien - Norwegen' but is '%s'", epgEntry.Subtitle)
	}

}

func TestSplitEPGNowByChannel(t *testing.T) {
	channelBlocks := splitEPGNowByChannel(testResponse)

	if len(channelBlocks) != 3 {
		t.Errorf("Length should be 3 but is '%d'", len(channelBlocks))
	}

	if channelBlocks[0] != "215-C S19.2E-1-1082-20003 Kabel 1 Schweiz" {
		t.Errorf("First Channelblock should be '...Kabel 1 Schweiz' but is '%s'", channelBlocks[0])
	}
	if channelBlocks[1] != "215-C S19.2E-1-1082-20005 SAT.1 A" {
		t.Errorf("Second Channelblock should be '...SAT 1' but is '%s'", channelBlocks[1])
	}
	//fmt.Println("CONTENT OF CB[2]:", channelBlocks[2])
	if len(channelBlocks[2]) != len(eurosport) {
		t.Errorf("Third Channelblock is not valid")
	}
}
