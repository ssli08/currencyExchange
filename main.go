package main

import (
	"currencyExchange/currency"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	CurrencyAPI struct {
		ApiURL string `yaml:"apiURL"`
		ApiKEY string `yaml:"apiKey"`
	} `yaml:"currency"`
	TelegramBot struct {
		TelegramURL string `yaml:"tgUrl"`
		BotToken    string `yaml:"botToken"`
		ChatGroupID string `yaml:"chatGroupID"`
	} `yaml:"telegram"`
	ReferenceValue struct {
		ReferredMinValue float32 `yaml:"referredMinValue"`
		ReferredMaxValue float32 `yaml:"referredMaxValue"`
	} `yaml:"reference"`
}

func main() {
	cfg := parseConfig()

	data, err := currency.GetCurrencyRates(cfg.CurrencyAPI.ApiURL, cfg.CurrencyAPI.ApiKEY)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	c := data.Rates.CNY / data.Rates.NZD
	u := data.Rates.NZD / data.Rates.USD

	// if current rate value greater than reference max value or less than reference min value, it'll send message to
	// telegram channel
	if c <= cfg.ReferenceValue.ReferredMinValue || c >= cfg.ReferenceValue.ReferredMaxValue {
		// im.SendMSG(msg)
		msg := fmt.Sprintf("%v: (RMB/NZD: %f) - (NZD/USD: %f)", time.Unix(data.Timestamp, 0), c, u)
		if err := currency.SendMSGViaProxy(cfg.TelegramBot.TelegramURL, cfg.TelegramBot.BotToken, cfg.TelegramBot.ChatGroupID, msg); err != nil {
			fmt.Println(err)
		}
	}
}

func parseConfig() Config {
	f, err := os.Open("config.yml")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	// fmt.Println(f.Stat())
	// os.Exit(1)
	var cfg Config
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&cfg); err != nil {
		panic(err)
	}
	return cfg
}
