#!/bin/sh
# Ждёт, пока DDL-очередь на всех узлах кластера опустеет.
# Без этого: после teardown отставшие узлы доигрывают DROP-команды
# и удяляют деревья таблиц, которые успел создать в Keeper.
# Использование: ./script/sh/wait-ddl-queue.sh [timeout_seconds]

CH="energosbyt_ch1"
TIMEOUT="${1:-60}"
i=0

while [ "$i" -lt "$TIMEOUT" ]; do
    busy=$(docker exec "$CH" clickhouse-client --user default --query \
        "SELECT count() FROM clusterAllReplicas('energosbyt', system.distributed_ddl_queue) \
         WHERE status != 'Finished' AND query_create_time > now() - INTERVAL 15 MINUTE" 2>/dev/null || echo 1)
    if [ "$busy" = "0" ]; then
        exit 0
    fi
    sleep 2
    i=$((i + 2))
done

echo "Error: DDL-очередь не опустела за $TIMEOUT секунд" >&2
exit 1
