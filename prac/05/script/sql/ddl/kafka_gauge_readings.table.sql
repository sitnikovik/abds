-- Таблица-движок Kafka для чтения сообщений из топика gauges.
-- ВНИМАНИЕ: Таблица создана ТОЛЬКО на ch1 (не ON CLUSTER).
-- Это демо-архитектура: один узел читает Kafka, остальные узлы кластера не участвуют в чтении.
--
-- Для продакшена и в лабораторной работе необходимо реализовать шардированное чтение,
-- где каждый шард должен иметь свою consumer group и читать свои партиции из Kafka.
CREATE TABLE IF NOT EXISTS energosbyt.kafka_gauge_readings (
    gauge_id UInt64,
    realty_id UInt64,
    user_id UInt64,
    tariff_id UInt64,
    day_kw UInt32,
    night_kw UInt32,
    month String,
    created_at DateTime
) ENGINE = Kafka
SETTINGS
    -- Адрес брокера Kafka внутри Docker-сети.
    kafka_broker_list = 'energosbyt_kafka:9092',
    -- Список топиков для подписки.
    kafka_topic_list = 'gauges',
    -- Имя consumer group (одна группа для одного узла).
    kafka_group_name = 'clickhouse_reader',
    -- Формат входящих JSON-сообщений.
    kafka_format = 'JSONEachRow',
    -- Количество параллельных потоков чтения на узле.
    -- 3 потока = 3 members в consumer group (по одному на партицию).
    kafka_num_consumers = 3,
    -- Максимальный размер блока (строк) для одной пачки INSERT из Kafka.
    -- Маленький блок = частые вставки: на демо счётчик fact_meter_readings
    -- заметно растёт при каждом повторном запросе, а не раз в ~10 секунд.
    kafka_max_block_size = 1000;
