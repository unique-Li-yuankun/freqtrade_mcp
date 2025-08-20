package tool

import (
	"context"
	"freqtrade_mcp/freqtrade"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func BackTest(ctx context.Context, session *mcp.ServerSession, req *mcp.CallToolParamsFor[freqtrade.BackTestParams]) (*mcp.CallToolResultFor[any], error) {
	output, err := freqtrade.BackTest(req.Arguments)
	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{&mcp.TextContent{Text: output}},
	}, err
}
