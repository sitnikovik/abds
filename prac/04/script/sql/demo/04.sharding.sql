-- Шардирование
-- Как Distributed-таблица (координатор) раскладывает строки и собирает их на чтении.

-- Вывод таблиц-координаторов в кластере.
SELECT
    hostName() AS node,
    table,
    engine
FROM clusterAllReplicas('energosbyt', system.tables)
WHERE database = 'energosbyt' AND engine LIKE '%Distributed%'
ORDER BY node, table;

-- По какому критерию строка идёт в шард: hash(realty_id) % 3.
-- Все строки одного объекта лежат на одном шарде.
SELECT
    _shard_num,
    uniqExact(realty_id) AS unique_objects,
    count() AS rows_in_shard
FROM cluster('energosbyt', view(
    SELECT realty_id FROM energosbyt.fact_invoice WHERE realty_id < 999000
))
GROUP BY _shard_num
ORDER BY _shard_num;

-- Подсчёт выставленных счетов на оплату.
-- Один объект живёт ровно на одном шарде (realty_id=3 — объект из датасета).
-- 
-- ВАЖНО про агрегаты.
-- Внутренний view даёт ОДНУ строку на шард (count()=n),
-- поэтому во внешнем запросе суммируем числа через sum(n).
-- Если бы чужой шард давал 0 строк, то count() посчитал бы их как 1.
-- Это ловушка, в которую легко попасть при любом clusterAllReplicas-запросе с агрегатом внутри view.
-- Ожидаем, что ровно один шард с полным числом счетов объекта (~93), у остальных 0.
SELECT
    _shard_num,
    sum(n) AS rows_for_object
FROM cluster('energosbyt', view(
    SELECT count() AS n FROM energosbyt.fact_invoice WHERE realty_id = 3
))
GROUP BY _shard_num
ORDER BY _shard_num;

-- То же самое одним запросом через координатора.
SELECT count() FROM energosbyt.fact_invoice_dist WHERE realty_id = 3;

-- Подсчёт всех через координатор = сумма трёх шардов.
-- internal_replication=true в конфиге: от каждой пары реплик считаем ОДИН раз.
SELECT count() FROM energosbyt.fact_invoice_dist;

-- Минимум/максимум объекта на каждом шарде.
-- Объекты не пересекаются между шардами.
SELECT
    _shard_num,
    min(realty_id) AS min_id,
    max(realty_id) AS max_id,
    count() AS rows
FROM cluster('energosbyt', view(
    SELECT realty_id FROM energosbyt.fact_invoice WHERE realty_id < 999000
))
GROUP BY _shard_num
ORDER BY _shard_num;
