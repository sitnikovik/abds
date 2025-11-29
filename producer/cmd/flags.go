package main

import "flag"

const (
	// configFlagName имя флага пути к файлу конфига.
	configFlagName = "config"
	// configFlagDescription опиcание для флага пути к файлу конфига.
	configFlagDescription = "Path to config file"
	// configFlagDefaultValue дефолтный путь к файлу конфига.
	configFlagDefaultValue = "./../configs/producer.yaml"
	// pathToFlatsCSVFlagName имя флага пути к файлу CSV с квартирами.
	pathToFlatsCSVFlagName = "flats"
	// pathToFlatsCSVDescription опиcание для флага пути к файлу CSV с квартирами.
	pathToFlatsCSVDescription = "Path to flats CSV file"
	// pathToFlatsCSVDefaultValue дефолтный путь к файлу CSV с квартирами.
	pathToFlatsCSVDefaultValue = "./flats.energosbyt.csv"
)

// flags содержит флаги и параметры команды.
type flags struct {
	// configPath путь к файлу конфига.
	configPath string
	// flatsCSV путь к CSV с квартирами.
	flatsCSV string
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
	flag.StringVar(
		&f.flatsCSV,
		pathToFlatsCSVFlagName,
		pathToFlatsCSVDefaultValue,
		pathToFlatsCSVDescription,
	)
	flag.Parse()
	return f
}
