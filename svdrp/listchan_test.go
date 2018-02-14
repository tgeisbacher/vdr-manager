package svdrp

import "testing"

func TestParseListMultipleChannelsResponse(t *testing.T) {
	responseFromVDRServer := `250-1 ORF1 HD;ORF:11302:HC23M5O35P0S1:S19.2E:22000:1920=27:0;1921=deu@106,1922=mis@106:1925:648,650,D95,D98,6E2,500,9C4,98C:49
250-2 ORF2St HD;ORF:11273:hC23M5O35S1:S19.2E:22000:3010=27:0;3011=deu@106:3005:D98,650,D95,648,6E2,98C,9C4,500:13301:1:1005:0
250 15 DMAX;BetaDigital:12480:VC34M2S0:S19.2E:27500:3327=2:3328=deu@3:44:0:63:133:33:0`

	parsedChannels := parseListChannelsResponse(responseFromVDRServer)

	if len(parsedChannels) != 3 {
		t.Error("Length of parsedChannels should be 3 but is ", len(parsedChannels))
	}

	if parsedChannels[0].name != "ORF1 HD" {
		t.Errorf("Channelname should be ORF1 HD but is '%s'", parsedChannels[0].name)
	}
	if parsedChannels[0].position != 1 {
		t.Error("Channelposition should be 1 but is", parsedChannels[0].position)
	}
	/*if parsedChannels[0].group != "ORF" {
		t.Error("Channelgroup should be ORF but is", parsedChannels[0].group)
	}
	*/
	if parsedChannels[1].name != "ORF2St HD" {
		t.Error("Channelname should be ORF2St HD but is ", parsedChannels[1].name)
	}

	if parsedChannels[1].position != 2 {
		t.Error("Channelposition should be 2 but is", parsedChannels[1].position)
	}

	/*	if parsedChannels[1].group != "ORF" {
		t.Error("Channelgroup should be ORF but is", parsedChannels[1].group)
	}*/

	//Testing third line
	if parsedChannels[2].name != "DMAX" {
		t.Error("Channelname should be DMAX but is ", parsedChannels[1].name)
	}

	if parsedChannels[2].position != 15 {
		t.Error("Channelposition should be 15 but is", parsedChannels[1].position)
	}

	/*	if parsedChannels[1].group != "ORF" {
		t.Error("Channelgroup should be ORF but is", parsedChannels[1].group)
	}*/

}
