# CRYPTO Wallet App
## Assumption
- A user can own multiple wallets.
- Each user is allowed only __one default wallet__.
- A user __cannot have more than one wallet per currency__.
- Money transfers are supported between __wallets__ and from a __wallet to another user__.
- __Transfers from a user's wallet to another wallet owned by the same user are allowed__, but __not from a wallet to the same user__ (i.e., wallet-to-user transfer is not allowed within the same user).
- __Intra-user transfers__ must be performed __wallet-to-wallet__.
- __Currency conversion is automatically applied__ when transferring funds between wallets of __different currencies__.
- When transferring money from one wallet to __another user__:
    - The system will first try to use a wallet with the __same currency__.
    - If the recipient doesn't have a wallet in that currency, the recipient's __default wallet__ is used instead.
- To mock the currency conversion service, a database is used to store __currency conversion rates__. In a real-world application, this would typically involve calling an external service to fetch __live exchange rates__.
---
## End Points
## GET /users/{id}/wallets/balance
Returns the balance(s) of the wallet(s) belonging to the specified user.

### Path Parameters

| Parameter | Type    | Mandatory | Description          |
|-----------|---------|-----------|----------------------|
| `id`      | integer | yes       | ID of the user       |

### Query Parameters (Optional)

| Parameter   | Type    | Mandatory | Description                                  |
|-------------|---------|-----------|----------------------------------------------|
| `wallet_id` | integer | no        | Filter the result to a specific wallet ID    |


### Example Request
- GET /users/123/wallets/balance
- GET /users/123/wallets/balance?wallet_id=456

Sample Response
```json
{
  "user_info": {
    "id": 4,
    "name": "Danny"
  },
  "wallets": [
    {
      "id": 8,
      "is_default": false,
      "type": "saving",
      "currency": "USD",
      "balance": "5112.00"
    },
    {
      "id": 9,
      "is_default": true,
      "type": "trading",
      "currency": "SGD",
      "balance": "100.23"
    }
  ],
  "total": {
    "currency": "USD",
    "amount": "5186.24"
  }
}
```

## GET /users/{id}/wallets/transactions
Retrieve all wallets and their corresponding transactions history for a given user. Supports optional filtering by wallet_id.

### Path Parameters

| Parameter | Type    | Mandatory | Description                   |
|-----------|---------|-----------|-------------------------------|
| `id`      | integer | yes       | ID of the source wallet       |

### Query Parameters (Optional)

| Parameter   | Type    | Mandatory | Description                               |
|-------------|---------|-----------|-------------------------------------------|
| `wallet_id` | integer | no        | Filter the result to a specific wallet ID |

### Example Request
- GET /users/123/wallets/transactions
- GET /users/123/wallets/transactions?wallet_id=456

Sample Response
```json
{
  "user_info": {
    "id": 4,
    "name": "Danny"
  },
  "wallets": [
    {
      "id": 8,
      "is_default": false,
      "type": "saving",
      "currency": "USD",
      "balance": "5112.00",
      "transactions": [
        {
          "id": 3,
          "type": "deposit",
          "amount": "5022.00",
          "time": "2025-05-19T22:40:20.396731Z"
        },
        {
          "id": 2,
          "type": "withdraw",
          "amount": "10.00",
          "time": "2025-05-19T22:39:52.627552Z"
        },
        {
          "id": 1,
          "type": "deposit",
          "amount": "100.00",
          "time": "2025-05-19T22:39:46.291502Z"
        }
      ]
    },
    {
      "id": 9,
      "is_default": true,
      "type": "trading",
      "currency": "SGD",
      "balance": "100.23",
      "transactions": [
        {
          "id": 4,
          "type": "deposit",
          "amount": "100.23",
          "time": "2025-05-19T22:40:53.729115Z"
        }
      ]
    }
  ],
  "total": {
    "currency": "USD",
    "amount": "5186.24"
  }
}
```

## POST /wallets/{id}/deposit
Deposit funds into the wallet specified by the id. The wallet must exist. This operation increases the wallet's balance and logs a deposit transaction.

### Path Parameters

| Parameter | Type    | Mandatory | Description                |
|-----------|---------|-----------|----------------------------|
| `id`      | integer | yes       | Wallet ID to deposit into  |

### Request Body
| Field    | Type           | Mandatory | Description                     |
|----------|----------------|-----------|---------------------------------|
| `amount` | Decimal Number | yes       | Amount of money to deposit into |

```json
{
  "amount": 100.00
}
```
### Response
204 No Content


## POST /wallets/{id}/withdraw
Withdraw funds from the wallet specified by the id. The wallet must exist and the balance should be greate or equal to the withdraw amount. This operation decreases the wallet's balance and logs a withdraw transaction.

### Path Parameters
| Parameter | Type    | Mandatory | Description                |
|-----------|---------|-----------|----------------------------|
| `id`      | integer | yes       | Wallet ID to withdraw from |


### Request Body
| Field    | Type           | Mandatory | Description                      |
|----------|----------------|-----------|----------------------------------|
| `amount` | Decimal Number | yes       | Amount of money to withdraw from |

```json
{
  "amount": 101.20
}
```
### Response
204 No Content


## POST /wallets/{id}/transfer
Transfer funds from the wallet specified by the id. The wallet must exist and the balance should be greate or equal to the transfer amount. This operation decreases the source wallet's balance and increases the target wallet's balance. It logs a transfer-out transaction in the source wallet and a transfer-in transaction in the target wallet.
Transfer money:
- wallet to wallet
- wallet to user

### Path Parameters

| Parameter | Type    | Mandatory | Description                       |
|-----------|---------|-----------|-----------------------------------|
| `id`      | integer | yes       | Source Wallet ID to transfer from |

### Request Body
| Field                   | Type           | Mandatoryv | Description                                                         |
|-------------------------|----------------|------------|---------------------------------------------------------------------|
| `amount`                | Decimal Number | yes        | Amount of money to transfer from                                    |
| `destination_wallet_id` | integer        | yes or no* | For a `wallet-to-wallet` transfer, this field represents the destination wallet ID  |
| `destination_user_id`   | integer        | yes or no* | For a `wallet-to-user transfer`, it represents the destination user ID. |

> * Either `destination_wallet_id` or `destination_user_id` must be provided — but not both.

Sample `Wallet to Wallet` Request
```json
{
  "amount": 100.00,
  "destination_wallet_id": 456
}
```

Sample `Wallet to User` Request
```json
{
  "amount": 100.00,
  "destination_user_id": 102
}
```
### Response
204 No Content

## Possible Future Improvements
1. Authentication Middleware
    - Add middleware to authenticate API requests using JWT tokens.
    - Useful for protecting endpoints like /withdraw, /transfer, /deposit.
1. Idempotency Keys for Transaction APIs
    - Allow clients to pass an Idempotency-Key header to prevent duplicate deposits or transfers on retry.
    - Store keys temporarily in Redis or a DB table.
1. Central Error Handler / Custom Error Types
    - Use a unified error response format.
    - Define custom error types to improve error propagation and API clarity.
1. Transaction Queue (Async Processing)
    - Queue large or long-running transactions to Redis for background processing (e.g., withdrawal approval).
    - Helps handle future scaling or business rules like fraud checks.
1. Rate Limiting Middleware
    - Prevent abuse on sensitive endpoints (e.g., max 5 withdrawals/min).
    - Use Redis for token bucket or sliding window logic.
1. Pagination for Transaction History
    - Add limit and offset query params to /users/{id}/wallets/transactions.
    - Useful for large datasets and frontend integrations.
1. Context & Timeout Management
    - Add context with timeout for all DB queries and HTTP handlers to avoid resource leaks.
1. Health Check
    - Add /healthz endpoint
    - Helps with deployment and monitoring.



## 🛠️ Setup Guide
### 🧩 Go Setup
__1. Install Go__
- Download and install Go from the official website: golang.org/dl
- Follow the installation instructions for your operating system.
- Ensure Go is added to your system’s `PATH`.

__2. Install Project Dependencies__
Use the following commands to install required packages:
```
go get github.com/jackc/pgx/v5
go get github.com/jackc/pgx/v5/stdlib
go get github.com/gorilla/mux
go get github.com/jackc/pgx/v5/stdlib
go get github.com/shopspring/decimal
go get github.com/spf13/vipe
go get github.com/stretchr/testify/assert
go get github.com/stretchr/testify/require
go get github.com/DATA-DOG/go-sqlmock
```


### 🐘 PostgreSQL
__1. Download & Install__
- Get PostgreSQL from postgresapp.com or use your preferred method.

__2. Create a Database__
- You can use an existing database or create a new one.

__3. Initialize Schema & Data__
- Run the provided SQL scripts to create tables and populate initial data:

```bash
./db/scripts/01_createTables.sql
./db/scripts/02_prepData.sql
```
__4. Configure Connection__
- Update your database settings in:
```
./config/config.yaml
```
Set the following:
- Database name
- User ID
- Password
- Host
- Port

### 💼 Wallet Application
#### ✅ Run Test Cases
```
go test -v ./test/...
# or use Makefile
make test
```
#### ⚙️ Build Application
```
go build -o CRYPTO-WalletApp ./main.go
# or
make build
```
#### 🚀 Run Without Building
```
go run ./main.go
# or
make run
```
#### ▶️ Run After Building
```
./CRYPTO-WalletApp
```
