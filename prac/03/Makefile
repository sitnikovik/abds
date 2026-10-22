.PHONY: run rund setup schema teardown load-datasets gen-datasets

# Запускает проект в фоновом режиме.
run:
	@echo "🐳 🔕 Запускаю проект в фоновом режиме..."
	docker-compose up -d --build

# Запускает проект в режиме отладки с просмотром всех логов.
rund:
	@echo "🐳 📣 Запускаю проект..."
	docker-compose up --build

# Настройка и прогрев системы.
# Создаст схему БД, сгенерирует и загрузит все датасеты в таблицы.
setup: schema gen-datasets load-datasets

# Очистка системы
teardown: .drop-ch

# Создаёт схему базы данных.
schema:
	@echo "Создаём таблицы в ClickHouse..."
	@./script/sh/run_sql.sh \
		script/sql/ddl/dim_gauges.table.sql \
		script/sql/ddl/dim_realty.table.sql \
		script/sql/ddl/dim_tariff.table.sql \
		script/sql/ddl/dim_users.table.sql \
		script/sql/ddl/fact_gauge_bindings.table.sql \
		script/sql/ddl/fact_invoice.table.sql \
		script/sql/ddl/fact_meter_readings.table.sql \
		script/sql/ddl/fact_payments.table.sql

.drop-ch:
	@echo "Удаляем таблицы в ClickHouse..."
	@./script/sh/run_sql.sh \
		script/sql/dml/drop_all.sql

# Загружает датасеты во все таблицы
# с предварительной очисткой базы и пересозданием схемы.
load-datasets: teardown schema
	@for t in dim_users dim_realty dim_gauges dim_tariffs fact_gauge_bindings fact_meter_readings fact_invoice fact_payments; do \
		echo "Загружаем $$t..."; \
		docker exec -i energosbyt_ch1 clickhouse-client --user default --database energosbyt \
			--query "$$(cat script/sql/dml/load/$$t.sql)" < data/$$t.dataset.csv || exit 1; \
	done

# Генерирует большие датасеты в data/.
gen-datasets: .gen-users .gen-all
	@python3 script/py/gen_energosbyt.py

# Генерирует датасет пользователей в data/.
.gen-dataset-users:
	@python3 script/py/gen_dim_users.py > data/dim_users.dataset.csv

# Генерирует все датасеты (нужен data/dim_users.csv).
.gen-dataset-all:
	@python3 script/py/gen_energosbyt.py

