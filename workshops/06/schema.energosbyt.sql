CREATE DATABASE IF NOT EXISTS energosbyt;

-- Таблица для чтения из Kafka
CREATE TABLE IF NOT EXISTS energosbyt.kafka_gauges (
    t1 UInt32,
    t2 UInt32,
    flat_id UInt32,
    id UInt32,
    sent_at DateTime
) ENGINE = Kafka(
    'kafka:9092',
    'gauges',
    'clickhouse_gauges_consumer',
    'JSONEachRow'
);

-- Таблица для хранения показаний
CREATE TABLE IF NOT EXISTS energosbyt.gauges (
    t1 UInt32,
    t2 UInt32,
    flat_id UInt32,
    id UInt32,
    sent_at DateTime
) ENGINE = MergeTree()
ORDER BY (flat_id, id);

-- Материализованное представление для передачи данных из Kafka в таблицу gauges
CREATE MATERIALIZED VIEW IF NOT EXISTS energosbyt.kafka_to_gauges_mv TO energosbyt.gauges AS
SELECT * FROM energosbyt.kafka_gauges;


-- Справочник квартир
CREATE TABLE IF NOT EXISTS energosbyt.flats(
    id UInt32,
    address String,
    district String,
    area Float32,
    building_type String,
) ENGINE = MergeTree()
ORDER BY id;

-- Справочник квартирантов
CREATE TABLE IF NOT EXISTS energosbyt.residents(
    id UInt32,
    fio String,
    flat_id UInt32,
    move_in_date Date,
    move_out_date Date DEFAULT '2077-01-01' -- дата из будущего чтобы избежать NULL
) ENGINE = MergeTree()
ORDER BY (id, flat_id);