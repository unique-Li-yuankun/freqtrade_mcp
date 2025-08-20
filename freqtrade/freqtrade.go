package freqtrade

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var Dir string

type BackTestParams struct {
	TimeFrame       string   `json:"timeframe" jsonschema:"Specify timeframe 1m,5m,30m,1h,1d"`
	TimeRange       string   `json:"timerange" jsonschema:"Specify what timerange of data to use"`
	MaxOpenTrades   int      `json:"max-open-trades" jsonschema:"Override the value of the max_open_trades configuration setting"`
	StakeAmount     int      `json:"stake-amount" jsonschema:"Override the value of the stake_amount configuration setting"`
	Pairs           []string `json:"pairs" jsonschema:"Limit command to these pairs"`
	StartingBalance int      `json:"starting-balance" jsonschema:"Starting balance, used for backtesting / hyperopt and dry-runs"`
	StrategyList    []string `json:"strategy-list" jsonschema:"Strategy that need to be backtested"`
}

func (p *BackTestParams) Param() []string {
	var params []string
	if p.TimeFrame == "1m" || p.TimeFrame == "5m" || p.TimeFrame == "30m" || p.TimeFrame == "1h" || p.TimeFrame == "1d" {
		params = append(params, fmt.Sprintf("--timeframe %s", p.TimeFrame))
	}
	if p.TimeRange != "" {
		params = append(params, fmt.Sprintf("--timerange %s", p.TimeRange))
	}
	if p.MaxOpenTrades != 0 {
		params = append(params, fmt.Sprintf("--max-open-trades %d", p.MaxOpenTrades))
	}
	if p.StakeAmount != 0 {
		params = append(params, fmt.Sprintf("--stake-amount %d", p.StakeAmount))
	}
	if len(p.Pairs) != 0 {
		pairs := strings.Join(p.Pairs, " ")
		params = append(params, fmt.Sprintf("--pairs %s", pairs))
	}
	if p.StartingBalance != 0 {
		params = append(params, fmt.Sprintf("--starting-balance %d", p.StartingBalance))
	}
	if len(p.StrategyList) != 0 {
		strategies := strings.Join(p.StrategyList, " ")
		params = append(params, fmt.Sprintf("--strategy-list %s", strategies))
	}
	return params
}

func BackTest(p BackTestParams) (string, error) {
	output, err := ExecuteCommandInNewConsole("freqtrade backtesting", p.Param()...)
	return string(output), err
}

func ExecuteCommandInNewConsole(command string, args ...string) ([]byte, error) {
	var cmd *exec.Cmd
	var fullCommand string

	userCommand := command
	if len(args) > 0 {
		userCommand += " " + strings.Join(args, " ")
	}

	cdCommand := fmt.Sprintf("cd %s", Dir)
	if runtime.GOOS == "windows" {
		activateScript := filepath.Join(Dir, ".venv", "Scripts", "Activate.ps1")
		fullCommand = fmt.Sprintf(`& "%s"; %s; %s`, activateScript, cdCommand, userCommand)
		cmd = exec.Command("powershell", "-Command", fullCommand)
	} else {
		activateScript := filepath.Join(Dir, ".venv", "bin", "activate")
		fullCommand = fmt.Sprintf("source %s && %s && %s", activateScript, cdCommand, userCommand)
		cmd = exec.Command("/bin/bash", "-c", fullCommand)
	}

	return cmd.CombinedOutput()
}
