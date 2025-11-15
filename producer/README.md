# Универсальный Kafka продюсер (Go)

Запускает поток событий для разных учебных сценариев: `energy`, `iot` или произвольный `template`.

Примеры:

```bash
go run . --brokers redpanda:9092 --topic energy_readings --scenario energy --rate 20/s --acks all

go run . --brokers redpanda:9092 --topic iot_events --scenario iot --rate 50/s --key-template "{{ .Key }}"

go run . --brokers redpanda:9092 --topic events --scenario template --template ./payload.tmpl --rate 10/s
```

Флаги (основные):
-
- --brokers (строка) — список брокеров, по умолчанию `redpanda:9092`
- --topic (строка) — имя топика
- --scenario (energy|iot|template)
- --rate (например `20/s`, `100/m`)
- --duration (например `10m`, `0` для бесконечно)
- --max-messages — ограничить количество сообщений
- --acks (none|leader|all)
- --compression (none|gzip|snappy|lz4|zstd)
- --key-template — шаблон для ключа, по умолчанию `{{ .Key }}`
- --template — путь к файлу шаблона для scenario=template

Docker:

```bash
docker build -t abds/producer:latest ./producer
docker run --rm --network abds-private_default abds/producer:latest \
  --brokers redpanda:9092 --topic energy_readings --scenario energy --rate 20/s
```
