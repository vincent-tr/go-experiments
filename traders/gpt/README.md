# 📈 EMA Crossover + RSI Filter Trading Strategy (Long/Short with SL/TP)

## Indicators
- **EMA(5)** — Fast exponential moving average  
- **EMA(20)** — Slow exponential moving average  
- **RSI(14)** — Relative Strength Index (14 periods)

## Position Entry Rules

### Long Position (Buy)
Enter a **long position** when:  
1. EMA(5) crosses **above** EMA(20) (bullish crossover)  
2. RSI(14) is **between 30 and 70**

- **Stop-Loss:** Set stop-loss **3 pips below** the recent swing low over the past **15 minutes**  
- **Take-Profit:** Set take-profit at a **2:1 reward-to-risk ratio** relative to the stop-loss distance (e.g., if stop-loss is 3 pips below entry, take-profit at 6 pips above)

### Short Position (Sell)
Enter a **short position** when:  
1. EMA(5) crosses **below** EMA(20) (bearish crossover)  
2. RSI(14) is **between 30 and 70**

- **Stop-Loss:** Set stop-loss **3 pips above** the recent swing high over the past **15 minutes**  
- **Take-Profit:** Set take-profit at a **2:1 reward-to-risk ratio** relative to the stop-loss distance (e.g., if stop-loss is 3 pips above entry, take-profit at 6 pips below)

## Position Size Calculation

To manage risk effectively, calculate your position size based on your stop-loss and the amount of capital you are willing to risk:

- Determine your **account risk** (e.g., 1% of your total capital)  
- Measure the **stop-loss distance** in pips (here, 3 pips)  
- Calculate the **pip value** per lot for your currency pair (for EUR/USD, approximately $10 per pip per standard lot)  

Use the formula:

\[
\text{Position Size (lots)} = \frac{\text{Account Risk (in currency)}}{\text{Stop-Loss (in pips)} \times \text{Pip Value per Lot}}
\]

This ensures your potential loss does not exceed your risk tolerance on any trade.

## Trading Schedule ⏰

This strategy should only be executed:
- **Between Tuesday and Thursday**, when market volume and volatility are highest  
- **During the London–New York session overlap** (typically **13:00 to 17:00 CET** or **12:00 to 16:00 GMT**, depending on DST)
- **Avoid trading during public holidays** in the **UK or US**, as liquidity is lower and market behavior is less predictable

## Parameters to Tune During Backtesting

- **EMA periods:** Try varying fast EMA (3 to 7) and slow EMA (15 to 30)  
- **RSI period and thresholds:** Adjust RSI period (10 to 20) and entry range (e.g., 25–75 instead of 30–70)  
- **Stop-loss buffer:** Test 2 to 5 pips below/above swing points  
- **Swing lookback period:** Experiment with 10 to 30 minutes to define recent swing highs/lows  
- **Take-profit ratio:** Try ratios from 1.5:1 up to 3:1 for reward-to-risk  
