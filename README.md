# BOOKIES

### What is this?
This repo contains an app for **Altech** recruitment test.

### Stack
read [stack.md](./stack.md)

### Running instruction
1. Ensure you have [migrate](https://github.com/golang-migrate/migrate) cli installed on your device
2. Run development database. I provide docker compose file to easily running development `PostgreSQL` server.
   * Database name `bookies`
   * Exported port `15432`
   * Default password `bookiesDBpass`
   * Default user `postgres`
3. Run migrations
   ```bash
   migrate -source file://migrations -database postgres://{user}:{password}@{host}:{port}/{db_name} up
   ```
4. **[OPTIONAL]** You can run `sample-data.sql` to populate dev DB with sample data
5. Copy `sample.env` to `.env` in root project and modify it as you need
6. Running app by
   ```bash
   go run cmd/server/main.go
   ```
