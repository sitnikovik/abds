CREATE TABLE IF NOT EXISTS energosbyt.dim_users (
    id UInt64, -- Showflake ID
    fio String, -- ФИО полностью, например, Иванов Иван Иванович (SCD 1: перезаписываем)
    phone String, -- Телефон (SCD 3: текущее значение)
    phone_prev String, -- Предыдущее значение phone (SCD 3)
    created_at DateTime, -- Дата создания записи
    updated_at DateTime -- Дата обновления записи (версия строки)
) ENGINE = ReplacingMergeTree(updated_at)
ORDER BY id;
