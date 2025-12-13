.DEFAULT_GOAL = run

.PHONY: run
run: .cfg
	@echo "🐳 🔕 Запускаю проект в фоновом режиме..."
	docker-compose up -d --build
	@echo "✅ Всё ОК"

.PHONY: rund
rund: .cfg
	@echo "🐳 📣 Запускаю проект..."
	docker-compose up --build

.PHONY: stop
stop:
	@echo "🛑 Останавливаю проект..."
	docker-compose --profile manual down producer
	docker-compose down
	@echo "✅ Всё ОК"

.PHONY: venv
venv:
	@echo "🐍 Создаю виртуальное окружение Python и устанавливаю зависимости..."
	python3 -m venv .venv && \
	source .venv/bin/activate && \
	pip install --upgrade pip && \
	pip install -r requirements.txt --upgrade
	
.PHONY: .energosbyt
energosbyt: .schema .loadcsv .run-gauge-producer
	@echo "✅ Проект 'energosbyt' готов к работе!"

.PHONY: .run-gauge-producer
.run-gauge-producer:
	@echo "🐳 🔕 Запускаю продюсера показаний в фоновом режиме..."
	docker-compose --profile manual up -d --build producer
	@echo "✅ Всё ОК"

.PHONY: .rund-gauge-producer
.rund-gauge-producer:
	@echo "🐳 📣 Запускаю продюсера показаний..."
	docker-compose --profile manual up --build producer

.PHONY: .loadcsv
.loadcsv:
	@echo "🐳 🔕 Загружаю CSV данные в ClickHouse..."
	docker exec -i ch1 clickhouse-client --query \
	"TRUNCATE TABLE energosbyt.residents"
	docker exec -i ch1 clickhouse-client --query \
	 "INSERT INTO energosbyt.residents FORMAT CSVWithNames" < workshops/06/residents.energosbyt.csv
	docker exec -i ch1 clickhouse-client --query \
	"TRUNCATE TABLE energosbyt.flats"
	docker exec -i ch1 clickhouse-client --query \
	 "INSERT INTO energosbyt.flats FORMAT CSVWithNames" < workshops/06/flats.energosbyt.csv
	@echo "✅ Всё ОК"

.PHONY: .purge-kafka-gauges
.purge-kafka-gauges:
	@echo "🐳 🔕 Очищаю топики показаний в ClickHouse..."
	docker exec -it kafka kafka-topics \
	--bootstrap-server localhost:9092 --delete --topic gauges
	docker exec -it kafka kafka-topics \
	--bootstrap-server localhost:9092 --create --topic gauges --partitions 3 --replication-factor 1

.PHONY: .schema
.schema:
	@echo "🐳 Создаю новую БД 'energosbyt' ClickHouse..."
	docker exec -i ch1 clickhouse-client --multiquery < workshops/06/schema.energosbyt.sql

.PHONY: .cfg
.cfg:
	@echo "⚙️  Создаю файл конфигурации .env для Airflow..."
	cp airflow/.env.example airflow/.env
