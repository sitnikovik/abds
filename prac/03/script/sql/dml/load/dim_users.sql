-- Загрузка data/dim_users.csv в dim_users.
-- Данные приходят через stdin (см. цель .load в Makefile).
-- person_no из CSV становится id. Это временное решение для демо: в приложении id назначает snowflake.
INSERT INTO energosbyt.dim_users (id, fio, phone, phone_prev, created_at, updated_at)
SELECT person_no, fio, phone, phone_prev, created_at, updated_at
FROM input('person_no UInt64, fio String, phone String, phone_prev String, created_at DateTime, updated_at DateTime')
FORMAT CSVWithNames
