.PHONY: run rund setup schema load-datasets gen-datasets teardown demo-reset kafka-schema streamer-start streamer-stop build-streamer demo-smoke-test

# Запускает проект в фоновом режиме.
run:
	@echo "🐳 🔕 Запускаю проект в фоновом режиме..."
	docker compose up -d

# Запускает проект в режиме отладки с просмотром всех логов.
rund:
	@echo "🐳 📣 Запускаю проект..."
	docker compose up

# Настройка и прогрев системы.
# Создаст схему БД и загрузит все датасеты в таблицы.
setup: schema load-datasets

# Полная очистка перед семинаром.
demo-reset: teardown schema load-datasets

# Создаёт схему реплицируемого кластера.
schema:
	@echo "Создаём таблицы на кластере..."
	@./script/sh/run_sql.sh \
		script/sql/ddl/dim_gauges.table.sql \
		script/sql/ddl/dim_realty.table.sql \
		script/sql/ddl/dim_tariff.table.sql \
		script/sql/ddl/dim_users.table.sql \
		script/sql/ddl/fact_gauge_bindings.table.sql \
		script/sql/ddl/fact_invoice.table.sql \
		script/sql/ddl/fact_meter_readings.table.sql \
		script/sql/ddl/fact_payments.table.sql \
		script/sql/ddl/dim_gauges.dist.table.sql \
		script/sql/ddl/dim_realty.dist.table.sql \
		script/sql/ddl/dim_tariff.dist.table.sql \
		script/sql/ddl/dim_users.dist.table.sql \
		script/sql/ddl/fact_gauge_bindings.dist.table.sql \
		script/sql/ddl/fact_invoice.dist.table.sql \
		script/sql/ddl/fact_meter_readings.dist.table.sql \
		script/sql/ddl/fact_payments.dist.table.sql

# Создаёт таблицу Kafka и материализованное представление.
kafka-schema:
	@echo "Создаём Kafka интеграцию..."
	@./script/sh/run_sql.sh \
		script/sql/ddl/kafka_gauge_readings.table.sql \
		script/sql/ddl/mv_kafka_to_meter_readings.sql

# Генерирует большие датасеты в data/.
gen-datasets:
	@mkdir -p data
	@python3 script/py/gen_dim_users.py > data/dim_users.dataset.csv
	@python3 script/py/gen_energosbyt.py

# Загружает датасеты во все таблицы через Distributed-обёртки.
# Данные размазываются по шардам и реплицируются внутри шарда.
load-datasets:
	@for t in dim_users dim_realty dim_gauges dim_tariffs fact_gauge_bindings fact_meter_readings fact_invoice fact_payments; do \
		echo "Загружаем $$t..."; \
		docker exec -i energosbyt_ch1 clickhouse-client --user default --database energosbyt \
			--query "$$(cat script/sql/dml/load/$$t.sql)" < data/$$t.dataset.csv || exit 1; \
	done
	@echo "Проверяем распределение фактов по шардам..."
	@docker exec energosbyt_ch1 clickhouse-client --user default --database energosbyt \
		--query "SELECT _shard_num, n FROM cluster('energosbyt', view(SELECT hostName() AS host, count() AS n FROM energosbyt.fact_invoice)) ORDER BY _shard_num"

# Собирает бинарь стримера локально в bin/
build-streamer:
	@echo "🔨 Собираем стример..."
	@mkdir -p bin
	@cd app/energosbyt && go build -o ../../bin/streamer ./cmd/streamer
	@echo "✅ Бинарь собран: bin/streamer"

# Запускает стример показаний (отправка в Kafka).
streamer-start:
	@echo "🚀 Запускаем стример показаний..."
	@echo "⏳ Ждём готовности Kafka..."
	@until docker exec energosbyt_kafka kafka-topics --bootstrap-server localhost:9092 --list >/dev/null 2>&1; do \
		echo "Kafka ещё не готова, жду..."; \
		sleep 2; \
	done
	@echo "✅ Kafka готова"
	@echo "⏳ Ждём готовности ClickHouse..."
	@until docker exec energosbyt_ch1 clickhouse-client --query "SELECT 1" >/dev/null 2>&1; do \
		echo "ClickHouse ещё не готов, жду..."; \
		sleep 2; \
	done
	@echo "✅ ClickHouse готов"
	docker compose --profile manual up -d energosbyt_streamer

# Останавливает стример.
streamer-stop:
	@echo "⏹️  Останавливаем стример..."
	docker compose stop energosbyt_streamer

# Очистка системы.
# Удаляет все таблицы на кластере, чистит Keeper и сбрасывает очередь Kafka.
# Порядок важен: rmr в Keeper выполняется последним, иначе отстающая фоновая
# очистка DROP TABLE удалит пути, которые успел создать повторный make schema.
teardown: .drop-ch .drop-wait .drop-kafka .drop-keeper

# Удаляет таблицы на кластере.
.drop-ch:
	@echo "Удаляем таблицы на кластере..."
	@./script/sh/run_sql.sh \
		script/sql/dml/drop_all.sql

# Чистит дерево таблиц в Keeper: без этого повторный schema ловит REPLICA_ALREADY_EXISTS.
.drop-keeper:
	@echo "Чистим дерево таблиц в Keeper..."
	@./script/sh/drop-keeper.sh

# Сбрасывает очередь Kafka: удаляет топик gauges и consumer group.
# Топик пересоздастся автоматически при первом сообщении стримера,
# поэтому чтение всегда начинается с чистой очереди без накопленного lag.
.drop-kafka:
	@echo "Сбрасываем очередь Kafka (топик gauges и consumer group)..."
	@docker exec energosbyt_kafka kafka-topics --bootstrap-server localhost:9092 --delete --topic gauges >/dev/null 2>&1 || true
	@docker exec energosbyt_kafka kafka-consumer-groups --bootstrap-server localhost:9092 --delete --group clickhouse_reader >/dev/null 2>&1 || true

# Ждёт, пока DDL-очередь опустеет на всех узлах.
.drop-wait:
	@echo "Ждём опустошения DDL-очереди..."
	@./script/sh/wait-ddl-queue.sh

