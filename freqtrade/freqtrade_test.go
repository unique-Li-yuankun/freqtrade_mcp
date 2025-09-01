package freqtrade

import (
	"flag"
	"freqtrade_mcp/utils"
	"github.com/stretchr/testify/require"
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

func TestStructJsonParam(t *testing.T) {
	testStruct := struct {
		No          int      `json:"no"`
		Description string   `json:"description"`
		Price       float64  `json:"price"`
		Labels      []string `json:"labels"`
	}{
		No:          1,
		Description: "Coin",
		Price:       10,
		Labels:      []string{"stable", "pow"},
	}
	p := utils.StructJsonParams(&testStruct)
	require.Equal(t, p, []string{"--no 1", "--description Coin", "--price 10", "--labels stable pow"})
}
