.PHONY: run rund setup schema load-datasets gen-datasets teardown demo-reset

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

# Очистка системы.
# Удаляет все таблицы на кластере и снимки в Keeper.
teardown: .drop-ch .drop-keeper .drop-wait

# Удаляет таблицы на кластере.
.drop-ch:
	@echo "Удаляем таблицы на кластере..."
	@./script/sh/run_sql.sh \
		script/sql/dml/drop_all.sql

# Чистит дерево таблиц в Keeper: без этого повторный schema ловит REPLICA_ALREADY_EXISTS.
.drop-keeper:
	@echo "Чистим дерево таблиц в Keeper..."
	@./script/sh/drop-keeper.sh

# Ждёт, пока DDL-очередь опустеет на всех узлах.
.drop-wait:
	@echo "Ждём опустошения DDL-очереди..."
	@./script/sh/wait-ddl-queue.sh
