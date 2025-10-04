.DEFAULT_GOAL = run

.PHONY: run
run:
	@echo "🐳 🔕 Запускаю проект в фоновом режиме..."
	docker-compose up -d --build
	@echo "✅ Всё ОК"

.PHONY: rund
rund:
	@echo "🐳 📣 Запускаю проект..."
	docker-compose up --build

.PHONY: stop
stop:
	@echo "🛑 Останавливаю проект..."
	docker-compose down
	@echo "✅ Всё ОК"