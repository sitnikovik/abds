-- Загрузка data/dim_users.dataset.csv в dim_users_dist (Distributed обёртка).
-- Данные приходят через stdin (см. цель load-datasets в Makefile).
-- person_no из CSV становится id. Это временное решение для демо: в приложении id назначает snowflake.

-- insert_distributed_sync: фоновая отправка блоков через DistributedInsertQueue дублирует блоки при пересогласованиях. Синхронно — ровно одна копия на шард.
INSERT INTO energosbyt.dim_users_dist (id, fio, phone, phone_prev, created_at, updated_at) SETTINGS insert_distributed_sync = 1
SELECT person_no, fio, phone, phone_prev, created_at, updated_at
FROM input('person_no UInt64, fio String, phone String, phone_prev String, created_at DateTime, updated_at DateTime')
FORMAT CSVWithNames
