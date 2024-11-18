DBName = asplay

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
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/asplay?sslmode=disable" -verbose up
	
#down
migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/asplay?sslmode=disable" -verbose down

migrateup1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/asplay?sslmode=disable" -verbose up 1
	
#down
migratedown1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/asplay?sslmode=disable" -verbose down 1


# Redis　ダウンロードと作成
redis:
	docker run --name redis -d -p 6379:6379 redis:7.4.1-alpine

# redis コンテナ 削除
dropRedis:
	docker stop redis || true
	docker rm redis || true

# reset docker db
resetDB:
	make dropPsql
	make dropRedis
	make postgres
	make redis
	sleep 3
	make createDB
	make migrateup
	make test

#Sqlc
sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go


.PHONY: postgres dropPsql createDB dropDB migrateup migratedown migrateup1 migratedown1 sqlc 
		redis dropRedis
		resetDB
		test server
