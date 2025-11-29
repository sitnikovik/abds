FROM apache/superset:4.1.4

# Устанавливаем clickhouse-connect
USER root
RUN pip install --no-cache-dir clickhouse-connect

# Копируем скрипт инициализации в докер
COPY init.sh /app/superset-init.sh
RUN chmod +x /app/superset-init.sh

# Инициализация Superset
USER superset
ENTRYPOINT ["/bin/bash", "/app/superset-init.sh"]
