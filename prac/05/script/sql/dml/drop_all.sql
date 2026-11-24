-- Полная очистка кластера «ЭнергоСбыт» (семинар 5).
-- Порядок: сначала Distributed-обёртки, затем локальные таблицы. Базу не трогаем — с ней может работать приложение.
-- SYNC обязателен: без него фоновая очистка DROP TABLE удаляет пути реплик в Keeper
-- асинхронно и может удалить пути, созданные повторным запуском make schema.
DROP TABLE IF EXISTS energosbyt.mv_kafka_to_meter_readings ON CLUSTER energosbyt SYNC;
DROP TABLE IF EXISTS energosbyt.kafka_gauge_readings ON CLUSTER energosbyt SYNC;
DROP TABLE IF EXISTS energosbyt.dim_gauges_dist ON CLUSTER energosbyt SYNC;
DROP TABLE IF EXISTS energosbyt.dim_realty_dist ON CLUSTER energosbyt SYNC;
DROP TABLE IF EXISTS energosbyt.dim_tariffs_dist ON CLUSTER energosbyt SYNC;
DROP TABLE IF EXISTS energosbyt.dim_users_dist ON CLUSTER energosbyt SYNC;
DROP TABLE IF EXISTS energosbyt.fact_gauge_bindings_dist ON CLUSTER energosbyt SYNC;
DROP TABLE IF EXISTS energosbyt.fact_invoice_dist ON CLUSTER energosbyt SYNC;
DROP TABLE IF EXISTS energosbyt.fact_meter_readings_dist ON CLUSTER energosbyt SYNC;
DROP TABLE IF EXISTS energosbyt.fact_payments_dist ON CLUSTER energosbyt SYNC;
DROP TABLE IF EXISTS energosbyt.dim_gauges ON CLUSTER energosbyt SYNC;
DROP TABLE IF EXISTS energosbyt.dim_realty ON CLUSTER energosbyt SYNC;
DROP TABLE IF EXISTS energosbyt.dim_tariffs ON CLUSTER energosbyt SYNC;
DROP TABLE IF EXISTS energosbyt.dim_users ON CLUSTER energosbyt SYNC;
DROP TABLE IF EXISTS energosbyt.fact_gauge_bindings ON CLUSTER energosbyt SYNC;
DROP TABLE IF EXISTS energosbyt.fact_invoice ON CLUSTER energosbyt SYNC;
DROP TABLE IF EXISTS energosbyt.fact_meter_readings ON CLUSTER energosbyt SYNC;
DROP TABLE IF EXISTS energosbyt.fact_payments ON CLUSTER energosbyt SYNC;
