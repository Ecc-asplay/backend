DBName = sssa 

# PSQL　ダウンロード　と　作成
postgres:
	docker run --name psql -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:latest

# コンテナ 削除
dropPsql:
	docker stop psql || true
	docker rm psql || true

# DB 作成
createDB:
	docker exec -it psql createdb --username=root --owner=root $(DBName)

# DB 削除
dropDB:
	docker exec -it psql dropdb $(DBName)



.PHONY: postgres dropPsql createDB dropDB