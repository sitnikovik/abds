-- Локальная таблица установок счётчиков на кластере (аннулирование через sign).
CREATE TABLE IF NOT EXISTS energosbyt.fact_gauge_bindings ON CLUSTER energosbyt (
    gauge_id UInt64, -- Идентификатор счётчика
    realty_id UInt64, -- Идентификатор объекта недвижимости
    created_at DateTime, -- Момент установки счётчика на объект
    sign Int8 -- +1 установка, -1 аннулирование
) ENGINE = ReplicatedMergeTree('/clickhouse/tables/{shard}/fact_gauge_bindings', '{replica}')
ORDER BY (gauge_id, created_at);
