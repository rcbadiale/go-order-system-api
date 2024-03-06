createmigration:
	migrate create -ext=sql -dir=migrations -seq init

migrate:
	migrate -path=migrations -database "mysql://root:root@tcp(localhost:3306)/orders" -verbose up

migratedown:
	migrate -path=migrations -database "mysql://root:root@tcp(localhost:3306)/orders" -verbose down

.PHONY: createmigration migrate migratedown
