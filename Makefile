db/local:
	@docker run --name local-pg -p 5432:5432 -e POSTGRES_PASSWORD=postgres -d postgres
	@sleep 2
	@psql -f db/create_tables.sql "postgresql://postgres:postgres@localhost:5432/postgres" 
	@psql -f db/initial_seed.sql "postgresql://postgres:postgres@localhost:5432/postgres"
psql:
	@psql "postgresql://postgres:postgres@localhost:5432/postgres"
