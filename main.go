package main

import (
	"currencyExchange/currency"
	"currencyExchange/im"
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

	msg := currency.GetCurrencyRates(cfg.CurrencyAPI.ApiURL, cfg.CurrencyAPI.ApiKEY)
	// im.SendMSG(msg)
	im.SendMSGViaProxy(cfg.TelegramBot.TelegramURL, cfg.TelegramBot.BotToken, cfg.TelegramBot.ChatGroupID, msg)
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
