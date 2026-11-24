-- Загрузка data/dim_tariffs.dataset.csv в dim_tariffs_dist (Distributed обёртка). Данные приходят через stdin (цель load-big в Makefile).

-- insert_distributed_sync: фоновая отправка блоков через DistributedInsertQueue дублирует блоки при пересогласованиях. Синхронно — ровно одна копия на шард.
INSERT INTO energosbyt.dim_tariffs_dist (id, name, price_day_kwph, price_night_kwph, valid_from, valid_to, created_at, updated_at) SETTINGS insert_distributed_sync = 1
SELECT id, name, price_day_kwph, price_night_kwph, valid_from, valid_to, created_at, updated_at
FROM input('id UInt64, name String, price_day_kwph UInt32, price_night_kwph UInt32, valid_from DateTime, valid_to DateTime, created_at DateTime, updated_at DateTime')
FORMAT CSVWithNames
