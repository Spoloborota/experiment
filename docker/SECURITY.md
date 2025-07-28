# Container Security Guide

## 🔒 Обзор безопасности

Наши Docker образы реализуют множественные уровни защиты согласно OWASP Container Security рекомендациям.

## 📊 Сравнение уровней безопасности

| Аспект | Standard | Secure | Описание |
|--------|----------|--------|----------|
| **Базовый образ** | Alpine Linux | Distroless | Distroless не содержит shell, пакетных менеджеров |
| **Размер образа** | ~15MB | ~8MB | Меньше компонентов = меньше уязвимостей |
| **Shell доступ** | ❌ Доступен | ✅ Отсутствует | Невозможно выполнять команды в контейнере |
| **Пакеты** | Минимальные | Отсутствуют | Нет apt/apk для установки malware |
| **Версии** | Фиксированные | SHA-пинированные | Гарантированная неизменность |

## 🛡️ Реализованные защиты

### 1. Multi-stage Build
- **Что защищает**: Исходный код, build инструменты не попадают в production
- **Как работает**: Разделение build и runtime окружений
- **Атаки**: Source code leakage, build tool exploitation

### 2. Non-root User
```dockerfile
# Standard version
USER appuser  # UID 1001

# Secure version  
USER nonroot  # UID 65532 (стандарт distroless)
```
- **Что защищает**: Container escape, privilege escalation
- **Как работает**: Процесс не имеет root привилегий в контейнере
- **Атаки**: Kernel exploits, host file system access

### 3. Static Compilation
```dockerfile
CGO_ENABLED=0 GOOS=linux go build -ldflags='-extldflags "-static"'
```
- **Что защищает**: Dynamic library vulnerabilities
- **Как работает**: Все зависимости встроены в бинарник
- **Атаки**: Shared library poisoning, dependency confusion

### 4. Minimal Attack Surface
- **Alpine**: ~5MB, minimal пакеты
- **Distroless**: ~2MB, только runtime
- **Что защищает**: CVE exploits, lateral movement
- **Атаки**: Package vulnerabilities, system tool abuse

### 5. SHA-pinned Images
```dockerfile
FROM golang:1.24.3-alpine3.22@sha256:b4f875e...
```
- **Что защищает**: Supply chain attacks
- **Как работает**: Криптографическая проверка образа
- **Атаки**: Malicious image replacement, registry compromise

## 🔍 Типы атак и защита

### Container Escape
**Атака**: Выход из контейнера на host систему
**Защита**: 
- Non-root user
- Minimal capabilities
- Read-only filesystem (опционально)

### Supply Chain Compromise
**Атака**: Подмена базовых образов или зависимостей
**Защита**:
- SHA-pinned образы
- Dependency verification
- Private registry использование

### Code Injection
**Атака**: Выполнение произвольного кода в контейнере
**Защита**:
- Статическая компиляция
- Отсутствие shell
- Minimal runtime

### Privilege Escalation
**Атака**: Получение root привилегий
**Защита**:
- Non-root user
- Dropped capabilities
- Security contexts

## 📋 Дополнительные рекомендации

### Runtime Security
```bash
# Запуск с дополнительными ограничениями
docker run \
  --read-only \                    # Read-only filesystem
  --no-new-privileges \            # Запрет privilege escalation
  --cap-drop=ALL \                 # Удаление всех capabilities
  --security-opt=no-new-privileges \
  --user=65532:65532 \            # Явно указываем non-root
  your-image
```

### Kubernetes Security Context
```yaml
securityContext:
  runAsNonRoot: true
  runAsUser: 65532
  readOnlyRootFilesystem: true
  allowPrivilegeEscalation: false
  capabilities:
    drop:
      - ALL
  seccompProfile:
    type: RuntimeDefault
```

### Сканирование уязвимостей
```bash
# Триvy сканер
trivy image your-image:latest

# Docker Scout
docker scout cves your-image:latest

# Snyk
snyk container test your-image:latest
```

## 🚨 Мониторинг безопасности

### Файловая система
- Мониторинг изменений в read-only контейнере
- Отслеживание попыток записи в критические директории

### Сетевая активность
- Неожиданные исходящие соединения
- Подключения к внешним ресурсам

### Процессы
- Запуск процессов от root пользователя
- Выполнение системных команд

## 🎯 Выбор версии

### Используйте `server.dockerfile` если:
- Нужна совместимость со старыми системами
- Требуется shell доступ для отладки
- Простота важнее максимальной безопасности

### Используйте `server-secure.dockerfile` если:
- Production окружение
- Высокие требования к безопасности
- Compliance требования (SOC2, ISO27001)
- Zero-trust архитектура

## 📚 Ресурсы

- [OWASP Container Security](https://owasp.org/www-project-container-security/)
- [CIS Docker Benchmark](https://www.cisecurity.org/benchmark/docker)
- [NIST Container Security](https://csrc.nist.gov/publications/detail/sp/800-190/final)
- [Distroless Images](https://github.com/GoogleContainerTools/distroless) 