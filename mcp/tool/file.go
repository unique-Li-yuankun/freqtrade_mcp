package tool

import (
	"context"
	"freqtrade_mcp/freqtrade"
	"freqtrade_mcp/utils"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"os"
	"path"
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
