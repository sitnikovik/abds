CREATE TABLE IF NOT EXISTS energosbyt.fact_payments (
    realty_id UInt64, -- Идентификатор объекта недвижимости
    user_id UInt64, -- Идентификатор пользователя
    price UInt32, -- Уплаченная сумма в копейках
    created_at DateTime -- Дата создания записи
) ENGINE = MergeTree
PARTITION BY (toYYYYMM(created_at))
ORDER BY (created_at, realty_id);