-- Загрузка data/dim_tariffs.csv в dim_tariffs. Данные приходят через stdin (цель load-big в Makefile).
INSERT INTO energosbyt.dim_tariffs (id, name, price_day_kwph, price_night_kwph, valid_from, valid_to, created_at, updated_at)
SELECT id, name, price_day_kwph, price_night_kwph, valid_from, valid_to, created_at, updated_at
FROM input('id UInt64, name String, price_day_kwph UInt32, price_night_kwph UInt32, valid_from DateTime, valid_to DateTime, created_at DateTime, updated_at DateTime')
FORMAT CSVWithNames
