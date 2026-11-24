-- Материализованное представление для автоматической вставки данных из Kafka.
-- Таблица Kafka создана только на ch1, поэтому MV тоже только на ch1 (не ON CLUSTER).
-- MV пишет в локальную таблицу fact_meter_readings на ch1.
-- fact_meter_readings использует ReplicatedMergeTree, который автоматически
-- реплицирует данные на ch2 (внутри шарда 1).
-- Для чтения со всех шардов используйте Distributed-таблицу fact_meter_readings_dist.
CREATE MATERIALIZED VIEW IF NOT EXISTS energosbyt.mv_kafka_to_meter_readings
TO energosbyt.fact_meter_readings -- Целевая таблица: локальная (реплицируется внутри шарда).
AS
SELECT
    realty_id,
    user_id,
    tariff_id,
    gauge_id,
    toStartOfMonth(created_at) AS reading_month, -- Агрегация показаний по месяцам.
    day_kw,
    night_kw,
    created_at
FROM energosbyt.kafka_gauge_readings; -- Источник: локальная Kafka-таблица на ch1.
