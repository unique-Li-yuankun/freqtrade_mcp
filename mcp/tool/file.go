package tool

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"freqtrade_mcp/freqtrade"
	"freqtrade_mcp/utils"
	jsoniter "github.com/json-iterator/go"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type UpsertStrategyParams struct {
	FileName string `json:"filename" json:"Strategy filename"`
	Strategy string `json:"strategy" jsonschema:"Strategy code"`
	UserDir  string `json:"userdir" jsonschema:"Path to userdata directory"`
}

func (p *UpsertStrategyParams) Param() []string {
	return utils.StructJsonParams(p)
}

func UpsertStrategy(ctx context.Context, session *mcp.ServerSession, req *mcp.CallToolParamsFor[UpsertStrategyParams]) (*mcp.CallToolResultFor[any], error) {
	result := &mcp.CallToolResultFor[any]{}
	_, err := os.Stat(path.Join(freqtrade.Dir, req.Arguments.UserDir))
	if err != nil {
		result.Content = append(result.Content, &mcp.TextContent{Text: err.Error() + "userdir not exist. You Should Use create-user-dir tool"})
		return result, nil
	}
	err = os.WriteFile(path.Join(freqtrade.Dir, req.Arguments.UserDir, "strategies", req.Arguments.FileName), []byte(req.Arguments.Strategy), 0644)
	if err != nil {
		result.Content = append(result.Content, &mcp.TextContent{Text: err.Error()})
		return result, nil
	}
	result.Content = append(result.Content, &mcp.TextContent{Text: "upload success"})
	return result, nil
}

type UpsertConfigParams struct {
	Config  string `json:"config" jsonschema:"config.json content"`
	UserDir string `json:"userdir" jsonschema:"Path to userdata directory"`
}

func (p *UpsertConfigParams) Param() []string {
	return utils.StructJsonParams(p)
}

func UpsertConfig(ctx context.Context, session *mcp.ServerSession, req *mcp.CallToolParamsFor[UpsertConfigParams]) (*mcp.CallToolResultFor[any], error) {
	result := &mcp.CallToolResultFor[any]{}
	_, err := os.Stat(path.Join(freqtrade.Dir, req.Arguments.UserDir))
	if err != nil {
		result.Content = append(result.Content, &mcp.TextContent{Text: err.Error() + "userdir not exist"})
		return result, nil
	}
	err = os.WriteFile(path.Join(freqtrade.Dir, req.Arguments.UserDir, "config.json"), []byte(req.Arguments.Config), 0644)
	if err != nil {
		result.Content = append(result.Content, &mcp.TextContent{Text: err.Error()})
		return result, nil
	}
	result.Content = append(result.Content, &mcp.TextContent{Text: "upload success"})
	return result, nil
}

type GetBacktestingResultParams struct {
	UserDir string `json:"userdir" jsonschema:"Path to userdata directory"`
}

func (p *GetBacktestingResultParams) Param() []string {
	return utils.StructJsonParams(p)
}

func GetBacktestingTrades(ctx context.Context, session *mcp.ServerSession, req *mcp.CallToolParamsFor[GetBacktestingResultParams]) (*mcp.CallToolResultFor[any], error) {
	result := &mcp.CallToolResultFor[any]{}
	content, err := latestBacktestingJson(req.Arguments.UserDir)
	if err != nil {
		result.Content = append(result.Content, &mcp.TextContent{Text: err.Error()})
		return result, nil
	}
	formattedContent := formatBacktestingJson(content)
	result.Content = append(result.Content, &mcp.TextContent{Text: formattedContent})
	result.Content = append(result.Content, &mcp.TextContent{
		Text: "This is a csv file. you should save it to your system without changing anything.",
	})
	return result, nil
}

func formatBacktestingJson(content string) string {
	var data map[string]interface{}
	if err := jsoniter.UnmarshalFromString(content, &data); err != nil {
		return content
	}
	strategy, ok := data["strategy"].(map[string]interface{})
	if !ok {
		return content
	}
	var resultBuilder strings.Builder
	for strategyName, strategyData := range strategy {
		strategyInfo, ok := strategyData.(map[string]interface{})
		if !ok {
			continue
		}
		trades, ok := strategyInfo["trades"].([]interface{})
		if !ok || len(trades) == 0 {
			continue
		}
		csv := convertTradesToColumnCSV(trades)
		if resultBuilder.Len() > 0 {
			resultBuilder.WriteString("\n\n")
		}
		resultBuilder.WriteString(fmt.Sprintf("Strategy: %s\n\n%s", strategyName, csv))
	}
	if resultBuilder.Len() > 0 {
		return resultBuilder.String()
	}
	return content
}

func convertTradesToColumnCSV(trades []interface{}) string {
	if len(trades) == 0 {
		return ""
	}
	result := make(map[string][]interface{})
	for i, trade := range trades {
		tradeMap, ok := trade.(map[string]interface{})
		if !ok {
			continue
		}
		for key, value := range tradeMap {
			if result[key] == nil {
				result[key] = make([]interface{}, len(trades))
			}
			result[key][i] = value
		}
	}
	var csvBuilder strings.Builder
	keys := []string{
		"open_date",
		"close_date",
		"open_timestamp",
		"close_timestamp",
		"trade_duration",
		"pair",
		"amount",
		"stake_amount",
		"max_stake_amount",
		"open_rate",
		"close_rate",
		"min_rate",
		"max_rate",
		"stop_loss_abs",
		"stop_loss_ratio",
		"initial_stop_loss_abs",
		"initial_stop_loss_ratio",
		"leverage",
		"fee_open",
		"fee_close",
		"funding_fees",
		"profit_ratio",
		"profit_abs",
		"exit_reason",
		"is_open",
		"is_short",
		"enter_tag",
		"orders",
	}
	for _, key := range keys {
		values, ok := result[key]
		if !ok {
			continue
		}
		csvBuilder.WriteString(key)
		for _, value := range values {
			csvBuilder.WriteString(",")
			csvBuilder.WriteString(fmt.Sprintf("%v", value))
		}
		csvBuilder.WriteString("\n")
	}
	return csvBuilder.String()
}

func latestBacktestingFilePath(userDir string) (string, error) {
	recordPath := filepath.Join(freqtrade.Dir, userDir, "backtest_results", ".last_result.json")
	content, err := os.ReadFile(recordPath)
	if err != nil {
		return "", err
	}
	result := filepath.Join(freqtrade.Dir, userDir, "backtest_results", jsoniter.Get(content, "latest_backtest").ToString())
	return result, nil
}

func latestBacktestingJson(userDir string) (string, error) {
	backtestingFilePath, err := latestBacktestingFilePath(userDir)
	if err != nil {
		return "", err
	}
	jsonFilename := strings.TrimSuffix(filepath.Base(backtestingFilePath), filepath.Ext(backtestingFilePath)) + ".json"
	if filepath.Ext(backtestingFilePath) == ".zip" {
		f, err := os.ReadFile(backtestingFilePath)
		if err != nil {
			return "", err
		}
		zipReader, err := zip.NewReader(bytes.NewReader(f), int64(len(f)))
		if err != nil {
			return "", err
		}
		for _, file := range zipReader.File {
			if file.Name == jsonFilename {
				jsonFile, err := file.Open()
				if err != nil {
					return "", err
				}
				defer jsonFile.Close()
				content, err := io.ReadAll(jsonFile)
				if err != nil {
					return "", err
				}
				return string(content), nil
			}
		}
	}
	return "", fmt.Errorf("latestBacktestingJson error")
}
