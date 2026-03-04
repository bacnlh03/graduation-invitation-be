# Đọc các biến từ file .env
include .env
export

# Chuỗi kết nối DB dùng cho Goose
DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

.PHONY: run migrate-up migrate-down migrate-status create-migration build docker-build docker-run

# Chạy ứng dụng Go
run:
	go run cmd/main.go

# Lệnh Goose: Cập nhật database lên bản mới nhất
migrate-up:
	goose -dir migrations postgres "$(DB_URL)" up

# Lệnh Goose: Rollback 1 phiên bản
migrate-down:
	goose -dir migrations postgres "$(DB_URL)" down

# Lệnh Goose: Kiểm tra trạng thái các bản migration
migrate-status:
	goose -dir migrations postgres "$(DB_URL)" status

# Lệnh tạo file migration mới (Sử dụng: make create-migration name=ten_file)
create-migration:
	goose -dir migrations create $(name) sql

# Build ứng dụng thành file thực thi
build:
	go build -o bin/app cmd/main.go

# Docker build
docker-build:
	docker build -t graduation-invitation-be .

# Docker run: tự động nạp file .env
docker-run:
	docker run --rm --env-file .env -p $(APP_PORT):$(APP_PORT) graduation-invitation-be