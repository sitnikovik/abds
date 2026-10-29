CREATE TABLE IF NOT EXISTS energosbyt.dim_realty (
    id UInt64, -- Showflake ID
    address String, -- Адрес полностью
    square UInt32, -- Площаль (кв. см.)
    type LowCardinality(String), -- тип недвижимости
    number String, -- Номер квартиры (если есть)
    floor UInt8, -- Этаж
    porch UInt8, -- Подъезд
    created_at DateTime -- Дата создания записи
) ENGINE = MergeTree()
ORDER BY (address, id)
-- SETTINGS нужны для лёгкого UPDATE: служебные колонки номера блока и смещения строки, по ним патч находит строку.
SETTINGS enable_block_number_column = 1, enable_block_offset_column = 1;