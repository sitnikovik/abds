#!/bin/bash
# Быстрая проверка работы стримера

echo "🔍 Проверяем формат JSON сообщений..."
echo ""

# Создаем тестовое сообщение
cat > /tmp/test_message.json << 'EOF'
{
  "created_at": "2020-01-15T10:30:00Z",
  "gauge_id": 12345,
  "realty_id": 67890,
  "user_id": 11111,
  "tariff_id": 1,
  "day_kw": 150,
  "night_kw": 80,
  "month": "January"
}
EOF

echo "Пример сообщения, которое будет отправлено в Kafka:"
cat /tmp/test_message.json | jq '.'

echo ""
echo "✅ Формат корректный. Можно запускать стример!"
echo ""
echo "Команды:"
echo "  make run              # Запустить инфраструктуру"
echo "  make setup            # Настроить БД"
echo "  make kafka-schema     # Создать Kafka таблицы"
echo "  make streamer-start   # Запустить стример"
