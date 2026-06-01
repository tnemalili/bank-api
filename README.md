# Bank API

A lightweight banking API built with Go, Fiber, GORM, SQLite, and AWS SNS integration.

## For Non-Technical Readers

### What this system does

This API helps with simple bank account operations:

- Create a bank account
- Check account details
- Deposit money
- Withdraw money
- Track transaction outcomes (success or failure)

### Typical user flow

1. A customer account is created.
2. The account starts with an initial balance and currency.
3. The customer can deposit or withdraw funds.
4. The system returns a clear result message, such as:
   - Withdrawal successful
   - Insufficient funds
   - Account not found

### Reliability and safety basics

- The API has crash recovery middleware to avoid full service crashes on handler panics.
- Basic validation is done for malformed request bodies.
- Transaction events are published asynchronously so user responses are not blocked.

### Important note

This project is currently a development-oriented API and not a production banking core.
Some enterprise controls (strong auth, full auditing, retries, ledger reconciliation) are not fully implemented yet.

## For Technical Readers

### Stack

- Language: Go 1.24
- Web framework: Fiber v2
- ORM: GORM
- Database: SQLite
- Messaging: AWS SNS (AWS SDK v2)
- Logging: Logrus

### Project structure

- `main.go`: application bootstrap and dependency wiring
- `server/`: HTTP server and middleware setup
- `handlers/`: HTTP handlers (account + transaction endpoints)
- `core/services/`: business service layer
- `repository/relationaldb/`: GORM repositories and SQLite client
- `core/models/`: request/response/domain models
- `messaging/`: SNS publish client

### Runtime configuration

Set the following environment variables before starting the service:

- `API_PORT`: port to run the API on (example: `3540`)
- `API_VERSION`: route version prefix (example: `v1`)
- `DB_HOST`: SQLite database path (example: `/tmp/bank.db`)
- `AWS_REGION`: AWS region for SNS (defaults to `us-east-1` if missing)
- `TRANSACTION_TOPIC`: SNS Topic ARN for transaction event publishing

Example local setup:

```bash
export API_PORT=3540
export API_VERSION=v1
export DB_HOST=/tmp/bank.db
export AWS_REGION=us-east-1
export TRANSACTION_TOPIC=arn:aws:sns:us-east-1:123456789012:transaction-topic
```

### Run locally

```bash
go mod tidy
go run .
```

Base URL example:

```text
http://localhost:3540/api/v1/
```

### API endpoints

1. Health check

- Method: `GET`
- Path: `/api/{version}/health`

```bash
curl -X GET http://localhost:3540/api/v1/health
```

2. Create account

- Method: `POST`
- Path: `/api/{version}/accounts`
- Body:

```json
{
  "accountHolder": "Tshaks",
  "initiationAmount": 300,
  "currency": "R"
}
```

```bash
curl -X POST http://localhost:3540/api/v1/accounts \
	-H "Content-Type: application/json" \
	-d '{"accountHolder":"Tshaks","initiationAmount":300,"currency":"R"}'
```

Returns created account fields including generated `id` and `accountNumber`.

3. Get account

- Method: `GET`
- Path: `/api/{version}/accounts`
- Required header: `X-Account-Id`

```bash
curl -X GET http://localhost:3540/api/v1/accounts \
	-H "X-Account-Id: 3356884001"
```

Lookup supports:

- `id`
- `account_number`
- legacy `account_id` column (if present)

4. Deposit

- Method: `POST`
- Path: `/api/{version}/deposit`
- Required header: `X-Idempotency-Key`
- Body:

```json
{
  "accountId": "3356884001",
  "amount": 50,
  "currency": "R"
}
```

```bash
curl -X POST http://localhost:3540/api/v1/deposit \
	-H "X-Idempotency-Key: dep-001" \
	-H "Content-Type: application/json" \
	-d '{"accountId":"3356884001","amount":50,"currency":"R"}'
```

5. Withdraw

- Method: `POST`
- Path: `/api/{version}/withdraw`
- Required header: `X-Idempotency-Key`
- Body:

```json
{
  "accountId": "3356884001",
  "amount": 100,
  "currency": "R"
}
```

```bash
curl -X POST http://localhost:3540/api/v1/withdraw \
	-H "X-Idempotency-Key: wd-001" \
	-H "Content-Type: application/json" \
	-d '{"accountId":"3356884001","amount":100,"currency":"R"}'
```

### Transaction response shape

Deposit and withdraw currently return a transaction result with fields such as:

- `amount` (value + currency)
- `newBalance` (value + currency)
- `eventId`
- `status`
- `message`
- `createdAt`
- `success`

### Messaging behavior (SNS)

- After deposit/withdraw, handlers publish the transaction event in a goroutine.
- Publish is non-blocking for API responses.
- Current code reads the SNS topic from `TRANSACTION_TOPIC` and passes it as `TopicArn`. This must be a valid SNS Topic ARN in real environments.

### Database behavior

- GORM auto-migrates account and transaction models at startup.
- SQLite is used via `DB_HOST` path.

### Docker

The repository includes a multi-stage Dockerfile.

Build:

```bash
docker build -t bank-api .
```

Run:

```bash
docker run --rm -p 3540:3540 \
	-e API_PORT=3540 \
	-e API_VERSION=v1 \
	-e DB_HOST=/tmp/bank.db \
	-e AWS_REGION=us-east-1 \
	-e TRANSACTION_TOPIC=arn:aws:sns:us-east-1:123456789012:transaction-topic \
	bank-api
```

### Known gaps and improvement backlog

- No authentication/authorization on endpoints.
- Idempotency keys are accepted via `X-Idempotency-Key` header but not yet enforced at persistence level.
- No distributed lock or transaction ledger model.
- SNS topic should be environment-configured as full ARN.
- Input validation can be expanded (currency rules, positive amount checks).
- Observability Telemetry.

## Quick glossary

- **Account ID:** internal account identifier.
- **Account Number:** secondary identifier (same as Account ID) used in account retrieval compatibility logic.
- **Idempotency key:** a request key intended to prevent duplicate processing.
- **SNS:** AWS service for publishing events/messages.

@copyright: Tshakule Nemalili tnemalili@gmail.com +2773-5845-995
