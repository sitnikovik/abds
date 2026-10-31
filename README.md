# Archtecture of Big Data Systems

[![Quality](https://github.com/sitnikovik/abds/actions/workflows/quality.yml/badge.svg)](https://github.com/sitnikovik/abds/actions/workflows/quality.yml)
[![Tests](https://github.com/sitnikovik/abds/actions/workflows/tests.yml/badge.svg)](https://github.com/sitnikovik/abds/actions/workflows/tests.yml)
![Lab Progress](https://img.shields.io/badge/Lab-3%20of%207-blue)

Проект содержит все необходимые материалы для курса по NoSQL базам данных
и пайплайн в GitHub Actions для проверки лабораторных работ.

## Семинары

- [Знакомство c ClickHouse](/prac/02/)
- [SCD - медленно изменяющиеся измерения](/prac/02/)

## Лабораторные работы

**ЭнергоСбыт** — backend‑сервис вымышленной ЖКХ-компании, предназначенный
для изучения архитектутрных принципов систем обработки больших данных.
Проект выполняется поэтапно в рамках лабораторных работ и семинаров
и представляет собой единый сервис, который развивается в течение всего курса.

> Задания публикуются постепенно по мере прохождения курса

- [1. ЭнергоСбыт. Приём показаний электричества и выставление счетов](/lab/01/)

## Требования к реализации

- Используйте предоставленный шаблон репозитория
[abds-template](https://github.com/sitnikovik/abds-template)
- Соблюдайте требования к структуре проекта, описанные в
[CONTRIBUTING.md](https://github.com/sitnikovik/abds-template?tab=contributing-ov-file)
- Можно использовать **любой язык программирования**
- Фреймворки и библиотеки - на ваше усмотрение, но не забывайте о производительности
- Все лабораторные работы должны полностью проходить
автоматическую проверку в GitHub Actions
- Каждая последующая лабораторная работа должна
на основе предыдущей
- Очередная лабораторная работа не должна ломать
функциональность предыдущих лабораторных работ
(проверяется в процессе автоматической проверки)

## Система оценивания

За выполнение лабораторных работ вы получаете баллы,
из которых складывается итоговая оценка за семестр.

У каждой лабораторной работы есть своя максимальная сумма баллов и **срок сдачи**.
Если вы сдаёте работу после срока,
то из полученных баллов вычитается "штраф за просрочку".

Сроки и баллы указаны в таблице успеваемости (уточняйте у преподавателя).

## Рекомендации

- **Документируйте** исходный код, соблюдая стандарты, например:
  - [Javadoc](https://docs.oracle.com/en/java/javase/17/docs/specs/javadoc/doc-comment-spec.html)
  - [PEP8](https://peps.python.org/pep-0008/)
  - [Go doc](https://go.dev/doc/comment)
  - [PHPDoc](https://github.com/php-fig/fig-standards/blob/master/proposed/phpdoc.md)
  - [JSDoc](https://jsdoc.app/)
- Собюдайте [конвенцию коммитов](https://www.conventionalcommits.org/ru/v1.0.0-beta.2/)
- Каждый коммит - рабочая и протестированная версия проекта. Это облегчит поиск последней рабочей версии при поиске ошибки
- Пишите документацию на [Markdown](https://www.markdownguide.org/) в [Obsidian](https://obsidian.md/) или в вашей IDE.
- Используйте [Makefile](https://makefiletutorial.com/) для удобного взаимодействия с проектом
- Используйте [PlantUML](https://plantuml.com/ru/) для создания блок-схемы и диаграмм.
Есть возможность предпросмотра в IDE
(в VSCode через [расширение PlantUML](https://marketplace.visualstudio.com/items?itemName=jebbs.plantuml))
и утилиты graphviz и Java, установленных на вашем компьютере ([подробнее](https://plantuml.com/ru-dark/faq-install))
- Схемы БД реализуйте через [DBML](https://dbml.dbdiagram.io/home/)
и [dbdiagram.io](https://dbdiagram.io/home).
В VSCode есть [расширение](https://marketplace.visualstudio.com/items?itemName=dbdiagram.dbdiagram-vscode) для предпросмотра
- Ведите коллекцию [Postman](https://www.postman.com/)
или [Insomnia](https://insomnia.rest/) c примерами запросов и ответов для **каждого эндпоинта**.
Можно развернуть [Swagger](https://swagger.io/), но не обязательно.
- Используйте публичные **датасеты** с тестовыми данными,
например, из [Kaggle](https://www.kaggle.com/) или генерируйте свои на Python, Go или любом другом языке программирования.
- Если нужен текст-рыба - используйте
  - [генератор текста "рыба"](https://fish-text.ru/) если нужен текст на русском
  - или [Lorem Ipsum](https://getlorem.com/) если нужен текст на английском
