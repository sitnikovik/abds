-- Репликация на боевой таблице.
-- Пишем в любую реплику, а строка реплицируется внутри шарда (ch1 и ch2 — это реплики шарда 1).

-- Состояние репликации dim_gauges на всех узлах.
-- Если absolute_delay = 0 и is_readonly = 0 на всех 6 узлах означает,
-- что система готова новые принимать записи.
SELECT
    hostName() AS node,
    is_leader,
    absolute_delay,
    is_readonly
FROM clusterAllReplicas('energosbyt', system.replicas)
WHERE database = 'energosbyt' AND table = 'dim_gauges'
ORDER BY node;

-- Добавляем ноый счётчик на ch1 в ЛОКАЛЬНУЮ таблицу.
-- (Идент-р 999001 не из датасета и предполагает, что его нет в таблице).
INSERT INTO energosbyt.dim_gauges
    (id, name, country, created_at)
VALUES (999001, 'Демо-счётчик из ch1', 'Россия', now());

-- 2.3. Смотрим с ch1, что строка есть в локальной таблице узла.
SELECT id, name FROM energosbyt.dim_gauges WHERE id = 999001;

-- Добавляем второй счётчик на ch2 в локальную таблицу.
-- ch1 увидит её же через несколько секунд.
INSERT INTO energosbyt.dim_gauges
    (id, name, country, created_at)
VALUES (999002, 'Демо-счётчик из ч2', 'Китай', now());

-- ============================================================================
-- Проверяем на ВСЕХ шести узлах.
-- У реплик (ch1, ch2) из шарда 1 должно быть 2 демо-строки; 
-- У остальных узлов 0, так как это "чужие" данные;
SELECT
    hostName() AS node,
    demo_ids
FROM clusterAllReplicas('energosbyt', view(
    SELECT hostName() AS h, uniqExactIf(id, id >= 999001) AS demo_ids
    FROM energosbyt.dim_gauges
    GROUP BY h
))
ORDER BY node;
