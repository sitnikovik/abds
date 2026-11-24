-- Distributed-обёртка fact_payments поверх локальных таблиц на кластере.
CREATE TABLE IF NOT EXISTS energosbyt.fact_payments_dist ON CLUSTER energosbyt
AS energosbyt.fact_payments
ENGINE = Distributed(energosbyt, energosbyt, fact_payments, realty_id);
