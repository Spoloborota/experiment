# ==========================================
# Docker commands as .mk file
# ==========================================
# Альтернативный подход с .mk расширением
# Включается в Makefile через: include scripts/docker.mk

# Базовые переменные
DOCKER_COMPOSE_FILE := docker-compose.yml
DOCKER_REGISTRY := your-registry.com
IMAGE_TAG := latest

# Docker команды (альтернативная реализация)
.docker.build:
	docker-compose -f $(DOCKER_COMPOSE_FILE) build

.docker.up:
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d

.docker.push:
	docker tag experiment-server $(DOCKER_REGISTRY)/experiment-server:$(IMAGE_TAG)
	docker push $(DOCKER_REGISTRY)/experiment-server:$(IMAGE_TAG) 