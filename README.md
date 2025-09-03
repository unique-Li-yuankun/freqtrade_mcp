# FreqTrade MCP

### 概述
FreqTrade MCP 是一个集成了模型上下文协议（MCP）的加密货币回测工具。

### 快速开始
写一个追涨杀跌策略，用户目录lyk，帮我回测并获取交易记录。

### 工具

| 工具名称 | 描述 | 功能说明 | 参数说明                                                                                                                                                                                                               |
|---------|------|---------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **create-user-dir** | 创建用户目录用于存放策略代码与回测结果 | 创建新的用户数据目录，包含策略、配置和回测结果的文件夹结构 | `userdir`: 用户数据目录的相对路径 (例如: lyk)                                                                                                                                                                                   |
| **download-data** | 从交易所下载OHLCV历史数据 | 下载指定交易对和时间范围的历史数据，用于回测 | `exchange`: 交易所名称<br>`timeframe`: 时间帧<br>`pairs`: 交易对列表<br>`timerange`: 时间范围<br>`userdir`: 用户数据目录路径                                                                                                                |
| **backtesting** | 使用历史数据测试策略 | 对指定策略进行历史数据回测，生成交易结果和统计信息 | `timeframe`: 时间帧 (1m,5m,30m,1h,1d)<br>`timerange`: 数据时间范围<br>`max-open-trades`: 最大开仓数量<br>`stake-amount`: 每笔投资金额<br>`pairs`: 交易对列表<br>`starting-balance`: 起始余额<br>`strategy-list`: 要回测的策略列表<br>`userdir`: 用户数据目录路径 |
| **backtesting-analysis** | 分析回测结果 | 对回测结果进行深入分析，生成详细的统计报告 | `userdir`: 用户数据目录路径                                                                                                                                                                                                |
| **upsert-strategy** | 上传/更新策略代码 | 上传Python策略文件到指定用户目录的strategies文件夹 | `filename`: 策略文件名<br>`strategy`: 策略代码内容<br>`userdir`: 用户数据目录路径                                                                                                                                                     |
| **upsert-config** | 上传/更新配置文件 | 上传配置文件(config.json)到指定用户目录 | `config`: config.json文件内容<br>`userdir`: 用户数据目录路径                                                                                                                                                                   |
| **get-backtesting-trades** | 获取回测交易数据 | 获取最新回测结果中的交易数据，以CSV格式返回 | `userdir`: 用户数据目录路径                                                                                                                                                                                                |

### 使用
PS ： 预先下载可执行文件或自行编译（https://github.com/unique-Li-yuankun/freqtrade_mcp/releases)

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
        "url": "http://127.0.0.1:8000"
      }
    }
}
```

- **Claude Code**
```bash
添加到当前项目：
claude mcp add --transport http freqtrade_mcp http://127.0.0.1:8000

添加到全局：
claude mcp add --scope user --transport http freqtrade_mcp http://127.0.0.1:8000
```
