package freqtrade

import (
	"strings"
	"testing"
)

func init() {
	Dir = "W:\\freqtrade"
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

func TestBackTesting(t *testing.T) {
	p := BackTestParams{
		TimeFrame:       "5m",
		TimeRange:       "20250101-20250801",
		MaxOpenTrades:   3,
		StakeAmount:     500,
		Pairs:           []string{"BTC/USDT"},
		StartingBalance: 10000,
		StrategyList:    []string{"SampleStrategy"},
	}
	output, err := BackTest(p)
	t.Log(output, err)
}
