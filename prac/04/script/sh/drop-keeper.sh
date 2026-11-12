#!/bin/sh
# Чистит дерево таблиц /clickhouse/tables в ClickHouse Keeper.
# Без этого после teardown повторный make schema ловит REPLICA_ALREADY_EXISTS.
# Использование: ./script/sh/drop-keeper.sh

KEEPER="energosbyt_keeper1"
CH="energosbyt_ch1"

if [ "$(docker inspect -f '{{.State.Running}}' "$KEEPER" 2>/dev/null)" != "true" ]; then
    echo "Error: контейнер '$KEEPER' не запущен. Поднимите его: docker compose up -d" >&2
    exit 1
fi

# rmr удаляет дерево рекурсивно, y — подтверждение интерактивной команды rmr
printf 'rmr /clickhouse/tables 100000\ny\n' \
    | docker exec -i "$KEEPER" clickhouse-keeper-client --host localhost --port 9181 >/dev/null 2>&1

left=$(docker exec "$CH" clickhouse-client --user default --query \
    "SELECT count() FROM system.zookeeper WHERE path='/clickhouse/tables'")
if [ "$left" = "0" ]; then
    echo "Keeper чист."
else
    echo "Warning: в /clickhouse/tables осталось узлов: $left" >&2
    exit 1
fi
