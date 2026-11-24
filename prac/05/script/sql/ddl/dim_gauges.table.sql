-- Локальная таблица счётчиков на кластере.
CREATE TABLE IF NOT EXISTS energosbyt.dim_gauges ON CLUSTER energosbyt (
    id UInt64, -- Snowflake ID
    name String, -- Название производителя и модели
    country LowCardinality(String), -- Страна производства
    created_at DateTime -- Дата создания записи
) ENGINE = ReplicatedMergeTree('/clickhouse/tables/{shard}/dim_gauges', '{replica}')
ORDER BY (name, id);
