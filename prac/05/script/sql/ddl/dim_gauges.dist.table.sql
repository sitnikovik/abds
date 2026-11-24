-- Distributed-обёртка dim_gauges поверх локальных таблиц на кластере.
CREATE TABLE IF NOT EXISTS energosbyt.dim_gauges_dist ON CLUSTER energosbyt
AS energosbyt.dim_gauges
ENGINE = Distributed(energosbyt, energosbyt, dim_gauges, id);
