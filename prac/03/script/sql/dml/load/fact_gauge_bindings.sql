-- Загрузка data/fact_gauge_bindings.csv в fact_gauge_bindings. Данные приходят через stdin (цель load-big в Makefile).
INSERT INTO energosbyt.fact_gauge_bindings (gauge_id, realty_id, created_at, sign)
SELECT gauge_id, realty_id, created_at, sign
FROM input('gauge_id UInt64, realty_id UInt64, created_at DateTime, sign Int8')
FORMAT CSVWithNames
