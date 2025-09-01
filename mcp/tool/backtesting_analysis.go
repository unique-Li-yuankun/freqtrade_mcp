package tool

import (
	"context"
	"freqtrade_mcp/freqtrade"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func BacktestingAnalysis(ctx context.Context, session *mcp.ServerSession, req *mcp.CallToolParamsFor[freqtrade.BacktestingAnalysisParams]) (*mcp.CallToolResultFor[any], error) {
	output, err := freqtrade.BacktestingAnalysis(req.Arguments)
	result := &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: output},
		},
	}
	if err != nil {
		result.Content = append(result.Content, &mcp.TextContent{Text: err.Error()})
	}
	return result, nil
}
