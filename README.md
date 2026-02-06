# FreqTrade MCP

### Overview
FreqTrade MCP is a cryptocurrency backtesting tool integrated with the Model Context Protocol (MCP).

### Quick Start
Write a momentum strategy for user directory `lyk`, help me backtest it and get trading records.

### Tools

| Tool Name | Description | Functionality | Parameters |
|-----------|-------------|---------------|-----------|
| **create-user-dir** | Create user directory for storing strategy code and backtest results | Create a new user data directory with folder structure for strategies, configurations, and backtest results | `userdir`: Relative path of user data directory (e.g., lyk) |
| **download-data** | Download OHLCV historical data from exchange | Download historical data for specified trading pairs and time range for backtesting | `exchange`: Exchange name<br>`timeframe`: Timeframe<br>`pairs`: List of trading pairs<br>`timerange`: Time range<br>`userdir`: User data directory path |
| **backtesting** | Test strategy using historical data | Backtest specified strategy with historical data, generate trading results and statistics | `timeframe`: Timeframe (1m,5m,30m,1h,1d)<br>`timerange`: Data time range<br>`max-open-trades`: Maximum open trades<br>`stake-amount`: Stake per trade<br>`pairs`: List of trading pairs<br>`starting-balance`: Starting balance<br>`strategy-list`: List of strategies to backtest<br>`userdir`: User data directory path |
| **backtesting-analysis** | Analyze backtest results | Perform in-depth analysis of backtest results, generate detailed statistical reports | `userdir`: User data directory path |
| **upsert-strategy** | Upload/update strategy code | Upload Python strategy file to strategies folder in user directory | `filename`: Strategy filename<br>`strategy`: Strategy code content<br>`userdir`: User data directory path |
| **upsert-config** | Upload/update configuration file | Upload configuration file (config.json) to user directory | `config`: config.json file content<br>`userdir`: User data directory path |
| **get-backtesting-trades** | Get backtest trading data | Get trading data from latest backtest results, return in CSV format | `userdir`: User data directory path |

### Usage
PS: Pre-download executable or compile yourself (https://github.com/unique-Li-yuankun/freqtrade_mcp/releases)

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

Configure MCP Server

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
Add to current project:
claude mcp add --transport http freqtrade_mcp http://127.0.0.1:8000

Add globally:
claude mcp add --scope user --transport http freqtrade_mcp http://127.0.0.1:8000
```
