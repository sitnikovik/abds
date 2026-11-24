-- Демонстрация работы Kafka + ClickHouse.
-- Проверяем, что ch1 читает из Kafka (таблица-движок Kafka и её потребители)
SELECT
    'Kafka consumers' AS component,
    table,
    num_messages_read,
    last_poll_time
FROM system.kafka_consumers
WHERE database = 'energosbyt';

-- Проверяем материализованное представление
SELECT 'Materialized view' AS component, name FROM system.tables WHERE database = 'energosbyt' AND name LIKE 'mv_kafka%';

-- Смотрим данные в основной таблице
-- (число записей растёт каждую секунду, а max(created_at) продвигается по симулированному времени).
SELECT 
    'Meter readings from Kafka' AS source,
    count() AS total_records,
    min(created_at) AS first_reading,
    max(created_at) AS last_reading,
    uniq(gauge_id) AS unique_gauges,
    uniq(realty_id) AS unique_realties
FROM energosbyt.fact_meter_readings;

-- Распределение по месяцам.
-- Месяцы добавляются примерно по одному в секунду
-- (5 лет симулируются за 60 секунд).
SELECT 
    reading_month,
    count() AS readings_count,
    sum(day_kw + night_kw) AS total_kwh
FROM energosbyt.fact_meter_readings
WHERE reading_month >= '2020-01-01'
GROUP BY reading_month
ORDER BY reading_month;

-- Топ-10 счётчиков по потреблению (симулированная эпоха с 2020-01-01).
SELECT 
    gauge_id,
    sum(day_kw + night_kw) AS total_consumption,
    count() AS reading_count
FROM energosbyt.fact_meter_readings
WHERE reading_month >= '2020-01-01'
GROUP BY gauge_id
ORDER BY total_consumption DESC
LIMIT 10;

-- Проверка ошибок в материализованном представлении.
SELECT 
    query_duration_ms,
    exception_code,
    exception
FROM system.query_log
WHERE database = 'energosbyt'
  AND query LIKE '%mv_kafka_to_meter_readings%'
  AND exception != ''
ORDER BY event_time DESC
LIMIT 10;

-- Распределение данных по шардам кластера (все реплики).
--
-- Потоковые данные (симул. эпоха с 2020-01-01) попадают только на шард 1:
-- ch1 и ch2 показывают одинаковые растущие числа (репликация),
-- остальные узлы содержат только статичный исторический датасет.
SELECT 
    hostName() AS host,
    count() AS local_rows
FROM clusterAllReplicas('energosbyt', energosbyt.fact_meter_readings)
WHERE reading_month >= '2020-01-01'
GROUP BY host
ORDER BY host;
