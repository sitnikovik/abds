-- Distributed-обёртка fact_gauge_bindings поверх локальных таблиц на кластере.
-- Ключ шардирования gauge_id: пара (установка, аннулирование) всегда в одном шарде.
CREATE TABLE IF NOT EXISTS energosbyt.fact_gauge_bindings_dist ON CLUSTER energosbyt
AS energosbyt.fact_gauge_bindings
ENGINE = Distributed(energosbyt, energosbyt, fact_gauge_bindings, gauge_id);
