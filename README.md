# VulParser - анализатор безопасности конфигураций

Утилита на Go для анализа конфигурационных файлов веб-приложений (JSON/YAML) и выявления потенциально опасных настроек.

## Возможности

- Анализ JSON и YAML конфигураций
- Поиск 5 типов уязвимостей (с возможностью расширения)
- Проверка прав доступа к файлам конфигураций
- Рекурсивный анализ директорий
- Цветной вывод с уровнями опасности (LOW/MEDIUM/HIGH)
- Поддержка stdin для конвейерной обработки
- Расширяемая система правил через YAML-конфигурацию
- Режим silent (не выходить с ошибкой)

## Установка

```bash
git clone https://github.com/at0mchik/vul-parser
cd vul-parser
go build -o vul_parser ./cmd/app
```

## Использование

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

## Флаги
| Флаг | Описание |
| --- | --- |
| `-s`, `--silent` |  Не выходить с ошибкой при наличии уязвимостей |
| `--stdin` |	Читать конфигурацию из STDIN вместо файла |
| `--rules` | Путь к файлу с кастомными правилами |

## Exit codes

| Код | Описание |
| --- | --- |
| `0` | Уязвимости не найдены ИЛИ найдены но включен silent режим |
| `1` | Найдены уязвимости (без silent) ИЛИ ошибка выполнения |

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
│   └── vul_parser/
│       └── main.go                         # Точка входа
├── internal/
│   ├── domain/models
│   │   ├── vulnerability.go            # Модель уязвимости
│   │   └── rule.go                     # Модель правила
│   ├── config/
│   │   └── flags.go                    # Парсинг флагов
│   ├── rules/
│   │   ├── loader.go                   # Загрузчик правил
│   │   └── builtin.yaml                # Встроенные правила
│   ├── parser/
│   │   └── parser.go                   # Парсер JSON/YAML
│   ├── checker/
│   │   └── checker.go                  # Движок проверки
│   ├── output/
│   │   └── printer.go                  # Вывод результатов
│   ├── permission/
│   │   └── checker.go             # Проверка прав доступа
|   |
│   └── test-configs                    # Готовые конфиги для тестирования
│       ├── all_vulnerabilites.json     # Конфиг со всеми уязвимостями
│       ├── edge_cases.json             # Конфиг с граничными случаями
|       ├── no_vulnerabilites.json      # Конфиг без уязвимостей
|       └── partial_vulnerabilites.yaml # Конфиг формата YAML
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

## Примеры

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