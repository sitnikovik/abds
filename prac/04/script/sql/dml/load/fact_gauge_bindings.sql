-- Загрузка data/fact_gauge_bindings.dataset.csv в fact_gauge_bindings_dist (Distributed обёртка). Данные приходят через stdin (цель load-big в Makefile).

-- insert_distributed_sync: фоновая отправка блоков через DistributedInsertQueue дублирует блоки при пересогласованиях. Синхронно — ровно одна копия на шард.
INSERT INTO energosbyt.fact_gauge_bindings_dist (gauge_id, realty_id, created_at, sign) SETTINGS insert_distributed_sync = 1
SELECT gauge_id, realty_id, created_at, sign
FROM input('gauge_id UInt64, realty_id UInt64, created_at DateTime, sign Int8')
FORMAT CSVWithNames
