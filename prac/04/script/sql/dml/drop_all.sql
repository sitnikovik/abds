-- Полная очистка кластера «ЭнергоСбыт» (семинар 4).
-- Порядок: сначала Distributed-обёртки, затем локальные таблицы. Базу не трогаем — с ней может работать приложение.
DROP TABLE IF EXISTS energosbyt.dim_gauges_dist ON CLUSTER energosbyt;
DROP TABLE IF EXISTS energosbyt.dim_realty_dist ON CLUSTER energosbyt;
DROP TABLE IF EXISTS energosbyt.dim_tariffs_dist ON CLUSTER energosbyt;
DROP TABLE IF EXISTS energosbyt.dim_users_dist ON CLUSTER energosbyt;
DROP TABLE IF EXISTS energosbyt.fact_gauge_bindings_dist ON CLUSTER energosbyt;
DROP TABLE IF EXISTS energosbyt.fact_invoice_dist ON CLUSTER energosbyt;
DROP TABLE IF EXISTS energosbyt.fact_meter_readings_dist ON CLUSTER energosbyt;
DROP TABLE IF EXISTS energosbyt.fact_payments_dist ON CLUSTER energosbyt;
DROP TABLE IF EXISTS energosbyt.dim_gauges ON CLUSTER energosbyt;
DROP TABLE IF EXISTS energosbyt.dim_realty ON CLUSTER energosbyt;
DROP TABLE IF EXISTS energosbyt.dim_tariffs ON CLUSTER energosbyt;
DROP TABLE IF EXISTS energosbyt.dim_users ON CLUSTER energosbyt;
DROP TABLE IF EXISTS energosbyt.fact_gauge_bindings ON CLUSTER energosbyt;
DROP TABLE IF EXISTS energosbyt.fact_invoice ON CLUSTER energosbyt;
DROP TABLE IF EXISTS energosbyt.fact_meter_readings ON CLUSTER energosbyt;
DROP TABLE IF EXISTS energosbyt.fact_payments ON CLUSTER energosbyt;
