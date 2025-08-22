# FreqTrade MCP

### 概述
FreqTrade MCP 是一个集成了模型上下文协议（MCP）的加密货币回测工具。

### 工具
- **数据下载**: 获取交易所OHLCV数据
- **回测功能**: 使用历史数据测试策略

### 使用
PS ： 预先下载可执行文件（https://github.com/unique-Li-yuankun/freqtrade_mcp/releases)

- **Windows**:
```bash 
git clone https://github.com/freqtrade/freqtrade.git
Set-ExecutionPolicy -ExecutionPolicy Bypass
cd freqtrade
.\setup.ps1
.\freqtrade_mcp -dir W:\freqtrade -addr localhost:8000
```

- **Linux/MacOS**
```bash 
git clone https://github.com/freqtrade/freqtrade.git
cd freqtrade
./setup.sh -i
./freqtrade_mcp -dir path/to/freqtrade -addr localhost:8000
```

配置MCP服务器

- **Cursor**
```json
{
    "mcpServers": {
      "freqtrade_mcp": {
        "url": "http://127.0.0.1:8000",
      }
    }
}
```