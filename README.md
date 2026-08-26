# Композиционный анализ ПО: Go + CycloneDX + Dependency-Track

Учебная работа по дисциплине «Организационно-правовые основы разработки программного обеспечения».

## Цель

На практике показать процесс Software Composition Analysis (SCA):

1. развернуть OWASP Dependency-Track;
2. создать небольшое приложение на Go с внешними зависимостями;
3. сформировать SBOM в формате CycloneDX;
4. загрузить SBOM в Dependency-Track и посмотреть состав компонентов и найденные риски.

## Структура

```text
composition-analysis/
├── dtrack/
│   └── docker-compose.yml
├── app/
│   ├── cmd/
│   │   └── hello/
│   │       └── main.go
│   ├── go.mod
│   ├── go.sum
│   └── bom.json
└── README.md
```

## 1. Dependency-Track

В работе используется ветка Dependency-Track 4.14.x и PostgreSQL 15. API-сервер и frontend запускаются отдельными контейнерами, а PostgreSQL хранит данные в отдельном volume.

Перейти в каталог `dtrack` и выполнить:

```bash
docker compose pull
docker compose up -d
docker compose ps
```

Веб-интерфейс:

```text
http://localhost:8080
```

API-сервер доступен на:

```text
http://localhost:8081
```

При первом запуске Dependency-Track выполняет начальную инициализацию и загрузку данных о лицензиях и уязвимостях. Это может занять 10–30 минут и более. Не следует останавливать контейнеры во время этой операции.

Для просмотра логов:

```bash
docker compose logs -f dtrack-apiserver
```

Первичные учетные данные:

```text
Логин:    admin
Пароль:   admin
```

После первого входа пароль рекомендуется изменить.

Остановить окружение:

```bash
docker compose down
```

Удалить также данные PostgreSQL и Dependency-Track:

```bash
docker compose down -v
```

## 2. Go-приложение

Перейти в каталог приложения:

```bash
cd app
```

Проверить зависимости:

```bash
go mod tidy
```

Запустить:

```bash
go run ./cmd/hello
```

Пример результата:

```text
Hello, World!
Version: 1.0.0
Run ID: <UUID>
```

Проверить версию:

```bash
go run ./cmd/hello --version
```

В приложении намеренно используются внешние библиотеки:

- `github.com/google/uuid` — генерация UUID;
- `github.com/rs/zerolog` — структурированное логирование;
- `github.com/spf13/cobra` — обработка CLI-команд и флагов.

Это позволяет получить содержательный SBOM для демонстрации композиционного анализа.

## 3. Генерация SBOM

Для Go-модулей удобно использовать официальный инструмент CycloneDX `cyclonedx-gomod`.

Установка:

```bash
go install github.com/CycloneDX/cyclonedx-gomod/cmd/cyclonedx-gomod@latest
```

Убедиться, что инструмент доступен:

```bash
cyclonedx-gomod version
```

Из каталога `app` сформировать SBOM приложения:

```bash
cyclonedx-gomod app -main ./cmd/hello -json -output bom.json .
```

Для более подробного BOM, например с лицензиями и пакетами:

```bash
cyclonedx-gomod app -main ./cmd/hello -licenses -packages -json -output bom.json .
```

Полученный файл:

```text
app/bom.json
```

имеет формат CycloneDX JSON.

Для учебной работы также можно сформировать BOM всего Go-модуля:

```bash
cyclonedx-gomod mod -licenses -json -output bom-module.json .
```

Разница:

- `app` показывает зависимости конкретного приложения с учетом импортируемых пакетов и build constraints;
- `mod` показывает агрегированный набор зависимостей модуля.

## 4. Загрузка SBOM в Dependency-Track

После входа в веб-интерфейс:

1. Создать проект, например `composition-analysis`.
2. Задать версию проекта, например `1.0.0`.
3. Открыть страницу проекта.
4. Выбрать импорт SBOM.
5. Загрузить `app/bom.json`.
6. Дождаться обработки BOM.
7. Изучить список Components и найденные Vulnerabilities.

В Dependency-Track можно увидеть, какие сторонние компоненты входят в приложение, какие версии используются, какие лицензии обнаружены и какие известные уязвимости связаны с компонентами.

## 5. Что демонстрирует работа

Получается следующий конвейер:

```text
Go source code
      |
      v
go.mod / go.sum
      |
      v
cyclonedx-gomod
      |
      v
CycloneDX SBOM (bom.json)
      |
      v
Dependency-Track
      |
      +--> Components
      +--> Licenses
      +--> Vulnerabilities
      +--> Risk analysis
```

SBOM является описанием состава программного продукта. Dependency-Track использует его для анализа компонентов и связанных с ними рисков.

## 6. Полезные команды

Проверка контейнеров:

```bash
docker compose ps
```

Логи API:

```bash
docker compose logs -f dtrack-apiserver
```

Логи frontend:

```bash
docker compose logs -f dtrack-frontend
```

Проверка приложения:

```bash
cd app
go run ./cmd/hello
go run ./cmd/hello --version
```

Сборка:

```bash
go build -o hello ./cmd/hello
./hello
```

## 7. Использованные ресурсы

- OWASP Dependency-Track — документация по Docker-развертыванию: https://docs.dependencytrack.org/getting-started/deploy-docker/
- OWASP Dependency-Track — первичный запуск: https://docs.dependencytrack.org/getting-started/initial-startup/
- CycloneDX — официальный каталог инструментов: https://cyclonedx.org/tool-center/
- CycloneDX Go Modules — https://github.com/CycloneDX/cyclonedx-gomod
- CycloneDX Specification — https://cyclonedx.org/specification/overview/

## 8. Примечание по версиям

В учебной работе зафиксирована версия Dependency-Track 4.14.2, чтобы пример был воспроизводимым. В дальнейшем Dependency-Track 5.x имеет отдельную архитектуру и процедуру миграции с v4, поэтому для учебного стенда не следует бездумно заменять образы v4 на v5.

Пароли и настройки из данного `docker-compose.yml` предназначены только для локального учебного стенда и не должны использоваться в production.
