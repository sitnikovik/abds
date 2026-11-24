-- Загрузка data/dim_realty.dataset.csv в dim_realty_dist (Distributed обёртка). Данные приходят через stdin (цель load-big в Makefile).

-- insert_distributed_sync: фоновая отправка блоков через DistributedInsertQueue дублирует блоки при пересогласованиях. Синхронно — ровно одна копия на шард.
INSERT INTO energosbyt.dim_realty_dist (id, address, square, type, number, floor, porch, created_at) SETTINGS insert_distributed_sync = 1
SELECT id, address, square, type, number, floor, porch, created_at
FROM input('id UInt64, address String, square UInt32, type String, number String, floor UInt8, porch UInt8, created_at DateTime')
FORMAT CSVWithNames
