-- Distributed-обёртка dim_tariffs поверх локальных таблиц на кластере.
-- Ключ шардирования id, версии одного тарифа с разными valid_from идут вместе.
CREATE TABLE IF NOT EXISTS energosbyt.dim_tariffs_dist ON CLUSTER energosbyt
AS energosbyt.dim_tariffs
ENGINE = Distributed(energosbyt, energosbyt, dim_tariffs, id);
