-- Distributed-обёртка dim_realty поверх локальных таблиц на кластере.
CREATE TABLE IF NOT EXISTS energosbyt.dim_realty_dist ON CLUSTER energosbyt
AS energosbyt.dim_realty
ENGINE = Distributed(energosbyt, energosbyt, dim_realty, id);
