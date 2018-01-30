package main

import "testing"

func TestIsLastLine(t *testing.T) {

	// testcase 1
	lastLine := isLastLine("250 78 3sat;ZDFvision:11953:HC34M2S0:S19.2E:27500:210=2:220=deu@3,221=mis@3,222=mul@3;225=deu@106:230;231=deu:0:28007:1:1079:0")
	if lastLine == false {
		t.Error("lastLine should be detected as last line but is not")
	}

	// testcase 2
	normalLine := isLastLine("250-77 ORF III;ORF:12662:HC56M2S0:S19.2E:22000:1010=2:1011=deu@4:1013:648,650,D95,D98,9C4,98C:13101:1:1115:0")
	if normalLine == true {
		t.Error("normalLine should not be detected as last line but is")
	}
}
