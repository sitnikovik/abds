-- Distributed-обёртка dim_users поверх локальных таблиц на кластере.
-- Ключ шардирования id: один ID всегда попадает в один и тот же шард.
CREATE TABLE IF NOT EXISTS energosbyt.dim_users_dist ON CLUSTER energosbyt
AS energosbyt.dim_users
ENGINE = Distributed(energosbyt, energosbyt, dim_users, id);
