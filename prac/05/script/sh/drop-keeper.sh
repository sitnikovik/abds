#!/bin/sh
# Чистит дерево таблиц /clickhouse/tables в ClickHouse Keeper.
# Без этого после teardown повторный make schema ловит REPLICA_ALREADY_EXISTS.
#
# Сначала дожидается, пока узлы доиграют фоновую очистку DROP TABLE:
# каждый узел асинхронно удаляет путь своей реплики в Keeper. Если сделать rmr
# раньше, отстающая очистка удалит пути, которые успел создать повторный
# make schema, и реплики останутся без метаданных (total_replicas = 0).
# Использование: ./script/sh/drop-keeper.sh [timeout_seconds]

KEEPER="energosbyt_keeper1"
CH="energosbyt_ch1"
TIMEOUT="${1:-30}"

if [ "$(docker inspect -f '{{.State.Running}}' "$KEEPER" 2>/dev/null)" != "true" ]; then
    echo "Error: контейнер '$KEEPER' не запущен. Поднимите его: docker compose up -d" >&2
    exit 1
fi

# Ждём, пока из всех шард-директорий исчезнут пути таблиц.
i=0
while [ "$i" -lt "$TIMEOUT" ]; do
    busy=$(docker exec "$CH" clickhouse-client --user default --query \
        "SELECT (SELECT count() FROM system.zookeeper WHERE path = '/clickhouse/tables/1') \
              + (SELECT count() FROM system.zookeeper WHERE path = '/clickhouse/tables/2') \
              + (SELECT count() FROM system.zookeeper WHERE path = '/clickhouse/tables/3')" \
        2>/dev/null || echo 1)
    if [ "$busy" = "0" ]; then
        break
    fi
    sleep 2
    i=$((i + 2))
done

# rmr удаляет дерево рекурсивно, y — подтверждение интерактивной команды rmr.
printf 'rmr /clickhouse/tables 100000\ny\n' \
    | docker exec -i "$KEEPER" clickhouse-keeper-client --host localhost --port 9181 >/dev/null 2>&1

left=$(docker exec "$CH" clickhouse-client --user default --query \
    "SELECT count() FROM system.zookeeper WHERE path='/clickhouse/tables'")
if [ "$left" = "0" ]; then
    echo "Keeper чист."
else
    echo "Error: в /clickhouse/tables осталось узлов: $left" >&2
    exit 1
fi
