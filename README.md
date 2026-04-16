# VulParser - анализатор безопасности конфигураций

Утилита на Go для анализа конфигурационных файлов веб-приложений (JSON/YAML) и выявления потенциально опасных настроек.

## Возможности

### CLI утилита
- Анализ JSON и YAML конфигураций
- Поиск 5 типов уязвимостей (с возможностью расширения)
- Проверка прав доступа к файлам конфигураций (MEDIUM/HIGH)
- Рекурсивный анализ директорий (`-r`, `--recursive`)
- Цветной вывод с уровнями опасности (LOW/MEDIUM/HIGH)
- Поддержка stdin для конвейерной обработки (`--stdin`)
- Интерактивный режим ввода с подсказкой
- Расширяемая система правил через YAML-конфигурацию (`--rules`)
- Режим silent (не выходить с ошибкой) (`-s`, `--silent`)

### HTTP сервер
- REST API для анализа конфигураций
- Поддержка JSON формата запросов/ответов
- Анализ конфига из тела запроса (`POST /api/analyze`)
- Анализ файла на сервере (`POST /api/analyze/file?path=`)
- Кастомные правила в запросе
- Health check endpoint (`GET /api/health`)
- Контейнеризация (Docker)
- Hot reload для разработки (Air)

### gRPC сервер
- Высокопроизводительный RPC API для анализа конфигураций
- Поддержка protobuf протокола (версия 3)
- Анализ конфига через gRPC метод (`AnalyzerService/Analyze`)
- Анализ файла на сервере (`AnalyzerService/AnalyzeFile`)
- Кастомные правила в виде `google.protobuf.Struct`
- Health check метод (`AnalyzerService/Health`)
- Контейнеризация (Docker) с hot reload

### Обнаруживаемые уязвимости

| Уязвимость | Уровень | Описание |
|------------|---------|----------|
| Debug режим | LOW | Логирование в debug-режиме раскрывает внутренние детали |
| Пароль в открытом виде | HIGH | Пароль или секрет хранится в конфигурации в открытом виде |
| Привязка к 0.0.0.0 | MEDIUM | Приложение доступно со всех сетевых интерфейсов |
| Отключенный TLS | HIGH | Отключена проверка SSL/TLS сертификатов |
| Слабые алгоритмы | HIGH | Использование MD5, SHA1, DES, 3DES, TLSv1.0, RC4 |
| Широкие права доступа | MEDIUM/HIGH | Файл доступен другим пользователям (777, 755, 644) |

### Расширяемость
- Добавление новых правил без перекомпиляции (YAML)
- Поддержка кастомных правил через флаг `--rules`
- Возможность передачи правил в HTTP запросе
- Операторы: `eq`, `contains`, `regex_key`, `regex_value`
- Фильтры: `exclude_value_regex`, `and_value_not_empty`

## Установка и билд CLI-утилиты

```bash
git clone https://github.com/at0mchik/vul-parser
cd vul-parser
go build -o vul_parser ./cmd/parser
```

## Установка, билд и запуск http сервиса

```bash
git clone https://github.com/at0mchik/vul-parser
cd vul-parser
go build -o vul_server ./cmd/server
./vul_server
```

## Установка, билд и запуск grpc сервиса

```bash
git clone https://github.com/at0mchik/vul-parser
cd vul-parser
go build -o vul_grpc ./cmd/grpc
./vul_grpc
```

## Установка, билд и запуск http и grpc сервисов с использованием docker и docker-compose

```bash
git clone https://github.com/at0mchik/vul-parser
cd vul-parser
docker-compose -f docker/docker-compose.yml up -d   
```

## .env файл

Для запуска сервисов необходимо создать .env файл с таким содержимым:

```.env
HTTP_SERVER_PORT=8080
GRPC_SERVER_PORT=9054
```

## Использование CLI-утилиты

```bash
# Анализ файла
./vul_parser config.json

# Анализ YAML файла
./vul_parser config.yaml

# Рекурсивный анализ директории
./vul_parser -r ./configs/
./vul_parser --recursive ./configs/

# Анализ директории без рекурсии (только корень)
./vul_parser ./configs/

# Чтение из stdin
cat config.json | ./vul_parser --stdin
./vul_parser --stdin < config.json

# Интерактивный ввод (Ctrl+D для завершения)
./vul_parser --stdin
Enter configuration (press Ctrl+D on empty line to finish):
{"debug": true, "password": "secret"}
# нажмите Ctrl+D

# Режим silent (не выходить с ошибкой при наличии уязвимостей)
./vul_parser -s config.json
./vul_parser --silent config.json

# Использование кастомных правил
./vul_parser --rules my_rules.yaml config.json

# Комбинирование флагов
./vul_parser -r -s --rules custom.yaml ./configs/
cat config.json | ./vul_parser --stdin -s --rules custom.yaml
```

## Флаги CLI-утилиты
| Флаг | Описание |
| --- | --- |
| `-s`, `--silent` |  Не выходить с ошибкой при наличии уязвимостей |
| `--stdin` |	Читать конфигурацию из STDIN вместо файла |
| `--rules` | Путь к файлу с кастомными правилами |
| `-r`, `--recursive` |	Рекурсивный анализ директории |

## Exit codes CLI-утилиты

| Код | Описание |
| --- | --- |
| `0` | Уязвимости не найдены ИЛИ найдены но включен silent режим |
| `1` | Найдены уязвимости (без silent) ИЛИ ошибка выполнения |

## HTTP API Endpoints

| Метод | Endpoint | Описание | 
| --- | --- | --- | 
| POST | `/api/analyze` | Анализ JSON конфига из тела запроса | 
| POST | `/api/analyze/file?path=` | Анализ файла на сервере | 
| GET | `/api/health` | Проверка состояния сервера | 

Примеры запросов находятся ниже

### gRPC API методы

| Метод | Описание |
|-------|----------|
| `AnalyzerService/Analyze` | Анализ конфига из запроса |
| `AnalyzerService/AnalyzeFile` | Анализ файла на сервере |
| `AnalyzerService/Health` | Проверка состояния сервера |

Примеры запросов находятся ниже

## Обнаруживаемые уязвимости

### 1. Режим отладки (LOW)

- Поля: `debug`, `log.level`, `app.mode`, `verbose`

- Пример: `{"debug": true} или {"log": {"level": "debug"}}`

### 2. Пароль в открытом виде (HIGH)

- Поля с ключами: `password`, `pass`, `pwd`, `secret`, `api_key`, `token`, `jwt_secret`

- Исключаются значения в формате переменных окружения (например `DB_PASSWORD)`

- Пример: `{"database": {"password": "admin123"}}`

### 3. Привязка к 0.0.0.0 (MEDIUM)

- Поля: `host`, `bind`, `listen`, `address` со значением `0.0.0.0`

- Пример: `{"server": {"host": "0.0.0.0"}}`

### 4. Отключенный TLS (HIGH)

- Поля: `insecure_skip_verify`, `verify_ssl`, `disable_verification`, `tls_verify`

- Пример: `{"http": {"tls": {"insecure_skip_verify": true}}}`

### 5. Слабые алгоритмы (HIGH)

- Алгоритмы: `MD5`, `SHA1`, `DES`, `3DES`, `RC4`, `TLSv1.0`, `TLSv1.1`, `RSA-SHA1`

- Пример: `{"security": {"hash_algorithm": "MD5"}}`

### 6. Слишком широкие права доступа (MEDIUM/HIGH)

- HIGH (777, 666): любой пользователь может читать/писать конфиг

- MEDIUM (755, 644): любой пользователь может читать конфиг

- Рекомендация: `chmod 600` для конфигов с секретами, `chmod 640` для остальных

## Структура проекта

```text
vul_parser/
├── cmd/
│   ├── grpc/
│   │   └── main.go                # Точка входа gRPC-сервиса 
│   ├── server/
│   │   └── main.go                # Точка входа HTTP-сервиса                
│   └── parser/
│       └── main.go                # Точка входа CLI-утилиты
├── gen/                           # Сгенерированные файлы по proto-контракту
│   └── proto/
│       └── analyzer/
│           ├── analyzer.pb.go
│           └── analyzer_grpc.pb.go
├── internal/
│   ├── domain/
│   │   ├── dto/
│   │   │   └── dto.go             # Модели для http-сервиса
│   │   └── models/
│   │       ├── rules.go           # Модель правила
│   │       └── vulnerability.go   # Модель уязвимости
│   ├── config/
│   │   └── flags.go               # Парсинг флагов
│   ├── rules/
│   │   ├── loader.go              # Загрузчик правил
│   │   └── builtin.yaml           # Встроенные правила
│   ├── parser/
│   │   └── parser.go              # Парсер JSON/YAML
│   ├── checker/
│   │   └── checker.go             # Движок проверки
│   ├── output/
│   │   └── printer.go             # Вывод результатов
│   ├── permission/
│   │   └── checker.go             # Проверка прав доступа
|   |
│   ├── handler/                   # Слой хендлеров
│   │   ├── analyze_grpc.go
│   │   ├── analyze_http.go
│   │   ├── handler_grpc.go
│   │   └── handler_http.go
│   └── service/                   # Слой сервисов
│       ├── analyze_grpc_service.go
│       ├── analyze_service.go     
│       └── service.go 
├── pkg                            # Допонительные пакеты
|   └── server/                    # Пакет для запуска серверов
│       ├── grpc_server.go
│       └── http_server.go 
|
├── docker 
│   ├── Dockerfile.grpc             # Докерфайл gRPC сервиса
|   ├── Dockerfile                  # Докерфайл http сервиса
|   └── docker-compose.yml          # Докер композ для сервисов
|
├── test-configs                    # Готовые конфиги для тестирования
│   ├── all_vulnerabilites.json     # Конфиг со всеми уязвимостями
│   ├── edge_cases.json             # Конфиг с граничными случаями
|   ├── no_vulnerabilites.json      # Конфиг без уязвимостей
|   └── partial_vulnerabilites.yaml # Конфиг формата YAML
|
├── .air.toml           # Конфигурация .air для hot-reload контейнера http
├── .air.grpc.toml      # Конфигурация .air для hot-reload контейнера gRPC
├── .env                # Переменные окружения
├── go.mod
├── go.sum
└── README.md
```

## Расширение
### Добавление нового правила

Создайте файл `custom_rules.yaml`:

```yaml
rules:
  - id: weak_encryption
    name: Слабое шифрование
    severity: HIGH
    description: Обнаружен слабый метод шифрования
    recommendation: Используйте AES-256
    conditions:
      - path: "$.encryption.algorithm"
        operator: eq
        value: "blowfish"
      - path: "*"
        operator: regex_value
        value: "(?i)(blowfish|twofish)"
```

Запуск:
```bash
./vul_parser --rules custom_rules.yaml config.json
```

### Доступные операторы условий

| Оператор | Описание |
| --- | --- |
|`eq`	| Точное совпадение значения | 
|`contains` | Значение содержит подстроку | 
|`regex_key` | Регулярное выражение по имени поля | 
|`regex_value` | Регулярное выражение по значению |

### Параметры условий

| Параметр | Описание | 
| --- | --- |
| `path` | Путь к полю (JSONPath стиль, * для всех)
| `operator` | Оператор сравнения | 
| `value` | Значение для сравнения | 
| `and_value_not_empty` | Значение не должно быть пустым | 
| `exclude_value_regex` | Исключить значения, подходящие под regex | 

## Примеры запросов для HTTP-сервиса

```bash
# Health check
curl -X GET http://localhost:8080/api/health

# Анализ конфига с уязвимостями
curl -X POST http://localhost:8080/api/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "config": {
      "debug": true,
      "password": "secret"
    }
  }'

# Анализ конфига с кастомными правилами
curl -X POST http://localhost:8080/api/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "config": {
      "debug": true,
      "app": {
        "mode": "production"
      }
    },
    "rules": {
      "rules": [
        {
          "id": "custom_debug_check",
          "name": "Debug mode enabled",
          "severity": "HIGH",
          "description": "Debug mode is enabled",
          "recommendation": "Disable debug mode in production",
          "conditions": [
            {
              "path": "$.debug",
              "operator": "eq",
              "value": true
            }
          ]
        }
      ]
    }
  }'

# Анализ файла на сервере
curl -X POST "http://localhost:8080/api/analyze/file?path=./test-configs/edge_cases.json" \
  -H "Content-Type: application/json" \
  -d '{
    "check_permissions": true
  }'

# Анализ файла с кастомными правилами
curl -X POST "http://localhost:8080/api/analyze/file?path=./test-configs/config.json" \
  -H "Content-Type: application/json" \
  -d '{
    "rules": {
      "rules": [
        {
          "id": "file_permission_check",
          "name": "Strict permissions required",
          "severity": "MEDIUM",
          "description": "File permissions too broad",
          "recommendation": "Run chmod 600",
          "conditions": []
        }
      ]
    },
    "check_permissions": true
  }'
```

## Примеры запросов для gRPC-сервиса

```bash
# Health check
grpcurl -plaintext localhost:9054 analyzer.AnalyzerService/Health

# Анализ конфига с уязвимостями
grpcurl -plaintext -d '{
  "config": {
    "debug": true,
    "password": "secret"
  }
}' localhost:9054 analyzer.AnalyzerService/Analyze

# Анализ конфига с кастомными правилами
grpcurl -plaintext -d '{
  "config": {
    "debug": true,
    "app": {
      "mode": "production"
    }
  },
  "rules": {
    "rules": [
      {
        "id": "custom_debug_check",
        "name": "Debug mode enabled",
        "severity": "HIGH",
        "description": "Debug mode is enabled",
        "recommendation": "Disable debug mode in production",
        "conditions": [
          {
            "path": "$.debug",
            "operator": "eq",
            "value": true
          }
        ]
      }
    ]
  }
}' localhost:9054 analyzer.AnalyzerService/Analyze

# Анализ файла на сервере
grpcurl -plaintext -d '{
  "file_path": "./test-configs/edge_cases.json",
  "check_permissions": true
}' localhost:9054 analyzer.AnalyzerService/AnalyzeFile

# Анализ файла с кастомными правилами
grpcurl -plaintext -d '{
  "file_path": "./test-configs/config.json",
  "rules": {
    "rules": [
      {
        "id": "file_permission_check",
        "name": "Strict permissions required",
        "severity": "MEDIUM",
        "description": "File permissions too broad",
        "recommendation": "Run chmod 600",
        "conditions": []
      }
    ]
  },
  "check_permissions": true
}' localhost:9054 analyzer.AnalyzerService/AnalyzeFile
```

## Примеры CLI-утилиты

### Конфигурация `test-configs/all_vulnerabilities.json`

```json
{
  "debug": true,
  "log": {
    "level": "debug"
  },
  "database": {
    "password": "admin123",
    "host": "localhost"
  },
  "server": {
    "host": "0.0.0.0",
    "port": 3000
  },
  "http": {
    "tls": {
      "insecure_skip_verify": true
    }
  },
  "security": {
    "hash_algorithm": "MD5",
    "cipher": "DES"
  },
  "api": {
    "secret_key": "sk_test_12345"
  }
}
```

Запуск и вывод:

```bash
$ ./vul_parser test-configs/all_vulnerabilities.json                                                                          

LOW: Логирование в режиме отладки раскрывает детали внутреннего устройства приложения
  Location: debug
  Value: true
  Recommendation: Смените уровень логирования на info или error в production-среде

LOW: Логирование в режиме отладки раскрывает детали внутреннего устройства приложения
  Location: log.level
  Value: debug
  Recommendation: Смените уровень логирования на info или error в production-среде

MEDIUM: Приложение доступно со всех сетевых интерфейсов
  Location: server.host
  Value: 0.0.0.0
  Recommendation: Используйте 127.0.0.1 для локального доступа или добавьте белый список IP и firewall

HIGH: Отключена проверка SSL/TLS сертификатов, уязвимость для MITM-атак
  Location: http.tls.insecure_skip_verify
  Value: true
  Recommendation: Включите проверку TLS, установите insecure_skip_verify: false или verify_ssl: true

HIGH: Обнаружен устаревший или небезопасный алгоритм
  Location: security
  Value: map[cipher:DES hash_algorithm:MD5]
  Recommendation: Используйте SHA-256 или SHA-512 для хеширования, AES-256 для шифрования, TLS 1.2 или выше

HIGH: Обнаружен устаревший или небезопасный алгоритм
  Location: security.hash_algorithm
  Value: MD5
  Recommendation: Используйте SHA-256 или SHA-512 для хеширования, AES-256 для шифрования, TLS 1.2 или выше

HIGH: Обнаружен устаревший или небезопасный алгоритм
  Location: security.cipher
  Value: DES
  Recommendation: Используйте SHA-256 или SHA-512 для хеширования, AES-256 для шифрования, TLS 1.2 или выше
```

### Конфигурация `test-configs/no_vulnerabilities.json`

```json
{
  "debug": false,
  "log": {
    "level": "info"
  },
  "database": {
    "password_env": "DB_PASSWORD",
    "host": "127.0.0.1"
  },
  "server": {
    "host": "127.0.0.1",
    "port": 3000
  },
  "http": {
    "tls": {
      "insecure_skip_verify": false
    }
  },
  "security": {
    "hash_algorithm": "SHA-256",
    "cipher": "AES-256"
  }
}
```

Запуск и вывод:

```bash
$ ./vul_parser test-configs/no_vulnerabilities.json                                                                                    
No vulnerabilities found
```

### Конфигурация `test-configs/partial_vulnerabilities.yaml`

```yaml
version: 2.0
app:
  mode: debug
  verbose: true

storage:
  type: postgres
  password: postgres123
  host: 0.0.0.0

client:
  verify_ssl: false

tls_config:
  min_version: TLSv1.0
  cipher_suites:
    - TLS_RSA_WITH_3DES_EDE_CBC_SHA
```

Запуск и вывод:

```bash
$ ./vul_parser test-configs/partial_vulnerabilities.yaml                                                                       ✔ 
LOW: Логирование в режиме отладки раскрывает детали внутреннего устройства приложения
  Location: app.mode
  Value: debug
  Recommendation: Смените уровень логирования на info или error в production-среде

MEDIUM: Приложение доступно со всех сетевых интерфейсов
  Location: storage.host
  Value: 0.0.0.0
  Recommendation: Используйте 127.0.0.1 для локального доступа или добавьте белый список IP и firewall

HIGH: Отключена проверка SSL/TLS сертификатов, уязвимость для MITM-атак
  Location: client.verify_ssl
  Value: false
  Recommendation: Включите проверку TLS, установите insecure_skip_verify: false или verify_ssl: true

HIGH: Обнаружен устаревший или небезопасный алгоритм
  Location: tls_config
  Value: map[cipher_suites:[TLS_RSA_WITH_3DES_EDE_CBC_SHA] min_version:TLSv1.0]
  Recommendation: Используйте SHA-256 или SHA-512 для хеширования, AES-256 для шифрования, TLS 1.2 или выше

HIGH: Обнаружен устаревший или небезопасный алгоритм
  Location: tls_config.min_version
  Value: TLSv1.0
  Recommendation: Используйте SHA-256 или SHA-512 для хеширования, AES-256 для шифрования, TLS 1.2 или выше

HIGH: Обнаружен устаревший или небезопасный алгоритм
  Location: tls_config.cipher_suites
  Value: [TLS_RSA_WITH_3DES_EDE_CBC_SHA]
  Recommendation: Используйте SHA-256 или SHA-512 для хеширования, AES-256 для шифрования, TLS 1.2 или выше

HIGH: Обнаружен устаревший или небезопасный алгоритм
  Location: tls_config.cipher_suites[0]
  Value: TLS_RSA_WITH_3DES_EDE_CBC_SHA
  Recommendation: Используйте SHA-256 или SHA-512 для хеширования, AES-256 для шифрования, TLS 1.2 или выше

```