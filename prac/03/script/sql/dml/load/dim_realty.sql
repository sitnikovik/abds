-- Загрузка data/dim_realty.csv в dim_realty. Данные приходят через stdin (цель load-big в Makefile).
INSERT INTO energosbyt.dim_realty (id, address, square, type, number, floor, porch, created_at)
SELECT id, address, square, type, number, floor, porch, created_at
FROM input('id UInt64, address String, square UInt32, type String, number String, floor UInt8, porch UInt8, created_at DateTime')
FORMAT CSVWithNames
