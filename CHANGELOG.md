# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/2.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Fourth practice lesson:
  - ClickHouse cluster of 3 shards × 2 replicas coordinated by ClickHouse Keeper;
  - Monitoring with Prometheus and Grafana dashboards;
  - Python data generators for the energosbyt dims and facts;
  - DDL for distributed and replicated tables, DML load scripts
and demo SQL on replication, failover and sharding;

## [2.1.0] - 2026-11-01

### Added

- Autograder for checking student labs in GitHub Actions: step-based packages under `internal/step/...`,
ClickHouse repo layer and checks with `buggles`
- Check for Lab 1
- Coverage checking scripts and autograder CI workflow
- Code style guide `STYLE.md`
- `AGENTS.md` guidance for AI agents working in the repository

## [2.0.0] - 2026-10-29

### Added

- Docker Compose setup for `prac/02`
- Third practice lesson in `prac/03`: data generators, ClickHouse DDL and DML scripts, Makefile
- GitHub Actions workflows for Go linting, tests and releases
- Pull request template

### Changed

- Ignore datasets, temporary and harness directories in `.gitignore`

### Removed

- **BREAKING:** version 1 directories and files

## [1.6.0] - 2025-12-13

### Added

- Airflow setup and week 7 workshop materials

### Changed

- Updated Makefile, Docker Compose and README for week 7

## [1.5.0] - 2025-12-04

### Added

- README for term papers and README for the antifraud, clickstream and logging topics

### Changed

- Updated main README

## [1.4.0] - 2025-11-29

### Added

- Week 6 materials: extended Go producer services, repositories and common utilities
- Superset build

### Changed

- Reworked producer domain and infrastructure, updated Docker Compose and Makefile

## [1.3.0] - 2025-11-15

### Added

- Week 5 materials: Go Kafka producer for energosbyt gauges and its configuration
- Energosbyt domain models for flats and gauges

### Changed

- Updated Makefile and Docker Compose

## [1.2.1] - 2025-10-18

### Fixed

- ClickHouse configuration

## [1.2.0] - 2025-10-18

### Added

- Week 4 workshop materials
- ClickHouse configuration directory

### Changed

- Updated Docker Compose and README

## [1.1.0] - 2025-10-04

### Added

- Workshop 3 (SCD) tasks for groups 1 to 4

### Changed

- Updated Makefile, Docker Compose, `.gitignore` and README

## [1.0.0] - 2025-09-20

### Added

- Workshop materials for the ClickHouse lesson: Makefile, Docker Compose, ClickHouse in Docker, README and SQL data generator
