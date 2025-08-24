package freqtrade

import (
	"flag"
	"strings"
	"testing"
)

func init() {
	flag.StringVar(&Dir, "dir", "", "freqtrade directory")
}

func TestFreqtradeEnv(t *testing.T) {
	output, err := ExecuteCommandInNewConsole("freqtrade")
	if err == nil {
		t.Error("Execute freqtrade command error.")
	}
	if !strings.Contains(string(output), "To see the full list of options available, please use `freqtrade --help` or `freqtrade <command> --help`.") {
		t.Error("Execute freqtrade command error.")
	}
}
