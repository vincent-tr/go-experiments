# 🧠 Trading Bot Architecture: Core Components

A modular trading bot architecture helps keep logic clean and maintainable. Below are the three main components typically involved in trade execution logic:

## 📊 1. Strategy

**🎯 Purpose**  
The **Strategy** component is responsible for analyzing market data and deciding whether to **open a position**, and in which **direction** (long or short).

**🔍 Responsibilities**
- Evaluate indicators (e.g., EMA crossover, RSI)
- Use historical price patterns or statistical models
- Confirm market conditions (volatility, sessions, etc.)
- Output a trade signal

**📥 Input**
- Historical candles or ticks

**📤 Output**
- Signal: `Buy`, `Sell`, or `None`

## 🛡️ 2. RiskManager

**🎯 Purpose**  
The **RiskManager** is responsible for defining protective risk parameters for every trade, specifically the **stop-loss** and **take-profit** levels.

**🔍 Responsibilities**
- Compute **Stop-Loss (SL)** using:
  - ATR (Average True Range)
  - Recent swing highs/lows
  - Fixed pip offsets (e.g. 3 pips)
- Compute **Take-Profit (TP)** using:
  - Reward-to-Risk ratios (e.g. 2:1)
- Optionally account for:
  - Market volatility
  - Spread or slippage
  - Dynamic adjustment post-entry

**📥 Input**
- Entry price
- Trade direction (long or short)
- Market context:
  - Volatility indicators (e.g. ATR)
  - Historical prices (for swing highs/lows)
  - Instrument-specific spread

**📤 Output**
- `StopLoss` price
- `TakeProfit` price

## 💰 3. CapitalAllocator

**🎯 Purpose**  
The **CapitalAllocator** determines how much capital should be allocated to each trade by calculating the optimal **position size** based on risk parameters and account balance.

**📈 Responsibilities**
- Compute the **number of lots or units** to trade
- Ensure that position size aligns with:
  - Risk tolerance (e.g., risking 1% of capital per trade)
  - Leverage and margin requirements
  - Available account balance
- Prevent overexposure by applying caps based on available capital

**📥 Input**
- Account balance or equity
- Capital risk percentage (e.g., 1%)
- Stop-loss distance (in pips or price)
- Instrument pip value per lot or tick value
- Leverage (if applicable)

**📤 Output**
- `PositionSize` (e.g., in lots or units)
