package main

import (
	"flag"
	"fmt"
	"freqtrade_mcp/freqtrade"
	"freqtrade_mcp/mcp/tool"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var address string

func main() {
	flag.StringVar(&freqtrade.Dir, "dir", "", "freqtrade directory")
	flag.StringVar(&address, "addr", "localhost:8000", "address to listen on")
	flag.Parse()
	if freqtrade.Dir == "" {
		panic("Please set dir arg.")
	}

	server := mcp.NewServer(&mcp.Implementation{Name: "freqtrade_mcp", Version: "v1.0.0"}, nil)
	mcp.AddTool(server, &mcp.Tool{Name: "create-user-dir", Description: "create user-data directory"}, tool.CreateUserDir)
	mcp.AddTool(server, &mcp.Tool{Name: "backtesting", Description: "backtesting strategy"}, tool.Backtesting)
	mcp.AddTool(server, &mcp.Tool{Name: "backtesting-analysis", Description: "backtesting result analysis"}, tool.BacktestingAnalysis)
	mcp.AddTool(server, &mcp.Tool{Name: "download-data", Description: "download data from exchange"}, tool.DownloadData)

	mcp.AddTool(server, &mcp.Tool{Name: "get-backtesting-trades", Description: "get backtesting trades"}, tool.GetBacktestingResult)
	mcp.AddTool(server, &mcp.Tool{Name: "upsert-config", Description: "upsert config. Must use it to add or update config.json"}, tool.UpsertConfig)
	mcp.AddTool(server, &mcp.Tool{Name: "upsert-strategy", Description: "upsert strategy. Must use it to add or update a strategy"}, tool.UpsertStrategy)

	handler := mcp.NewStreamableHTTPHandler(func(request *http.Request) *mcp.Server {
		return server
	}, &mcp.StreamableHTTPOptions{})
	go func() {
		err := http.ListenAndServe(address, loggingMiddleware(handler))
		if err != nil {
			panic(err)
		}
	}()
	fmt.Println("freqtrade_mcp is running on", address)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	fmt.Println("freqtrade_mcp is shutting down")
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		log.Printf("%s %s %d", r.Method, r.URL.Path, duration)
	})
}
