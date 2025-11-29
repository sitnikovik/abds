#!/bin/bash

# Скрипт инициализации Superset с поддержкой ClickHouse
# Выполняется ВНУТРИ контейнера при запуске

echo "=== Инициализация Superset ==="

# Проверка, нужна ли инициализация (если БД уже существует, пропускаем)
if [ ! -f /app/superset_home/superset.db ]; then
    echo "Первый запуск - инициализация БД..."
    # Обновление БД
    echo "Обновление базы данных..."
    superset db upgrade
    # Создание администратора
    echo "Создание администратора..."
    superset fab create-admin \
        --username admin \
        --firstname Admin \
        --lastname User \
        --email admin@superset.com \
        --password admin
    # Инициализация Superset
    echo "Инициализация Superset..."
    superset init
    echo "=== Инициализация завершена ==="
else
    echo "БД уже существует, пропускаем инициализацию"
fi
echo "Запуск Superset на порту 8088..."
echo "Доступ: http://localhost:8088 (admin/admin)"
echo "ClickHouse: clickhousedb://default@ch1:8123/default"

# Запуск сервера
exec superset run -h 0.0.0.0 -p 8088
