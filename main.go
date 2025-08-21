package main

import (
	"flag"
	"freqtrade_mcp/freqtrade"
	"freqtrade_mcp/mcp/tool"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var address string

func main() {
	flag.StringVar(&freqtrade.Dir, "dir", "", "freqtrade directory")
	flag.StringVar(&address, "address", "localhost:8000", "address to listen on")
	flag.Parse()
	if freqtrade.Dir == "" {
		panic("Please set freqtradeDir variable")
	}

	server := mcp.NewServer(&mcp.Implementation{Name: "freqtrade_mcp", Version: "v1.0.0"}, nil)
	mcp.AddTool(server, &mcp.Tool{Name: "backtesting", Description: "backtesting strategy"}, tool.BackTest)
	mcp.AddTool(server, &mcp.Tool{Name: "download-data", Description: "download data from exchange"}, tool.DownloadData)
	handler := mcp.NewSSEHandler(func(request *http.Request) *mcp.Server {
		return server
	})
	err := http.ListenAndServe(address, handler)
	if err != nil {
		panic(err)
	}
}
