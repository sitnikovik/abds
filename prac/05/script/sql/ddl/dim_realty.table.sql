-- Локальная таблица материальных объектов на кластере.
CREATE TABLE IF NOT EXISTS energosbyt.dim_realty ON CLUSTER energosbyt (
    id UInt64, -- Snowflake ID
    address String, -- Адрес полностью
    square UInt32, -- Площадь (кв. м.)
    type LowCardinality(String), -- тип недвижимости
    number String, -- Номер квартиры (если есть)
    floor UInt8, -- Этаж
    porch UInt8, -- Подъезд
    created_at DateTime -- Дата создания записи
) ENGINE = ReplicatedMergeTree('/clickhouse/tables/{shard}/dim_realty', '{replica}')
ORDER BY (address, id)
SETTINGS enable_block_number_column = 1, enable_block_offset_column = 1;
