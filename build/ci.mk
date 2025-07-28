# ==========================================
# CI/CD Makefile
# ==========================================
# Пример отдельного Makefile для CI/CD задач
# Использование: make -f build/Makefile.ci command

.PHONY: ci-test ci-build ci-deploy ci-security-scan ci-lint

# === CI/CD КОМАНДЫ ===

# Полная CI проверка
ci-test:
	@echo "🧪 Запуск CI тестов..."
	go mod verify
	go vet ./...
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

# Сборка для CI
ci-build:
	@echo "🔨 CI сборка..."
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bin/server cmd/server/main.go
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bin/migrations cmd/migrations/main.go

# Lint проверки
ci-lint:
	@echo "📝 Lint проверки..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "❌ golangci-lint не установлен"; \
		exit 1; \
	fi

# Сканирование безопасности
ci-security-scan:
	@echo "🔍 Сканирование безопасности..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "❌ gosec не установлен"; \
	fi

# Deploy (пример)
ci-deploy:
	@echo "🚀 Deploy..."
	@echo "Здесь будут команды для deploy" 