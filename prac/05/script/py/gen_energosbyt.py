#!/usr/bin/env python3
# Генерирует CSV с большими данными модели «ЭнергоДом» (prac/03/data/*.dataset.csv).
# Зависит от data/dim_users.dataset.csv: person_no становится user_id.
#
# Что генерируем:
#   dim_realty           — квартиры и нежилые помещения трёх домов из README;
#   dim_gauges           — один счётчик на объект и замены счётчиков;
#   dim_tariffs          — версии тарифов, цена меняется раз в полгода (SCD 2);
#   fact_gauge_bindings  — установки и аннулирования (sign);
#   fact_meter_readings  — показания за каждый месяц, часть с корректировками;
#   fact_invoice         — счёт за каждый месяц по версии тарифа на дату счёта;
#   fact_payments        — оплаты счетов.
#
# Использование: python3 script/py/gen_energosbyt.py
import csv
import random
from datetime import datetime, timedelta
from pathlib import Path

SEED = 42
DATA_DIR = Path(__file__).resolve().parents[2] / "data"
OPEN_VALID_TO = datetime(2099, 1, 1)
# «Сейчас» для датасета: платежи позже этой даты не создаём.
NOW = datetime(2026, 10, 23)

FIRST_MONTH = (2019, 1)
LAST_MONTH = (2026, 9)

# Версии тарифов: раз в полгода с 1 января 2019 по 1 июля 2026.
# Конечные цены совпадают с README: День–Ночь 620 / 310, Базовый 720.
STEPS = [datetime(2019 + i // 2, 1 if i % 2 == 0 else 7, 1) for i in range(16)]
TARIFFS = {
    1: ("День–Ночь", 450, 620),
    2: ("Базовый", 520, 720),
}

GAUGE_MODELS = [
    ("Инкотекс Меркурий 230", "Россия"),
    ("Тайпит НЕВА МТ 324", "Россия"),
    ("Tuya ZMAi-90", "China"),
]

# Дома из README. Первый этаж нежилой, квартиры начинаются со второго.
HOUSES = [
    {"address": "г. Москва, улица Пушкина, дом Колотушкина", "porches": 4, "floors": 20,
     "flats": [40, 40, 65, 90], "commercial": [100, 50, 40]},
    {"address": "г. Москва, проспект Мира, дом Кефира", "porches": 2, "floors": 14,
     "flats": [33, 33, 55, 72], "commercial": [60, 70]},
    {"address": "г. Москва, улица Есенина, дом Каруселина", "porches": 7, "floors": 5,
     "flats": [45, 45, 68, 83, 102], "commercial": [500]},
]


def fmt(dt: datetime) -> str:
    return dt.strftime("%Y-%m-%d %H:%M:%S")


def months():
    y, m = FIRST_MONTH
    while (y, m) <= LAST_MONTH:
        yield y, m
        m += 1
        if m == 13:
            y, m = y + 1, 1


def step_index(dt: datetime) -> int:
    index = 0
    for i, start in enumerate(STEPS):
        if start <= dt:
            index = i
    return index


def tariff_price(tariff_id: int, index: int):
    _, d0, d1 = TARIFFS[tariff_id]
    day = round(d0 + (d1 - d0) * index / (len(STEPS) - 1))
    night = round(day / 2) if tariff_id == 1 else 0
    return day, night


def build_realties():
    rows = []
    for h in HOUSES:
        for sq in h["commercial"]:
            rows.append({"address": h["address"], "square": sq, "type": "nrcp", "number": "",
                         "floor": 1, "porch": 0, "commercial": True})
        for porch in range(1, h["porches"] + 1):
            for floor in range(2, h["floors"] + 1):
                for idx, sq in enumerate(h["flats"]):
                    rows.append({"address": h["address"], "square": sq, "type": "flat",
                                 "number": f"{porch}{floor:02d}{idx + 1}",
                                 "floor": floor, "porch": porch, "commercial": False})
    for rid, row in enumerate(rows, start=1):
        row["id"] = rid
    return rows


def build_tariffs():
    rows = []
    for tariff_id, (name, _, _) in TARIFFS.items():
        for i, start in enumerate(STEPS):
            end = STEPS[i + 1] if i + 1 < len(STEPS) else OPEN_VALID_TO
            day, night = tariff_price(tariff_id, i)
            rows.append([tariff_id, name, day, night, fmt(start), fmt(end), fmt(start), fmt(start)])
    return rows


def consumption(square: int, commercial: bool, vacant: bool):
    base = square * (6.0 if commercial else 3.5)
    total = base * random.uniform(0.7, 1.3) * (0.15 if vacant else 1.0)
    return round(total * 0.6), round(total * 0.4)


def read_user_ids():
    with open(DATA_DIR / "dim_users.dataset.csv", encoding="utf-8") as f:
        return [int(row["person_no"]) for row in csv.DictReader(f)]


def write_csv(name: str, header, rows):
    with open(DATA_DIR / name, "w", encoding="utf-8", newline="") as f:
        out = csv.writer(f, lineterminator="\n")
        out.writerow(header)
        out.writerows(rows)
    print(f"{name}: {len(rows)} строк")


def main() -> None:
    random.seed(SEED)
    user_ids = read_user_ids()
    realties = build_realties()
    repl_months = [(y, m) for (y, m) in months() if (2020, 1) <= (y, m) <= (2025, 12)]

    gauges, bindings, readings, invoices, payments = [], [], [], [], []
    next_gauge_id = len(realties) + 1
    month_list = list(months())

    for r in realties:
        rid = r["id"]
        owner = random.choice(user_ids)
        tariff_id = 1 if r["commercial"] or random.random() < 0.8 else 2

        model, country = random.choice(GAUGE_MODELS)
        gauges.append([rid, model, country, "2019-01-01 00:00:00"])
        bindings.append([rid, rid, "2019-01-01 00:00:00", 1])

        # Замена счётчика: старый аннулируем строкой со sign = -1, новый ставим.
        replacement = None
        if random.random() < 0.10:
            y, m = random.choice(repl_months)
            change = datetime(y, m, 1, 10, 0, 0)
            new_gid = next_gauge_id
            next_gauge_id += 1
            model2, country2 = random.choice(GAUGE_MODELS)
            gauges.append([new_gid, model2, country2, fmt(change)])
            bindings.append([rid, rid, "2019-01-01 00:00:00", -1])
            bindings.append([new_gid, rid, fmt(change), 1])
            replacement = ((change.year, change.month), new_gid)

        day_cum = random.randint(2000, 6000)
        night_cum = random.randint(1000, 3000)

        for y, m in month_list:
            if replacement and (y, m) >= replacement[0]:
                gauge_id = replacement[1]
            else:
                gauge_id = rid
            vacant = random.random() < 0.08
            day_c, night_c = consumption(r["square"], r["commercial"], vacant)
            day_cum += day_c
            night_cum += night_c
            reading_month = f"{y}-{m:02d}-01"

            # Корректировка: сначала ошибочное показание (15-го), потом верное (19-го).
            if random.random() < 0.03:
                wrong = datetime(y, m, 15, 10, 0, 0)
                fixed = datetime(y, m, 19, 12, 0, 0)
                readings.append([rid, owner, tariff_id, gauge_id, reading_month,
                                 day_cum + random.randint(30, 200), night_cum, fmt(wrong)])
                readings.append([rid, owner, tariff_id, gauge_id, reading_month,
                                 day_cum, night_cum, fmt(fixed)])
            else:
                created = datetime(y, m, random.randint(15, 20), random.randint(8, 22),
                                   random.randint(0, 59), random.randint(0, 59))
                readings.append([rid, owner, tariff_id, gauge_id, reading_month,
                                 day_cum, night_cum, fmt(created)])

            # Счёт выставляем 27-го того же месяца, срок оплаты — 10-е число следующего месяца.
            issued = datetime(y, m, 27, 10, 0, 0)
            ny, nm = (y + 1, 1) if m == 12 else (y, m + 1)
            deadline = datetime(ny, nm, 10)
            index = step_index(issued)
            day_price, night_price = tariff_price(tariff_id, index)
            if tariff_id == 2:
                price = (day_c + night_c) * day_price
            else:
                price = day_c * day_price + night_c * night_price
            invoices.append([rid, owner, tariff_id, gauge_id, day_c, night_c, price,
                             deadline.strftime("%Y-%m-%d"), fmt(issued)])

            roll = random.random()
            if roll < 0.80:
                paid = issued + timedelta(days=random.randint(1, 13), hours=random.randint(8, 20))
            elif roll < 0.92:
                paid = deadline + timedelta(days=random.randint(1, 30), hours=12)
            else:
                paid = None
            if paid and paid <= NOW:
                payments.append([rid, owner, price, fmt(paid)])

    write_csv("dim_realty.dataset.csv", ["id", "address", "square", "type", "number", "floor", "porch", "created_at"],
              [[r["id"], r["address"], r["square"], r["type"], r["number"], r["floor"], r["porch"],
                "2018-12-20 10:00:00"] for r in realties])
    write_csv("dim_gauges.dataset.csv", ["id", "name", "country", "created_at"], gauges)
    write_csv("dim_tariffs.dataset.csv",
              ["id", "name", "price_day_kwph", "price_night_kwph", "valid_from", "valid_to", "created_at", "updated_at"],
              build_tariffs())
    write_csv("fact_gauge_bindings.dataset.csv", ["gauge_id", "realty_id", "created_at", "sign"], bindings)
    write_csv("fact_meter_readings.dataset.csv",
              ["realty_id", "user_id", "tariff_id", "gauge_id", "reading_month", "day_kw", "night_kw", "created_at"],
              readings)
    write_csv("fact_invoice.dataset.csv",
              ["realty_id", "user_id", "tariff_id", "gauge_id", "day_kw", "night_kw", "price", "deadline", "created_at"],
              invoices)
    write_csv("fact_payments.dataset.csv", ["realty_id", "user_id", "price", "created_at"], payments)


if __name__ == "__main__":
    main()
