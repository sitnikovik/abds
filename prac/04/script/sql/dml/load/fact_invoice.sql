-- Загрузка data/fact_invoice.dataset.csv в fact_invoice. Данные приходят через stdin (цель load-big в Makefile).

-- insert_distributed_sync: фоновая отправка блоков через DistributedInsertQueue дублирует блоки при пересогласованиях. Синхронно — ровно одна копия на шард.
INSERT INTO energosbyt.fact_invoice_dist (realty_id, user_id, tariff_id, gauge_id, day_kw, night_kw, price, deadline, created_at) SETTINGS insert_distributed_sync = 1
SELECT realty_id, user_id, tariff_id, gauge_id, day_kw, night_kw, price, deadline, created_at
FROM input('realty_id UInt64, user_id UInt64, tariff_id UInt64, gauge_id UInt64, day_kw UInt32, night_kw UInt32, price UInt32, deadline Date, created_at DateTime')
FORMAT CSVWithNames
