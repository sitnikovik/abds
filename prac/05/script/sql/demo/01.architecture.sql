-- Топология: 3 шарда по 2 реплики (ch1..ch6).
-- Если убрать WHERE, появится и строка default/localhost — это сам сервер,
-- она всегда есть в system.clusters, но к нашему кластеру не относится.
SELECT cluster, shard_num, replica_num, host_name, host_address, port
FROM system.clusters
WHERE cluster = 'energosbyt'
ORDER BY shard_num, replica_num;

-- ============================================================================
-- Макросы текущего узла.
-- {shard} и {replica} подставляются в пути ReplicatedMergeTree в Keeper:
--   + На ch1 ожидаем cluster=energosbyt, shard=1, replica=1;
--   + На ch2 ожидаем cluster=energosbyt, shard=1, replica=2;
SELECT * FROM system.macros;

-- Все узлы видят друг друга
-- Cчётчик таблиц нашей базы на каждом узле.
-- Если план кластера применён везде, число одинаковое на всех 6 узлах.
SELECT
    hostName() AS node,
    count() AS tables_on_node
FROM clusterAllReplicas('energosbyt', view(
    SELECT name FROM system.tables WHERE database = 'energosbyt'
))
GROUP BY node
ORDER BY node;

-- Здоровье репликации по каждому узлу.
-- Здоровый кластер: absolute_delay = 0 (реплика не отстаёт) и is_readonly = 0.
SELECT
    table AS replica_table,
    hostName() AS node,
    replica_name,
    absolute_delay,
    is_readonly
FROM clusterAllReplicas('energosbyt', system.replicas)
WHERE database = 'energosbyt'
ORDER BY replica_table, node;
