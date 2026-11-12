-- Загрузка data/fact_payments.dataset.csv в fact_payments_dist (Distributed обёртка). Данные приходят через stdin (цель load-big в Makefile).

-- insert_distributed_sync: фоновая отправка блоков через DistributedInsertQueue дублирует блоки при пересогласованиях. Синхронно — ровно одна копия на шард.
INSERT INTO energosbyt.fact_payments_dist (realty_id, user_id, price, created_at) SETTINGS insert_distributed_sync = 1
SELECT realty_id, user_id, price, created_at
FROM input('realty_id UInt64, user_id UInt64, price UInt32, created_at DateTime')
FORMAT CSVWithNames
