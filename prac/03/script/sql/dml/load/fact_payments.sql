-- Загрузка data/fact_payments.csv в fact_payments. Данные приходят через stdin (цель load-big в Makefile).
INSERT INTO energosbyt.fact_payments (realty_id, user_id, price, created_at)
SELECT realty_id, user_id, price, created_at
FROM input('realty_id UInt64, user_id UInt64, price UInt32, created_at DateTime')
FORMAT CSVWithNames
