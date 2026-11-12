-- Загрузка data/dim_gauges.dataset.csv в dim_gauges_dist (Distributed обёртка). Данные приходят через stdin (цель load-big в Makefile).

-- insert_distributed_sync: фоновая отправка блоков через DistributedInsertQueue дублирует блоки при пересогласованиях. Синхронно — ровно одна копия на шард.
INSERT INTO energosbyt.dim_gauges_dist (id, name, country, created_at) SETTINGS insert_distributed_sync = 1
SELECT id, name, country, created_at
FROM input('id UInt64, name String, country String, created_at DateTime')
FORMAT CSVWithNames
