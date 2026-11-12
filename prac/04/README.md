# Семинар 4

Тема: кластер ClickHouse — репликация и шардирование.

## С чего начать

### Запуск

**Запуск контейнеров:**

```sh
make run
```

Поднимается кластер из 6 узлов ClickHouse (3 шарда по 2 реплики), ClickHouse Keeper,
а также Prometheus и Grafana с дашбордом для наблюдения за кластером.

Все томы монтируются в корне локально для ручной отладки и изучения устройства БД.

### Настройка и прогрев

**Полная настройка:**

Создаст схему БД на кластере, сгенерирует и загрузит все датасеты в таблицы.

```sh
make setup
```

**Выборочная настройка:**

**Создание схемы:**

Создаёт локальные и Distributed-таблицы на кластере.

```sh
make schema
```

**Генерация датасетов:**

Сгенерирует датасеты в `data/`

```sh
make gen-datasets
```

**Загрузка датасетов:**

Загрузит сгенерированные датасеты из `data/` в таблицы через Distributed-обёртки.
Данные размазываются по шардам и реплицируются внутри шарда.

```sh
make load-datasets
```

### Очистка

Удалит все таблицы на кластере и снимки в Keeper и восстановит БД в изначальное состояние.

```sh
make teardown
```

**Сброс демо:**

Полная очистка и повторная настройка перед семинаром.

```sh
make demo-reset
```

## Конфиги

Все конфиги лежат в [configs/](/prac/04/configs) и монтируются контейнерам через `docker-compose.yml`.

**ClickHouse — [configs/clickhouse/](/prac/04/configs/clickhouse):**

- [remote-servers.xml](/prac/04/configs/clickhouse/remote-servers.xml) — топология кластера
  `energosbyt` (3 шарда по 2 реплики) и подключение к трём узлам ClickHouse Keeper:
  через них реплики находят друг друга, хранят метаданные и DDL-очередь;
- [macros/1.xml … 6.xml](/prac/04/configs/clickhouse/macros) — макросы каждого узла ch1…ch6:
  `{cluster}`, `{shard}`, `{replica}` подставляются в пути `ReplicatedMergeTree` в Keeper,
  чтобы реплики не претендовали на один путь;
- [prometheus.xml](/prac/04/configs/clickhouse/prometheus.xml) — встроенный экспортер метрик:
  ClickHouse сам отдаёт метрики в формате Prometheus на порту 9363, отдельный экспортер не нужен.

**ClickHouse Keeper — [configs/keeper/](/prac/04/configs/keeper):**

- [1.xml, 2.xml, 3.xml](/prac/04/configs/keeper) — конфиги трёх узлов Keeper (raft-кворум):
  список участников и их raft-порты, порты для клиентов, пути к данным.

**Prometheus — [configs/prometheus/prometheus.yml](/prac/04/configs/prometheus/prometheus.yml):**

Собирает метрики со всех шести узлов ClickHouse (`energosbyt_ch1…ch6:9363`) и узлов Keeper каждые 5 секунд.

**Grafana — [configs/grafana/](/prac/04/configs/grafana):**

- [provisioning/datasources/prometheus.yml](/prac/04/configs/grafana/provisioning/datasources/prometheus.yml) —
  автоматическая регистрация Prometheus как источника данных Grafana;
- [provisioning/dashboards/dashboards.yml](/prac/04/configs/grafana/provisioning/dashboards/dashboards.yml) —
  провайдер дашбордов: загружает JSON-дашборды из каталога на диске;
- [dashboards/energosbyt.json](/prac/04/configs/grafana/dashboards/energosbyt.json) —
  дашборд с метриками кластера energosbyt.

## SQL

Все используемые скрипты SQL смотрите в [script/sql/](/prac/04/script/sql).

Скрипты для демо вживую — в [script/sql/demo/](/prac/04/script/sql/demo):

1. [01.architecture.sql](/prac/04/script/sql/demo/01.architecture.sql)
— топология кластера, макросы, здоровье репликации;
2. [02.replication.sql](/prac/04/script/sql/demo/02.replication.sql)
— репликация внутри шарда: пишем в любую реплику, читаем со всех;
3. [03.replication_failover.sql](/prac/04/script/sql/demo/03.replication_failover.sql)
— отказ реплики и восстановление;
4. [04.sharding.sql](/prac/04/script/sql/demo/04.sharding.sql)
— шардирование и распределение данных по шардам.
