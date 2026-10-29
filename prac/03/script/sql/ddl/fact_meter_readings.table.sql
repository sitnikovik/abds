CREATE TABLE IF NOT EXISTS energosbyt.fact_meter_readings (
    realty_id UInt64, -- Идентификатор объекта недвижимости
    user_id UInt64, -- Идентификатор пользователя
    tariff_id UInt64, -- Идентификатор тарифа, по которому производится расчёт
    gauge_id UInt64, -- Идентификатор счётчика
    reading_month Date, -- Месяц, за который передано показание
    day_kw UInt32, -- Показания электросчётчика за дневное время
    night_kw UInt32, -- Показания электросчётчика за ночное время
    created_at DateTime -- Дата создания записи (побеждает самая свежая)
) ENGINE = ReplacingMergeTree(created_at)
PARTITION BY (toYYYYMM(reading_month))
ORDER BY (realty_id, gauge_id, reading_month);