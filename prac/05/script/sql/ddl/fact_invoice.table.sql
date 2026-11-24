-- Локальная таблица счетов на кластере.
CREATE TABLE IF NOT EXISTS energosbyt.fact_invoice ON CLUSTER energosbyt (
    realty_id UInt64, -- Идентификатор объекта недвижимости
    user_id UInt64, -- Идентификатор пользователя
    tariff_id UInt64, -- Идентификатор тарифа
    gauge_id UInt64, -- Идентификатор счётчика
    day_kw UInt32, -- Расход за дневное время, кВт/ч
    night_kw UInt32, -- Расход за ночное время, кВт/ч
    price UInt32, -- Сумма к оплате в копейках
    deadline Date, -- Срок оплаты: до 10-го числа следующего месяца
    created_at DateTime -- Дата выставления счёта
) ENGINE = ReplicatedMergeTree('/clickhouse/tables/{shard}/fact_invoice', '{replica}')
PARTITION BY (toYYYYMM(created_at))
ORDER BY (created_at, realty_id);
