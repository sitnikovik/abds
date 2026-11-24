-- Distributed-обёртка fact_meter_readings поверх локальных таблиц на кластере.
-- Ключ шардирования realty_id: показания и версионные правки одного объекта вместе.
CREATE TABLE IF NOT EXISTS energosbyt.fact_meter_readings_dist ON CLUSTER energosbyt
AS energosbyt.fact_meter_readings
ENGINE = Distributed(energosbyt, energosbyt, fact_meter_readings, realty_id);
