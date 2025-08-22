package tool

import (
	"context"
	"freqtrade_mcp/freqtrade"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func DownloadData(ctx context.Context, session *mcp.ServerSession, req *mcp.CallToolParamsFor[freqtrade.DownloadDataParams]) (*mcp.CallToolResultFor[any], error) {
	output, err := freqtrade.DownloadData(req.Arguments)

	result := &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: output},
			&mcp.TextContent{Text: "if data length equal zero, maybe data already exists."},
		},
	}
	if err != nil {
		result.Content = append(result.Content, &mcp.TextContent{Text: err.Error()})
	}
	return result, nil
}
