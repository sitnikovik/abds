-- Отказ реплики и восстановление.
--
-- Демо-строки живут только в ch1 и ch2 (шард 1).
-- Падение ch1 бьёт ровно по одной из двух реплик шарда
-- и мы продолжаем читать и писать со второй реплики (ch2),
-- а после рестарта ch1 сама догоняет пропущенное.

-- Проверим строки до "аварии".
-- Только у ch1, ch2 есть 3 демо-строки (999001, 999002).
-- Проверяем на обоих веб-консолях — в ch1 и ch2.
SELECT
    hostName() AS node,
    demo_ids
FROM clusterAllReplicas('energosbyt', view(
    SELECT hostName() AS h, uniqExactIf(id, id >= 999001) AS demo_ids
    FROM energosbyt.dim_gauges
    GROUP BY h
))
ORDER BY node;

-- ================================
-- Роняем узел ch1 через Docker.
-- ================================


-- Читаем через ch2 тот же запрос выше — результат не изменился.
-- Приложению не нужно ничего перестраивать, так как отвечает живая реплика в шарде.

-- Пишем ВО ВРЕМЯ аварии на ch2.
INSERT INTO energosbyt.dim_gauges
    (id, name, country, created_at)
VALUES (999003, 'Демо-счётчик во время аварии', 'Тайланд', now());

-- =================================================================
-- Смотрим Grafana график "Живые узлы ClickHouse" — должно быть 5/6.
-- =================================================================

-- =================================================================
-- Поднимаем ch1 через Docker.
-- =================================================================

-- Проверяем состояние кластера.
-- Реплика догнала: absolute_delay = 0, очередь пустая везде.
SELECT
    hostName() AS node,
    is_readonly,
    absolute_delay,
    queue_size
FROM clusterAllReplicas('energosbyt', system.replicas)
WHERE database = 'energosbyt' AND table = 'dim_gauges'
ORDER BY node;

-- Проверяем консистентность.
-- У ch1 и ch2 — теперь 3 демо-строки (999001, 999002, 999003)
-- У остальных узлов — 0.
SELECT
    hostName() AS node,
    demo_ids
FROM clusterAllReplicas('energosbyt', view(
    SELECT hostName() AS h, uniqExactIf(id, id >= 999001) AS demo_ids
    FROM energosbyt.dim_gauges
    GROUP BY h
))
ORDER BY node;
