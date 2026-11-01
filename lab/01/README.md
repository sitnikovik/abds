# Лабораторная работа №1. SCD в ClickHouse

## Подготовка

⚠️ **Обязательно ознакомьтесь с
[CONTRIBUTING.md](https://github.com/sitnikovik/abds-template?tab=contributing-ov-file)** —
там описан процесс работы с репозиторием, Pull Requests и GitHub Actions.

## Цель работы

- Поднять один инстанс [ClickHouse](https://clickhouse.com/) в Docker.
- Смоделировать данные ЖКХ-компании «ЭнергоСбыт» по схеме «звезда».
- Применить SCD 1, SCD 2 и SCD 3 на практике.
- Научиться удалять и аннулировать данные в таблицах
  [MergeTree](https://clickhouse.com/docs/reference/engines/table-engines/mergetree-family/mergetree),
  [ReplacingMergeTree](https://clickhouse.com/docs/reference/engines/table-engines/mergetree-family/replacingmergetree),
  [VersionedCollapsingMergeTree](https://clickhouse.com/docs/reference/engines/table-engines/mergetree-family/versionedcollapsingmergetree).

## Предметная область

«ЭнергоСбыт» поставляет электричество в жилые дома и ведёт учёт показаний счётчиков.
Объектом учёта может быть квартира или нежилое помещение. К объекту привязан счётчик,
за объект платит собственник, тариф определяет стоимость киловатт-часа.

## Задание

### ClickHouse в Docker

Создайте `docker-compose.yml`, который поднимает **один** инстанс ClickHouse:

- Образ: `clickhouse/clickhouse-server:26.7`.
- Имя сервиса: `clickhouse`.
- Порты: `8123` (HTTP-интерфейс) и `9000` (нативный протокол).
- Данные ClickHouse хранятся в локальной папке `volume/` (она не коммитится).
- Все параметры подключения берутся из `.env.local`.

Репликация и шардирование в этой работе не нужны.

### Схема данных

Создайте таблицы в базе `energosbyt`. Имена таблиц и колонок, указанные ниже, менять нельзя.
Дополнительные колонки добавлять можно.

| Таблица | Назначение | Движок |
| --- | --- | --- |
| `dim_realty` | Объекты недвижимости | `MergeTree` |
| `dim_gauges` | Счётчики | `MergeTree` |
| `dim_tariffs` | Тарифы (SCD 2) | `ReplacingMergeTree(updated_at)` |
| `dim_users` | Собственники (SCD 1 и SCD 3) | `ReplacingMergeTree(updated_at)` |
| `fact_gauge_bindings` | Установка счётчика на объект | `VersionedCollapsingMergeTree(sign, version)` |
| `fact_meter_readings` | Показания счётчиков по дням | `ReplacingMergeTree(created_at)` |
| `fact_invoice` | Счета | `MergeTree` |

#### `dim_realty`

| Колонка | Тип | Описание |
| --- | --- | --- |
| `id` | `UInt64` | Идентификатор объекта |
| `address` | `String` | Адрес дома |
| `square` | `UInt32` | Площадь, м² |
| `type` | `String` | `flat` — квартира, `nrcp` — нежилое помещение |
| `number` | `String` | Номер квартиры (для нежилых помещений пусто) |
| `rooms` | `UInt8` | Количество комнат |
| `floor` | `UInt8` | Этаж |
| `porch` | `UInt8` | Подъезд |
| `created_at` | `DateTime` | Дата создания записи |

#### `dim_gauges`

| Колонка | Тип | Описание |
| --- | --- | --- |
| `id` | `UInt64` | Идентификатор счётчика |
| `name` | `String` | Модель счётчика |
| `country` | `String` | Страна производства |
| `created_at` | `DateTime` | Дата создания записи |

#### `dim_tariffs`

| Колонка | Тип | Описание |
| --- | --- | --- |
| `id` | `UInt64` | Идентификатор тарифа |
| `name` | `String` | Название тарифа |
| `price_day_kwph` | `UInt32` | Цена кВт·ч в дневное время, копейки |
| `price_night_kwph` | `UInt32` | Цена кВт·ч в ночное время, копейки |
| `valid_from` | `DateTime` | Начало действия версии тарифа |
| `valid_to` | `DateTime` | Конец действия версии тарифа (`2099-01-01 00:00:00` у открытой версии) |
| `created_at` | `DateTime` | Дата создания записи |
| `updated_at` | `DateTime` | Версия строки для `ReplacingMergeTree` |

#### `dim_users`

| Колонка | Тип | Описание |
| --- | --- | --- |
| `id` | `UInt64` | Идентификатор собственника |
| `fio` | `String` | ФИО (SCD 1: перезаписывается) |
| `phone` | `String` | Текущий телефон (SCD 3) |
| `phone_prev` | `String` | Предыдущий телефон, пусто если его не было (SCD 3) |
| `created_at` | `DateTime` | Дата создания записи |
| `updated_at` | `DateTime` | Версия строки для `ReplacingMergeTree` |

#### `fact_gauge_bindings`

| Колонка | Тип | Описание |
| --- | --- | --- |
| `gauge_id` | `UInt64` | Идентификатор счётчика |
| `realty_id` | `UInt64` | Идентификатор объекта |
| `created_at` | `DateTime` | Дата создания записи |
| `sign` | `Int8` | `1` — установка, `-1` — аннулирование |
| `version` | `UInt64` | Версия привязки, см. правила ниже |

#### `fact_meter_readings`

| Колонка | Тип | Описание |
| --- | --- | --- |
| `realty_id` | `UInt64` | Идентификатор объекта |
| `user_id` | `UInt64` | Идентификатор собственника |
| `tariff_id` | `UInt64` | Идентификатор тарифа, по которому считается счёт |
| `gauge_id` | `UInt64` | Идентификатор счётчика |
| `reading_date` | `Date` | Дата показания |
| `day_kw` | `UInt32` | Показание дневного счётчика, кВт·ч |
| `night_kw` | `UInt32` | Показание ночного счётчика, кВт·ч |
| `created_at` | `DateTime` | Время передачи показания, побеждает последняя запись |

#### `fact_invoice`

| Колонка | Тип | Описание |
| --- | --- | --- |
| `realty_id` | `UInt64` | Идентификатор объекта |
| `user_id` | `UInt64` | Идентификатор собственника |
| `tariff_id` | `UInt64` | Идентификатор тарифа |
| `gauge_id` | `UInt64` | Идентификатор счётчика |
| `day_kw` | `UInt32` | Потреблённые кВт·ч дневного тарифа |
| `night_kw` | `UInt32` | Потреблённые кВт·ч ночного тарифа |
| `price` | `UInt64` | Сумма счёта, копейки |
| `deadline` | `Date` | Срок оплаты |
| `created_at` | `DateTime` | Дата выставления счёта |

Ключи сортировки и партиционирование выбираете сами, но обоснуйте выбор в README вашего репозитория.

#### Правила для `VersionedCollapsingMergeTree`

- Установка счётчика: строка с `sign = 1` и `version = v`.
- Аннулирование: строка с `sign = -1`, тем же `version = v` и теми же ключевыми полями, что у аннулируемой строки.
- Новая установка после аннулирования: строка с `sign = 1` и `version = v + 1`.

#### Правила для `ReplacingMergeTree` у показаний

Корректировка показания за дату — новая строка с тем же ключом и более поздним `created_at`.
Старую строку не удаляем.

### Генерация датасетов

Напишите генератор датасетов на любом языке и в любой реализации. Главное, чтобы `make gen` запускался одной командой и без проблем работал в GitHub Actions.
Можно положить в репозиторий готовый бинарник (например, собранный из Go-проекта) или любой другой инструмент, если его запускает одна команда в CI.

Результат генерации должен соответствовать таблице:

| Таблица | Строк | Что должно быть |
| --- | --- | --- |
| `dim_realty` | 1 000 | Квартиры и нежилые помещения, несколько домов |
| `dim_gauges` | 1 000+ | Один счётчик на объект, часть объектов получает замену счётчика |
| `dim_tariffs` | 2+ | Минимум два интервала версий одного тарифа с разными `valid_from` и `valid_to`, открытая версия до `2099-01-01 00:00:00` |
| `dim_users` | 1 000 | Собственник на каждый объект |
| `fact_gauge_bindings` | 1 000+ | Установка счётчика на каждый объект, часть объектов с аннулированием и новой установкой |
| `fact_meter_readings` | ~1 100 000 | Показание на каждый объект за каждый день с 2023-11-01 по 2026-10-31 включительно, около 3% с корректировкой за ту же дату |
| `fact_invoice` | 36 000 | Счёт на каждый объект за каждый месяц периода (36 месяцев), срок оплаты (`deadline`) — 10-е число следующего месяца |

Требования к генератору:

- Результат пишется в папку `data/` в виде CSV.
- Папка `data/` добавлена в `.gitignore`. Датасеты в git не коммитятся.
- Генерация и загрузка в сумме укладываются в 5 минут на стандартном раннере GitHub Actions.

### Загрузка датасетов

Загрузите CSV из `data/` в таблицы ClickHouse. Способ выбирайте сами
(`clickhouse-client`, HTTP-интерфейс на порту 8123, `INSERT ... FORMAT CSV`).

Повторный запуск загрузки без очистки не должен создавать дубли в таблицах с
`ReplacingMergeTree` и `VersionedCollapsingMergeTree`, либо это должно быть явно описано в README.

### HTTP-сервис

Реализуйте HTTP-сервис на любом языке. Сервис работает без авторизации: пользователь
передаётся явно, через путь или тело запроса. Cookie, сессии и Redis не используются.

Общие правила:

- Тела запросов и ответов в формате JSON.
- Ошибки в формате `{"message": "..."}`.
- Поля называются в `snake_case`.
- Даты в ответах: `YYYY-MM-DD` для дат, `YYYY-MM-DD HH:MM:SS` для дат с временем.

#### Проверка доступности

**Запрос:**

```http
GET /ping HTTP/1.1
Host: localhost:8080
```

**Ответ:**

```http
HTTP/1.1 200 OK
Content-Length: 0
```

#### Собственник

**Запрос:**

```http
GET /users/{id} HTTP/1.1
Host: localhost:8080
```

**Ответ (найден):**

```http
HTTP/1.1 200 OK
Content-Type: application/json
```

```json
{
    "id": 1,
    "fio": "Иванов Иван Иванович",
    "phone": "+79000000001",
    "phone_prev": ""
}
```

**Ответ (собственник НЕ найден):**

```http
HTTP/1.1 404 Not Found
Content-Type: application/json
```

```json
{
    "message": "Not found"
}
```

**Запрос на изменение:**

```http
PATCH /users/{id} HTTP/1.1
Host: localhost:8080
Content-Type: application/json
```

```json
{
    "fio": "Петрова Анна Сергеевна",
    "phone": "+79111111111"
}
```

Оба поля необязательны, но хотя бы одно должно быть передано.

- `fio` перезаписывается (SCD 1), старое значение не сохраняется.
- При смене `phone` предыдущий телефон записывается в `phone_prev` (SCD 3).

**Ответ (изменение применено):**

```http
HTTP/1.1 204 No Content
Content-Length: 0
```

**Ответ (невалидные данные):**

```http
HTTP/1.1 400 Bad Request
Content-Type: application/json
```

```json
{
    "message": "invalid \"fio\" field"
}
```

**Ответ (собственник НЕ найден):**

```http
HTTP/1.1 404 Not Found
Content-Type: application/json
```

```json
{
    "message": "Not found"
}
```

#### Объект недвижимости

**Запрос:**

```http
GET /realty/{id} HTTP/1.1
Host: localhost:8080
```

**Ответ (найден):**

```json
{
    "id": 1,
    "address": "г. Москва, улица Пушкина, дом Колотушкина",
    "square": 41,
    "type": "flat",
    "number": "1021",
    "rooms": 2,
    "floor": 10,
    "porch": 1
}
```

**Ответ (объект НЕ найден):**

```http
HTTP/1.1 404 Not Found
Content-Type: application/json
```

```json
{
    "message": "Not found"
}
```

#### Передача показаний

**Запрос:**

```http
POST /realty/{id}/readings HTTP/1.1
Host: localhost:8080
Content-Type: application/json
```

```json
{
    "user_id": 1,
    "gauge_id": 1,
    "tariff_id": 1,
    "reading_date": "2026-10-01",
    "day_kw": 1250,
    "night_kw": 300
}
```

Все поля обязательны. Значения `day_kw` и `night_kw` — целые неотрицательные числа.

> ⚠️ Повторная передача за ту же дату создаёт корректировку. Старая строка остаётся в таблице.

**Ответ (показание сохранено):**

```http
HTTP/1.1 201 Created
Content-Length: 0
```

**Ответ (невалидные поля):**

```http
HTTP/1.1 400 Bad Request
Content-Type: application/json
```

```json
{
    "message": "invalid \"day_kw\" field"
}
```

**Ответ (объект, собственник, счётчик или тариф НЕ найдены):**

```http
HTTP/1.1 404 Not Found
Content-Type: application/json
```

```json
{
    "message": "Not found"
}
```

#### Показания за период

**Запрос:**

```http
GET /realty/{id}/readings?date_from=2026-10-01&date_to=2026-10-31 HTTP/1.1
Host: localhost:8080
```

Границы периода включительны. В ответе для каждой даты возвращается одно актуальное значение
(с учётом корректировок), дубликатов нет.

**Ответ:**

```json
{
    "readings": [
        {
            "reading_date": "2026-10-01",
            "user_id": 1,
            "tariff_id": 1,
            "gauge_id": 1,
            "day_kw": 1250,
            "night_kw": 300
        }
    ],
    "count": 1
}
```

- `count` — количество строк в `readings`.

**Ответ (показания НЕ найдены):**

```http
HTTP/1.1 200 OK
Content-Type: application/json
```

```json
{
    "readings": [],
    "count": 0
}
```

**Ответ (невалидный период):**

```http
HTTP/1.1 400 Bad Request
Content-Type: application/json
```

```json
{
    "message": "invalid \"date_from\" parameter"
}
```

**Ответ (объект НЕ найден):**

```http
HTTP/1.1 404 Not Found
Content-Type: application/json
```

```json
{
    "message": "Not found"
}
```

#### Счета объекта

**Запрос:**

```http
GET /realty/{id}/invoices?limit=10&offset=0 HTTP/1.1
Host: localhost:8080
```

- `limit` — количество счетов на странице, целое число от 0.
- `offset` — сколько счетов пропустить, целое число от 0.
- `limit` и `offset` по умолчанию: `limit = 10`, `offset = 0`.

**Ответ:**

```json
{
    "invoices": [
        {
            "tariff_id": 1,
            "gauge_id": 1,
            "day_kw": 1250,
            "night_kw": 300,
            "price": 93000,
            "deadline": "2026-11-10",
            "created_at": "2026-10-27 10:00:00"
        }
    ],
    "count": 36
}
```

- `count` — общее количество счетов объекта, без учёта `limit` и `offset`.
- Сортировка по `created_at` по убыванию.

**Ответ (невалидные параметры пагинации):**

```http
HTTP/1.1 400 Bad Request
Content-Type: application/json
```

```json
{
    "message": "invalid \"limit\" parameter"
}
```

**Ответ (объект НЕ найден):**

```http
HTTP/1.1 404 Not Found
Content-Type: application/json
```

```json
{
    "message": "Not found"
}
```

#### Тариф на дату

**Запрос:**

```http
GET /tariffs?date=2026-10-27 HTTP/1.1
Host: localhost:8080
```

Возвращается версия тарифа, действовавшая на указанную дату:
`valid_from <= date < valid_to`. Если дата не указана, возвращается `400`.

**Ответ:**

```json
{
    "id": 1,
    "name": "День–Ночь",
    "price_day_kwph": 620,
    "price_night_kwph": 310,
    "valid_from": "2026-07-01 00:00:00",
    "valid_to": "2099-01-01 00:00:00"
}
```

**Ответ (невалидная дата):**

```http
HTTP/1.1 400 Bad Request
Content-Type: application/json
```

```json
{
    "message": "invalid \"date\" parameter"
}
```

**Ответ (на эту дату нет версии тарифа):**

```http
HTTP/1.1 404 Not Found
Content-Type: application/json
```

```json
{
    "message": "Not found"
}
```

#### Установка счётчика на объект

**Запрос:**

```http
POST /realty/{id}/gauges HTTP/1.1
Host: localhost:8080
Content-Type: application/json
```

```json
{
    "gauge_id": 2
}
```

Эндпоинт делает две вещи:

1. Если на объекте уже есть активная установка, аннулирует её строкой с `sign = -1`
   и тем же `version`.
2. Вставляет новую установку счётчика `gauge_id` со `sign = 1` и `version = v + 1`.
   Если активной установки не было, `version = 1`.

**Ответ (установка сохранена):**

```http
HTTP/1.1 201 Created
Content-Length: 0
```

**Ответ (невалидный `gauge_id`):**

```http
HTTP/1.1 400 Bad Request
Content-Type: application/json
```

```json
{
    "message": "invalid \"gauge_id\" field"
}
```

**Ответ (объект или счётчик НЕ найдены):**

```http
HTTP/1.1 404 Not Found
Content-Type: application/json
```

```json
{
    "message": "Not found"
}
```

Формат ответа на аннулирование и активные установки студент проверяет сам через таблицу `fact_gauge_bindings`.

### SCD: сценарии, которые должны работать

Все сценарии проверяются вручную через HTTP-сервис и SQL-запросы.

| Сценарий | Какой SCD | Как проверить |
| --- | --- | --- |
| Собственник сменил фамилию | SCD 1 | `GET /users/{id}`: ФИО новое, старого значения нет нигде |
| Собственник сменил телефон дважды | SCD 3 | `GET /users/{id}`: `phone` — последний, `phone_prev` — предыдущий |
| Изменился тариф с 1 ноября | SCD 2 | `GET /tariffs?date=` до и после 1 ноября возвращает разные версии |
| Пересчёт счёта за старый месяц | SCD 2 | Старый счёт не меняется, а цена берётся по версии на дату счёта |
| Корректировка показания | ReplacingMergeTree | `GET /realty/{id}/readings` возвращает одно значение за дату |
| Замена счётчика | VersionedCollapsingMergeTree | `POST /realty/{id}/gauges` аннулирует старую установку |
| Ошибочная строка с лишним нулём | DELETE | Удалите строку через SQL, покажите, что её больше нет |

Для каждого сценария в README вашего репозитория опишите, какой SQL вы выполнили
и что он показывает.

## Конфигурация

Все параметры задаются в файле `.env.local`. Этот файл коммитится в ваш репозиторий: его читает GitHub Actions при проверке.
Не храните в нём секреты, которые нельзя публиковать.

```sh
# Порт HTTP-сервера, на котором доступен сервис для запросов
APP_PORT=8080

# Параметры подключения к ClickHouse с хоста
CLICKHOUSE_HOST=localhost
CLICKHOUSE_PORT=9000
CLICKHOUSE_HTTP_PORT=8123
CLICKHOUSE_USER=default
CLICKHOUSE_PASSWORD=
CLICKHOUSE_DATABASE=energosbyt
```

⚠️ Имена переменных менять нельзя. Значения можно менять, например, порт `APP_PORT`.

`CLICKHOUSE_HOST=localhost` — адрес ClickHouse с хоста: по нему проверки и `make` подключаются к контейнеру через опубликованные порты.
Контейнер приложения ходит в ClickHouse по имени сервиса, поэтому в `docker-compose.yml` для сервиса приложения задайте `CLICKHOUSE_HOST: clickhouse` в `environment`.

## Запуск

Проект запускается и останавливается командами `make`. Цели обязательны и называются ровно так.

| Команда | Что делает |
| --- | --- |
| `make run` | Поднимает ClickHouse и приложение через `docker compose` |
| `make stop` | Останавливает контейнеры |
| `make gen` | Генерирует датасеты в `data/` |
| `make load` | Загружает датасеты из `data/` в ClickHouse |
| `make setup` | Создаёт схему, генерирует и загружает данные. Требует запущенного ClickHouse (`make run`) |
| `make teardown` | Удаляет таблицы и данные, база возвращается в исходное состояние |

Типичный порядок работы:

```bash
make run
make setup
curl http://localhost:8080/ping
make teardown
make stop
```

После `make run` ClickHouse должен быть доступен на порту `8123`:

```bash
curl http://localhost:8123/ping
# Ожидается: Ok.
```

## Самопроверка

Перед отправкой PR проверьте свою работу: объёмы и сценарии ниже можно проверить вручную.
Количество строк после `make setup` должно соответствовать таблице в разделе «Генерация датасетов».
Для `fact_meter_readings` считайте с `FINAL`.

Дополнительно проверьте:

| Проверка | Ожидаемый результат |
| --- | --- |
| У тарифа есть две версии с разными `valid_from` | `SELECT id FROM dim_tariffs FINAL GROUP BY id HAVING count() >= 2` не пуст |
| В `fact_gauge_bindings` есть строка с `version ≥ 2` | есть (после замены счётчика). Считайте с `FINAL`: строки `sign = -1` при слиянии схлопываются |

Корректировки показаний проверяйте через `FINAL`: до слияния строк за дату может быть две, после — одна.

Пример запроса через `clickhouse-client` внутри контейнера:

```bash
docker compose exec clickhouse clickhouse-client \
  --query "SELECT count() FROM energosbyt.dim_realty"
```

Для `fact_meter_readings` используйте `FINAL`, иначе результат зависит от фонового слияния:

```bash
docker compose exec clickhouse clickhouse-client \
  --query "SELECT count() FROM energosbyt.fact_meter_readings FINAL"
```

### Критерии проверки

- ✅ `make run` поднимает ClickHouse и приложение.
- ✅ `make setup` завершается за 5 минут или меньше.
- ✅ `GET /ping` возвращает `200 OK`.
- ✅ Все таблицы созданы с движками из раздела «Схема данных».
- ✅ Датасеты не коммитятся в git, `data/` и `volume/` в `.gitignore`.
- ✅ Файл `.env.local` коммитится в репозиторий и не игнорируется `.gitignore`.
- ✅ `make teardown` очищает данные, повторный `make setup` работает без ошибок.
- ✅ Имена таблиц, колонок, переменных и make-целей совпадают с этой инструкцией.

## Проверка решения

После завершения работы:

1. **Pull Request** — создайте PR в **свой** публичный репозиторий.
2. **GitHub Actions** — при открытии PR запустится проверка лабы.
3. **Код-ревью** — ревьюер назначится автоматически.

В вашем репозитории должен быть workflow, который вызывает проверку лабы:

```yaml
name: Lab 1

on: pull_request

jobs:
  lab1:
    uses: sitnikovik/abds/.github/workflows/reusable_energosbyt_lab1.yml@main
```

Проверка выполняет по шагам:

1. Проверяет, что в репозитории есть `docker-compose.yml`, `Makefile` и `.env.local`.
2. Запускает `make run`, затем `make setup` с таймаутом 5 минут.
3. Читает переменные из `.env.local` и запускает автограйдер лабы 1 против вашего приложения и ClickHouse.
4. Останавливает контейнеры (`docker compose down -v`).

Подробнее см. в [CONTRIBUTING.md](https://github.com/sitnikovik/abds-template?tab=contributing-ov-file).

## FAQ

**Q: Почему нет авторизации?**  
A: Это учебная работа, фокус на данных и SCD. Пользователь передаётся явно в пути или теле запроса.
**Не используйте такой подход в продакшене**.

**Q: Какой язык выбрать для генератора датасетов?**  
A: Любой. Рекомендуем Python, но это не обязательно.

**Q: Что нужно, чтобы генератор работал в CI?**  
A: Генератор должен запускаться командой make gen на раннере GitHub Actions ubuntu-latest без ручных действий.

- Используйте инструменты, которые уже есть на раннере: `bash`, `make`, `python3` и другие.
- Внешние библиотеки Python (например, `pip install -r requirements.txt`) можно использовать, но установка должна выполняться внутри `make gen`, а не вручную. Учтите, что установка зависимостей замедляет CI, а библиотеки вроде `pandas` могут заметно снизить производительность генерации.
- Можно положить в репозиторий готовый бинарник, если его запускает `make gen`.

**Q: Нужна ли таблица `fact_payments`?**  
A: Нет. Она появится в следующих лабораторных работах.

**Q: Можно ли добавлять колонки в таблицы?**  
A: Да, но не меняйте имена и типы колонок, указанных в этой инструкции.

**Q: Что вернуть в `phone_prev`, если предыдущего телефона не было?**  
A: Пустую строку `""`.

**Q: Почему в `count` для счетов не учитывается `limit`?**  
A: `count` показывает общее количество счетов объекта, чтобы клиент мог строить пагинацию.

**Q: Нужно ли удалять старые показания при корректировке?**  
A: Нет. Корректировка — это новая строка с тем же ключом и более поздним `created_at`.
Старые строки остаются, ClickHouse схлопнет их при слиянии.

**Q: Почему `count()` без `FINAL` даёт другое число?**  
A: `ReplacingMergeTree` и `VersionedCollapsingMergeTree` схлопывают строки в фоне.
До слияния физически строк может быть больше. Используйте `FINAL`, чтобы видеть актуальное состояние.

**Q: Можно ли использовать `OPTIMIZE TABLE ... FINAL` вместо `FINAL` в запросах?**  
A: Можно для проверки, но в сервисе запросы должны работать корректно и без `OPTIMIZE`.

**Q: Можно ли хранить деньги в рублях с дробной частью?**  
A: Нет. Цены и суммы хранятся в копейках целыми числами.

**Q: Можно ли поменять порт сервиса?**  
A: Да, через `APP_PORT` в `.env.local`. Тогда в примерах замените `8080` на ваш порт.

**Q: Что делать, если `make setup` не укладывается в 5 минут?**  
A: Оптимизируйте генерацию (например, генерируйте батчами и пишите CSV потоково) и загрузку (пакетами по 100 000 строк и больше).
