package main

import "flag"

const (
	// configFlagName имя флага пути к файлу конфига.
	configFlagName = "config"
	// configFlagDescription опиcание для флага пути к файлу конфига.
	configFlagDescription = "Path to config file"
	// configFlagDefaultValue дефолтный путь к файлу конфига.
	configFlagDefaultValue = "./../configs/producer.yaml"
)

// flags содержит флаги и параметры команды.
type flags struct {
	// configPath путь к файлу конфига.
	configPath string
}

// parseFlags парсит флаги командной строки и возвращает их.
func parseFlags() flags {
	var f flags
	flag.StringVar(
		&f.configPath,
		configFlagName,
		configFlagDefaultValue,
		configFlagDescription,
	)
	flag.Parse()
	return f
}
