package tool

import (
	"context"
	"freqtrade_mcp/freqtrade"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type ExchangeDataFilesParams struct {
}

func ExchangeDataFiles(ctx context.Context, session *mcp.ServerSession, req *mcp.CallToolParamsFor[ExchangeDataFilesParams]) (*mcp.CallToolResultFor[any], error) {
	dir := path.Join(freqtrade.Dir, "user_data", "data")
	files, err := filesUnderDir(dir)
	result := &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: strings.Join(files, "\n")},
		},
	}
	if err != nil {
		result.Content = append(result.Content, &mcp.TextContent{Text: err.Error()})
	}
	return result, nil
}

type ReadExchangeDataFileParams struct {
	FilePath string `json:"file_path" jsonschema:"Exchange data file absolute path."`
}

func ReadExchangeDataFile(ctx context.Context, session *mcp.ServerSession, req *mcp.CallToolParamsFor[ReadExchangeDataFileParams]) (*mcp.CallToolResultFor[any], error) {
	content, err := readFile(req.Arguments.FilePath)
	result := &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: content},
		},
	}
	if err != nil {
		result.Content = append(result.Content, &mcp.TextContent{Text: err.Error()})
	}
	return result, nil
}

func filesUnderDir(dir string) ([]string, error) {
	var fileList []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			fileList = append(fileList, path)
		}
		return nil
	})
	return fileList, err
}

func readFile(path string) (string, error) {
	b, err := os.ReadFile(path)
	return string(b), err
}
