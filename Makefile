# Project metadata
APP_NAME := CRYPTO-WalletApp
PKG := ./test/...
CMD := ./main.go

DB_NAME := wallet_db
DB_USER := crypto_wallet
SQL_FILE_0 := ./db/scripts/00_dropTables.sql
SQL_FILE_1 := ./db/scripts/01_createTables.sql
SQL_FILE_2 := ./db/scripts/02_prepData.sql

# Go environment
GO := go

# Default target
all: build

run-sql:
	psql -U $(DB_USER) -d $(DB_NAME) -f $(SQL_FILE_0)
	psql -U $(DB_USER) -d $(DB_NAME) -f $(SQL_FILE_1)
	psql -U $(DB_USER) -d $(DB_NAME) -f $(SQL_FILE_2)


# Build the binary
build:
	$(GO) build -o $(APP_NAME) $(CMD)

# Run the app
run:
	$(GO) run $(CMD)

# Run unit tests
test:
	$(GO) test -v $(PKG)

# Clean up binaries
clean:
	@rm -f $(APP_NAME)