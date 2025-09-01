package tool

import (
	"context"
	"freqtrade_mcp/freqtrade"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func CreateUserDir(ctx context.Context, session *mcp.ServerSession, req *mcp.CallToolParamsFor[freqtrade.CreateUserDirParams]) (*mcp.CallToolResultFor[any], error) {
	output, err := freqtrade.CreateUserDir(req.Arguments)
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
