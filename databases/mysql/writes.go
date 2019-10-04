package database

import (
	"time"
)

// InsertNewTrade .
func InsertNewTrade(trade model.Trade) error {

	_, err := Db.NamedExec(`INSERT INTO trades (market, exchange_1_name, exchange_1_buy_price, exchange_1_amount, exchange_2_name, 
                    exchange_2_sell_price, exchange_2_amount, profit, profit_currency, when_unix)
	     VALUES (:market, :exchange_1_name, :exchange_1_buy_price, :exchange_1_amount, :exchange_2_name, 
                    :exchange_2_sell_price, :exchange_2_amount, :profit, :profit_currency, :when_unix)`,
		map[string]interface{}{
			"market":                trade.Market,
			"exchange_1_name":       trade.Exchange1Name,
			"exchange_1_buy_price":  trade.Exchange1BuyPrice,
			"exchange_1_amount":     trade.Exchange1Amount,
			"exchange_2_name":       trade.Exchange2Name,
			"exchange_2_sell_price": trade.Exchange2SellPrice,
			"exchange_2_amount":     trade.Exchange2Amount,
			"profit":                trade.Profit,
			"profit_currency":       trade.ProfitCurrency,
			"when_unix":             time.Now().Unix(),
		})

	if err != nil {
		return err
	}

	return nil
}

// StoreNewBalance ------ example with on duplicate key update
func StoreNewBalance(balance model.Balance) error {

	_, err := Db.NamedExec(`INSERT INTO balances (exchange_name, currency, balance, updated_at_unix)
	     VALUES (:exchange_name, :currency, :balance, :updated_at_unix)  ON DUPLICATE KEY UPDATE balance=:balance;`,
		map[string]interface{}{
			"exchange_name":   balance.ExchangeName,
			"currency":        balance.Currency,
			"balance":         balance.Balance,
			"updated_at_unix": time.Now().Unix(),
		})

	if err != nil {
		return err
	}

	return nil
}
