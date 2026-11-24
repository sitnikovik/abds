-- Distributed-обёртка fact_invoice поверх локальных таблиц на кластере.
-- Ключ шардирования realty_id: все счета одного объекта живут на одном шарде.
CREATE TABLE IF NOT EXISTS energosbyt.fact_invoice_dist ON CLUSTER energosbyt
AS energosbyt.fact_invoice
ENGINE = Distributed(energosbyt, energosbyt, fact_invoice, realty_id);
