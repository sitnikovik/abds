-- Локальная таблица тарифов на кластере (SCD 2 версии).
CREATE TABLE IF NOT EXISTS energosbyt.dim_tariffs ON CLUSTER energosbyt (
    id UInt64, -- Snowflake ID
    name String, -- Название тарифа
    price_day_kwph UInt32, -- Цена кВт/ч днём, в копейках (SCD 2: версия)
    price_night_kwph UInt32, -- Цена кВт/ч ночью, в копейках (SCD 2: версия)
    valid_from DateTime, -- Дата начала действия версии
    valid_to DateTime, -- Дата окончания действия версии (2099 у открытой)
    created_at DateTime, -- Дата создания записи
    updated_at DateTime -- Дата обновления записи (версия строки)
) ENGINE = ReplicatedMergeTree('/clickhouse/tables/{shard}/dim_tariffs', '{replica}')
ORDER BY (id, valid_from);
