package main

import (
	"flag"
	"freqtrade_mcp/freqtrade"
	"freqtrade_mcp/mcp/tool"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"net/http"
)

func main() {
	flag.StringVar(&freqtrade.Dir, "dir", "", "freqtrade directory")
	flag.Parse()
	if freqtrade.Dir == "" {
		panic("Please set freqtradeDir variable")
	}

	server := mcp.NewServer(&mcp.Implementation{Name: "freqtrade_mcp", Version: "v1.0.0"}, nil)
	mcp.AddTool(server, &mcp.Tool{Name: "backtesting", Description: ""}, tool.BackTest)
	handler := mcp.NewSSEHandler(func(request *http.Request) *mcp.Server {
		return server
	})
	err := http.ListenAndServe("localhost:8080", handler)
	if err != nil {
		panic(err)
	}
}
