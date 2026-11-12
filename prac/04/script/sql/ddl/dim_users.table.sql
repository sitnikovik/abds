-- Локальная таблица справочника пользователей на кластере.
-- ReplicatedMergeTree: каждая пара реплик одного шарда синхронизируется через Keeper.
-- Ключ таблицы в ZooKeeper: {shard} даёт свой путь каждому шарду, {replica} различает реплики внутри шарда.
CREATE TABLE IF NOT EXISTS energosbyt.dim_users ON CLUSTER energosbyt (
    id UInt64, -- Snowflake ID
    fio String, -- ФИО полностью (SCD 1: перезаписываем)
    phone String, -- Телефон (SCD 3: текущее значение)
    phone_prev String, -- Предыдущее значение phone (SCD 3)
    created_at DateTime, -- Дата создания записи
    updated_at DateTime -- Дата обновления записи (версия строки)
) ENGINE = ReplicatedMergeTree('/clickhouse/tables/{shard}/dim_users', '{replica}')
ORDER BY id;
