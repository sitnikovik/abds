-- Локальная таблица оплат счетов на кластере.
CREATE TABLE IF NOT EXISTS energosbyt.fact_payments ON CLUSTER energosbyt (
    realty_id UInt64, -- Идентификатор объекта недвижимости
    user_id UInt64, -- Идентификатор пользователя
    price UInt32, -- Уплаченная сумма в копейках
    created_at DateTime -- Дата создания записи
) ENGINE = ReplicatedMergeTree('/clickhouse/tables/{shard}/fact_payments', '{replica}')
PARTITION BY (toYYYYMM(created_at))
ORDER BY (created_at, realty_id);
