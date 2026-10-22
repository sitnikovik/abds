-- Загрузка data/dim_gauges.csv в dim_gauges. Данные приходят через stdin (цель load-big в Makefile).
INSERT INTO energosbyt.dim_gauges (id, name, country, created_at)
SELECT id, name, country, created_at
FROM input('id UInt64, name String, country String, created_at DateTime')
FORMAT CSVWithNames
