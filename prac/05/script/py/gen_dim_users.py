#!/usr/bin/env python3
# Генерирует CSV-датасет для dim_users (prac/03/script/sql/dim_users.table.sql).
#
# Таблица хранит текущее состояние собственника:
#   fio — SCD 1: исправления перезаписываются, история не хранится;
#   phone — SCD 3: текущий телефон и предыдущий в phone_prev.
#
# person_no в таблицу не попадает. Это номер собственника в датасете.
# Загрузка (script/sql/load_dim_users.sql) использует его как id.
#
# Предписание (README.md, раздел «Собственники»):
#   - один ответственный квартиросъёмщик на квартиру;
#   - заселённость объектов 70–100%;
#   - квартира может сменить жильца (нужны «запасные» пользователи).
#
# Всего квартир: 320 + 112 + 175 = 607. Генерируем 700 собственников
# (607 максимальная заселённость + ~15% запаса на смену жильцов).
#
# Использование: python3 script/py/gen_dim_users.py > data/dim_users.dataset.csv

import csv
import random
import sys
from datetime import datetime, timedelta

SEED = 42
TOTAL_USERS = 700
PHONE_CHANGE_RATIO = 0.3  # доля собственников, сменивших телефон (SCD 3: phone_prev)

CREATED_FROM = datetime(2019, 1, 1)
CREATED_TO = datetime(2025, 8, 31)
# «Сейчас» для датасета: события не могут быть позже этой даты.
NOW = datetime(2026, 10, 1)

SURNAMES = [
    "Смирнов", "Иванов", "Кузнецов", "Соколов", "Попов", "Лебедев",
    "Козлов", "Новиков", "Морозов", "Волков", "Петров", "Егоров",
    "Павлов", "Семёнов", "Голубев", "Виноградов", "Богданов", "Воробьёв",
    "Фёдоров", "Михайлов", "Беляев", "Тарасов", "Белов", "Комаров",
    "Орлов", "Киселёв", "Макаров", "Андреев", "Ковалёв", "Ильин",
]
MALE_NAMES = [
    "Александр", "Сергей", "Дмитрий", "Андрей", "Максим", "Иван",
    "Артём", "Николай", "Владимир", "Кирилл", "Роман", "Егор",
]
FEMALE_NAMES = [
    "Елена", "Ольга", "Мария", "Наталья", "Ирина", "Татьяна",
    "Светлана", "Юлия", "Анна", "Екатерина", "Дарья", "Анастасия",
]
MALE_PATRONYMICS = [
    "Александрович", "Сергеевич", "Дмитриевич", "Андреевич", "Максимович",
    "Иванович", "Артёмович", "Николаевич", "Владимирович", "Кириллович",
]
FEMALE_PATRONYMICS = [
    "Александровна", "Сергеевна", "Дмитриевна", "Андреевна", "Максимовна",
    "Ивановна", "Артёмовна", "Николаевна", "Владимировна", "Кирилловна",
]


def fmt(dt: datetime) -> str:
    return dt.strftime("%Y-%m-%d %H:%M:%S")


def gen_phone(used: set) -> str:
    while True:
        phone = "+79" + "".join(random.choices("0123456789", k=9))
        if phone not in used:
            used.add(phone)
            return phone


def main() -> None:
    random.seed(SEED)
    out = csv.writer(sys.stdout, lineterminator="\n")
    out.writerow(["person_no", "fio", "phone", "phone_prev", "created_at", "updated_at"])

    used_phones: set = set()
    used_fio: set = set()
    rows = []

    for person_no in range(1, TOTAL_USERS + 1):
        # ФИО без повторов: комбинируем, пока не получим уникальное.
        # Женская фамилия образуется от мужской добавлением «а».
        while True:
            female = random.random() < 0.5
            surname = random.choice(SURNAMES)
            fio = " ".join([
                surname + "а" if female else surname,
                random.choice(FEMALE_NAMES if female else MALE_NAMES),
                random.choice(
                    FEMALE_PATRONYMICS if female else MALE_PATRONYMICS
                ),
            ])
            if fio not in used_fio:
                used_fio.add(fio)
                break

        created_at = CREATED_FROM + timedelta(
            seconds=random.randint(
                0, int((CREATED_TO - CREATED_FROM).total_seconds())
            )
        )
        created_at = created_at.replace(
            hour=random.randint(8, 22),
            minute=random.randint(0, 59),
            second=random.randint(0, 59),
        )

        phone = gen_phone(used_phones)
        phone_prev = ""
        updated_at = created_at

        # Смена телефона: прежний номер уходит в phone_prev, updated_at фиксирует момент смены.
        if random.random() < PHONE_CHANGE_RATIO:
            change_at = created_at + timedelta(
                days=random.randint(30, 900),
                hours=random.randint(0, 23),
                minutes=random.randint(0, 59),
                seconds=random.randint(0, 59),
            )
            if change_at < NOW:
                phone_prev = phone
                phone = gen_phone(used_phones)
                updated_at = change_at

        rows.append((created_at, person_no, [
            person_no,
            fio,
            phone,
            phone_prev,
            fmt(created_at),
            fmt(updated_at),
        ]))

    # Порядок по дате регистрации — как будто пользователи регистрировались хронологически.
    rows.sort(key=lambda row: (row[0], row[1]))
    out.writerows(row[2] for row in rows)


if __name__ == "__main__":
    main()
