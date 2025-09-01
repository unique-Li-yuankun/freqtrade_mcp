package freqtrade

import (
	"fmt"
	"freqtrade_mcp/utils"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	Dir string
)

type BacktesingParams struct {
	TimeFrame       string   `json:"timeframe" jsonschema:"Specify timeframe 1m,5m,30m,1h,1d"`
	TimeRange       string   `json:"timerange" jsonschema:"Specify what timerange of data to use"`
	MaxOpenTrades   int      `json:"max-open-trades" jsonschema:"Override the value of the max_open_trades configuration setting"`
	StakeAmount     int      `json:"stake-amount" jsonschema:"Override the value of the stake_amount configuration setting"`
	Pairs           []string `json:"pairs" jsonschema:"Limit command to these pairs"`
	StartingBalance int      `json:"starting-balance" jsonschema:"Starting balance, used for backtesting / hyperopt and dry-runs"`
	StrategyList    []string `json:"strategy-list" jsonschema:"Strategy that need to be backtested. Make sure you have already upload them"`
	UserDir         string   `json:"userdir" jsonschema:"Path to userdata directory. Relative path"`
}

func (p *BacktesingParams) Param() []string {
	params := utils.StructJsonParams(p)
	params = append(params, fmt.Sprintf("--data-format-ohlcv %s", "json"))
	params = append(params, fmt.Sprintf("--export %s", "signals"))
	return params
}

type DownloadDataParams struct {
	Exchange  string   `json:"exchange" jsonschema:"Exchange to download data from"`
	Timeframe string   `json:"timeframe" jsonschema:"Timeframe to download data from"`
	Pairs     []string `json:"pairs" jsonschema:"Pairs to download data from, (example:'BTC/USDT')"`
	TimeRange string   `json:"timerange" jsonschema:"Time range to download data from, (example:'20240101-20240102')"`
	UserDir   string   `json:"userdir" jsonschema:"Path to userdata directory. Relative path"`
}

func (p *DownloadDataParams) Param() []string {
	params := utils.StructJsonParams(p)
	params = append(params, fmt.Sprintf("--data-format-ohlcv %s", "json"))
	params = append(params, "--erase")
	return params
}

type BacktestingAnalysisParams struct {
	UserDir string `json:"userdir" jsonschema:"Path to userdata directory. Relative path"`
}

func (p *BacktestingAnalysisParams) Param() []string {
	params := utils.StructJsonParams(p)
	params = append(params, fmt.Sprintf("--indicator-list %s", "close_date trade_duration amount profit_ratio profit_abs orders"))
	return params
}

type CreateUserDirParams struct {
	UserDir string `json:"userdir" jsonschema:"Path to userdata directory. Relative path (example: tt)"`
}

func (p *CreateUserDirParams) Param() []string {
	return utils.StructJsonParams(p)
}

func DownloadData(p DownloadDataParams) (string, error) {
	output, err := ExecuteCommandInNewConsole("freqtrade download-data", p.Param()...)
	return string(output), err
}

func Backtesting(p BacktesingParams) (string, error) {
	output, err := ExecuteCommandInNewConsole("freqtrade backtesting", p.Param()...)
	return string(output), err
}

func BacktestingAnalysis(p BacktestingAnalysisParams) (string, error) {
	output, err := ExecuteCommandInNewConsole("freqtrade backtesting-analysis", p.Param()...)
	return string(output), err
}

func CreateUserDir(p CreateUserDirParams) (string, error) {
	output, err := ExecuteCommandInNewConsole("freqtrade create-userdir", p.Param()...)
	return string(output), err
}

func ExecuteCommandInNewConsole(command string, args ...string) ([]byte, error) {
	var fullCommand string
	userCommand := command
	if len(args) > 0 {
		userCommand += " " + strings.Join(args, " ")
	}
	cdCommand := fmt.Sprintf("cd %s", Dir)
	if runtime.GOOS == "windows" {
		activateScript := filepath.Join(Dir, ".venv", "Scripts", "Activate.ps1")
		fullCommand = fmt.Sprintf(`& "%s"; %s; %s`, activateScript, cdCommand, userCommand)
	} else {
		activateScript := filepath.Join(Dir, ".venv", "bin", "activate")
		fullCommand = fmt.Sprintf("source %s && %s && %s", activateScript, cdCommand, userCommand)
	}
	return utils.ExecuteCommand(fullCommand)
}
