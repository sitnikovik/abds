CREATE TABLE IF NOT EXISTS energosbyt.fact_gauge_bindings (
    gauge_id UInt64, -- Идентификатор счётчика
    realty_id UInt64, -- Идентификатор объекта недвижимости
    created_at DateTime, -- Момент установки счётчика на объект
    sign Int8 -- +1 установка, -1 аннулирование
) ENGINE = CollapsingMergeTree(sign)
ORDER BY (gauge_id, created_at);
