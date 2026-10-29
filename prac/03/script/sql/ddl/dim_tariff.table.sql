CREATE TABLE IF NOT EXISTS energosbyt.dim_tariffs (
    id UInt64, -- Showflake ID
    name String, -- Название тарифа
    price_day_kwph UInt32, -- Стоимость одного кВт/ч в дневное время, в копейках (SCD 2: версия)
    price_night_kwph UInt32, -- Стоимость одного кВт/ч в ночное время, в копейках (SCD 2: версия)
    valid_from DateTime, -- Дата начала действия версии
    valid_to DateTime, -- Дата окончания действия версии (2099-01-01 00:00:00 у открытой версии)
    created_at DateTime, -- Дата создания записи
    updated_at DateTime -- Дата обновления записи (версия строки)
) ENGINE = ReplacingMergeTree(updated_at)
ORDER BY (id, valid_from);
