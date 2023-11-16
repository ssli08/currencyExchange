package main

import (
	"currencyExchange/currency"
	"fmt"
	"os"

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
}

func main() {
	cfg := parseConfig()

	msg, err := currency.GetCurrencyRates(cfg.CurrencyAPI.ApiURL, cfg.CurrencyAPI.ApiKEY)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// im.SendMSG(msg)
	if err := currency.SendMSGViaProxy(cfg.TelegramBot.TelegramURL, cfg.TelegramBot.BotToken, cfg.TelegramBot.ChatGroupID, msg); err != nil {
		fmt.Println(err)
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
