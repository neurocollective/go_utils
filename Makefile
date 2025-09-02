db/local:
	@docker run --name sql-test-pg -p 5432:5432 -e POSTGRES_PASSWORD=postgres -d postgres
	@sleep 2
	@psql -f db/create_tables.sql "postgresql://postgres:postgres@localhost:5432/postgres" 
	@psql -f db/initial_seed.sql "postgresql://postgres:postgres@localhost:5432/postgres"
db/local/down:
	@docker rm -f sql-test-pg
psql:
	@psql "postgresql://postgres:postgres@localhost:5432/postgres"
test/unit:
	go test -v ./sql/index_test.go ./sql/index.go
test/int:
	docker rm -f sql-test-pg
	make db/local
	go test -v ./sql/index.go ./sql/expenditure.go ./sql/sql_reporter_test.go -run TestInsertStructsWithSQLMetaStruct
ahab:
	@docker rm -f sql-test-pg
fmt:
	@go fmt ./db ./sql ./sqlv2
