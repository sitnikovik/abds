CREATE TABLE IF NOT EXISTS energosbyt.dim_gauges (
    id UInt64, -- Showflake ID
    name String, -- Название производителя и модели
    country LowCardinality(String), -- Страна производства
    created_at DateTime -- Дата создания записи
) ENGINE = MergeTree()
ORDER BY (name, id);