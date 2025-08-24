package freqtrade

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"reflect"
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
	Config          string   `json:"config" jsonschema:"Specify configuration file, path must be absolute where the config.json located"`
	StrategyPath    string   `json:"strategy-path" jsonschema:"Path to the strategy, directory must be absolute where the .py located"`
}

func (p *BackTestParams) Param() []string {
	params := structJsonParams(p)
	params = append(params, fmt.Sprintf("--data-format-ohlcv %s", "json"))
	params = append(params, fmt.Sprintf("--export %s", "trades"))
	return params
}

type DownloadDataParams struct {
	Exchange  string   `json:"exchange" jsonschema:"Exchange to download data from"`
	Timeframe string   `json:"timeframe" jsonschema:"Timeframe to download data from"`
	Pairs     []string `json:"pairs" jsonschema:"Pairs to download data from, (example:'BTC/USDT')"`
	TimeRange string   `json:"timerange" jsonschema:"Time range to download data from, (example:'20240101-20240102')"`
}

func (p *DownloadDataParams) Param() []string {
	params := structJsonParams(p)
	params = append(params, fmt.Sprintf("--data-format-ohlcv %s", "json"))
	return params
}

func structJsonParams(s any) []string {
	var params []string
	ts := reflect.TypeOf(s).Elem()
	vs := reflect.ValueOf(s).Elem()
	for i := 0; i < ts.NumField(); i++ {
		field := ts.Field(i)
		fieldV := vs.Field(i)
		key, ok := field.Tag.Lookup("json")
		if !ok {
			continue
		}

		var param string
		switch field.Type.Kind() {
		case reflect.Slice:
			var value []string
			for j := 0; j < fieldV.Len(); j++ {
				value = append(value, fmt.Sprintf("%v", fieldV.Index(j).Interface()))
			}
			param = fmt.Sprintf("--%s %v", key, strings.Join(value, " "))
		default:
			param = fmt.Sprintf("--%s %v", key, fieldV.Interface())
		}
		params = append(params, param)
	}
	return params
}

func DownloadData(p DownloadDataParams) (string, error) {
	output, err := ExecuteCommandInNewConsole("freqtrade download-data", p.Param()...)
	return string(output), err
}

func BackTest(p BackTestParams) (string, error) {
	var result string
	trades, err := ExecuteCommandInNewConsole("freqtrade backtesting", p.Param()...)
	result += string(trades)
	if err != nil {
		return result, err
	}
	return result, nil
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
