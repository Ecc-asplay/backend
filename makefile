DBName = asplay
DBSource = postgresql://root:secret@localhost:5432/asplay?sslmode=disable
# DBSource = postgresql://root:Secret123qwecc@asplaytest2.c9s8a6m6kots.us-east-1.rds.amazonaws.com:5432/asplay


# PSQL　ダウンロードと作成
postgres:
	docker run --name psql -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:17-alpine

# psql コンテナ 削除
dropPsql:
	docker stop psql || true
	docker rm psql || true

# DB 作成
createDB:
	docker exec -it psql createdb --username=root --owner=root $(DBName)

# DB 削除
dropDB:
	docker exec -it psql dropdb $(DBName)

# Migrate 初期設定
# createMigrate:
# 	migrate create -ext sql -dir db/migration -seq init_table

#up
migrateup:
	migrate -path db/migration -database "${DBSource}" -verbose up
migrateup1:
	migrate -path db/migration -database "${DBSource}" -verbose up 1
migrateup2:
	migrate -path db/migration -database "${DBSource}" -verbose up 2
migrateup3:
	migrate -path db/migration -database "${DBSource}" -verbose up 3

#down
migratedown:
	migrate -path db/migration -database "${DBSource}" -verbose down
migratedown1:
	migrate -path db/migration -database "${DBSource}" -verbose down 1
migratedown2:
	migrate -path db/migration -database "${DBSource}" -verbose down 2
migratedown3:
	migrate -path db/migration -database "${DBSource}" -verbose down 3

#psql reset
tablereset:
	make migratedown
	make migrateup
	make sqlc


# Redis　ダウンロードと作成
redis:
	docker run --name redis -d -p 6379:6379 redis:7.4.1-alpine

# redis コンテナ 削除
dropRedis:
	docker stop redis || true
	docker rm redis || true

#Sqlc
sqlc:
	sqlc generate

test:
	go test -v -cover ./...

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/Ecc-asplay/backend/db/sqlc Store

server:
	go run main.go


.PHONY: postgres dropPsql createDB dropDB migrateup migratedown migrateup1 migratedown1 migrateup2 migratedown2 migrateup3 migratedown3 sqlc tablereset
		redis dropRedis
		test mock server