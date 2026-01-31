# --- UPDATE BAGIAN INI ---
# URL Koneksi: postgres://user:password@host:port/dbname
DB_URL=postgresql://postgres:root@localhost:5432/go_gin_db?sslmode=disable

# --- PASTIKAN COMMAND DI BAWAH SUDAH PAKAI :Z ---

# Migrasi UP (Bikin Tabel)
migrateup:
	docker run --rm -v "$(PWD)/db/migration:/migrations:Z" --network host migrate/migrate \
		-path=/migrations/ -database "$(DB_URL)" up

# Migrasi DOWN (Hapus Tabel/Undo)
migratedown:
	docker run --rm -v "$(PWD)/db/migration:/migrations:Z" --network host migrate/migrate \
		-path=/migrations/ -database "$(DB_URL)" down 1

# Bikin file migrasi baru
new_migration:
	docker run --rm -v "$(PWD)/db/migration:/migrations:Z" migrate/migrate \
		create -ext sql -dir /migrations -seq $(name)

# Generate Code SQLC
sqlc:
	docker run --rm -v "$(PWD):/src:Z" -w /src sqlc/sqlc generate

.PHONY: migrateup migratedown new_migration sqlc