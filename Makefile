# Passkey Sandbox Makefile

.PHONY: up down build docker-up docker-down


# ローカル開発サーバーを起動（Air使用） - メインコマンド
up:
	@echo "🔧 ローカル開発サーバーを起動中（Air使用）..."
	@echo "Dockerコンテナを停止中..."
	@make docker-down 2>/dev/null || true
	@echo "Airでホットリロード開発サーバーを起動中..."
	cd server && air

# ローカル開発サーバーを停止
down:
	@echo "🛑 ローカル開発サーバーを停止中..."
	@pkill -f "air" || true
	@echo "✅ 開発サーバーが停止しました"

# Docker Composeでアプリケーションを起動
docker-up:
	@echo "🚀 Docker Composeでアプリケーションを起動中..."
	@make down 2>/dev/null || true
	docker-compose up -d
	@echo "✅ アプリケーションが起動しました: http://localhost:8080"

# Docker Composeでアプリケーションを停止
docker-down:
	@echo "🛑 Docker Composeでアプリケーションを停止中..."
	docker-compose down 2>/dev/null || true
	@echo "✅ アプリケーションが停止しました"

# Docker Composeでイメージをビルドして起動
build:
	@echo "🔨 Docker Composeでイメージをビルドして起動中..."
	@make down 2>/dev/null || true
	docker-compose up --build -d
	@echo "✅ アプリケーションがビルドされ起動しました: http://localhost:8080"
