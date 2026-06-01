# Code Analysis for the  Bank Account Withdrawal

This document describes **what the original Java snippet does**, then catalogues the
defects and design weaknesses found in it.

## 1. What the code is supposed to do (business capability)

The snippet exposes a single HTTP endpoint that performs a **cash withdrawal** from a
bank account and emits a **domain event** about it. The intended behaviour is:

1. Receive `POST /bank/withdraw` with an `accountId` and an `amount`.
2. Read the account's current balance from a relational database (`accounts` table).
3. If the balance is sufficient, **debit** the account (`balance = balance - amount`).
4. **Publish a `WithdrawalEvent`** to an AWS SNS topic so downstream consumers
   (ledgers, notifications, fraud, analytics) can react.
5. Return a human-readable status string to the caller.

This is the *core domain capability* that must be preserved: **conditionally debit an
account and reliably announce that the debit happened.**
Everything below is about how that capability is implemented, and not whether it should exist.

## 2. Control - Flow Walkthrough of the Original

```
withdraw(accountId, amount):
    SELECT balance FROM accounts WHERE id = accountId
    if balance != null AND balance >= amount:
        UPDATE accounts SET balance = balance - amount WHERE id = accountId
        if rowsAffected > 0: return "Withdrawal successful"
        else:               return "Withdrawal failed"
    else:
        return "Insufficient funds for withdrawal"

    # ---- everything below is unreachable ----

    publish WithdrawalEvent to SNS
    return "Withdrawal successful"
```

## 3. Defect catalogue

**Severity Legend:**

**Critical - Correctness  or Money loss)
Major - Reliability or Design
Minor - Quality or hygiene**

### Critical-1: The SNS event is never published (dead code)

Every branch of the `if/else` **returns**, so the block that builds and publishes the
`WithdrawalEvent` is **unreachable**. The headline feature of the snippet  emitting
the withdrawal event never executes. In most compilers this is also a
"unreachable statement" error. The business capability is *silently broken*.

### Critical-2: Race condition between balance check and debit

The balance is read in one statement and decremented in a *separate* statement, with no
transaction or locking. Two concurrent withdrawals on the same account can both read the
same balance, both pass the `>= amount` check, and both subtract and **overdrawing the
account**. This is a textbook *time-of-check-to-time-of-use* bug and, in a banking core
domain, it directly causes monetary loss. Correctness depends on the debit being a single
**atomic, conditional** operation.

### Critical-3: Dual-write / no atomicity between DB and SNS

Even once the dead code is removed, the design commits the DB update and *then* publishes
to SNS as two independent operations against two systems with no shared transaction:

* If the publish fails or the process crashes after the commit, the money moves but **no
  event is emitted,** hence downstream ledgers/notifications drift out of sync.
* If the publish were moved *before* the commit and the commit then rolled back, a
  **phantom event** would be emitted for a withdrawal that never happened.

There is no way to make a DB write and an SNS publish atomic directly; this needs a
pattern such as the **transactional outbox** (see the optimized design).

### Critical-4: No input validation (negative amount inverts the operation)

`amount` is taken straight from the request with no checks. A **negative** amount turns
`balance - amount` into an *addition*, a withdrawal of `-100` credits the account by
100. Zero, `null`, and absurdly large values are equally unhandled.

### Major-1:  `EmptyResultDataAccessException` on unknown account

`jdbcTemplate.queryForObject(...)` throws when the account id matches no row. That
exception is uncaught, so a request for a non-existent account returns an opaque
**HTTP 500** instead of a clean "not found". A missing account and an insufficient balance
are also conflated into the same `else` branch.

### Major-2: Wrong HTTP semantics

The method always returns a `String` with (implicitly) **HTTP 200 OK**, even for
"Insufficient funds" and "Withdrawal failed". Clients and gateways cannot distinguish
success from a business rejection from a server error. Money operations need explicit,
machine-readable status codes (e.g. `200/201`, `400`, `404`, `409/422`, `500`).

### Major-3: No idempotency

If a client times out and **retries**, the server happily debits a second time. Payment
and withdrawal APIs must be idempotent (e.g. an `Idempotency-Key`) so a retried request
returns the original result without moving money twice.

### Major-4: Business logic lives in the controller

The controller directly issues SQL, builds the SNS client, serializes JSON, and encodes
business rules. There is no service layer, domain model, or repository abstraction. This
hurts **testability, maintainability, and portability** (the rules are welded to Spring
MVC, JDBC, and the AWS SDK).

### Major-5: Hand-rolled JSON serialization

`WithdrawalEvent.toJson()` builds JSON with `String.format`. It performs **no escaping**,
so any string field could break the payload or allow injection, and the format silently
drifts from the type (e.g. `amount` quoted as a string, `accountId` numeric). A real
serializer with a defined schema is required.

### Major-6: SNS client constructed in the controller constructor

`SnsClient` is built inline with a hardcoded/placeholder `Region`, not injected. This
makes the class impossible to unit test without real AWS, couples deployment config to
code, and prevents reuse/mocking.

### Minor-1: Hardcoded configuration

`Region.YOUR_REGION` and the topic ARN `arn:aws:sns:YOUR_REGION:...` are literals. Region,
topic, and DB connection must come from configuration/environment per deployment.

### Minor-2: No observability or audit trail

No structured logging, metrics, tracing, or audit record. For a regulated banking
operation, every withdrawal attempt (accepted *and* rejected) should be auditable, and
the system should expose metrics/latency for monitoring.

### Minor-3: Money modelled loosely

`BigDecimal` avoids float error (good) but there is no **currency**, no scale/rounding
policy, and the DB column type is unspecified. Money should be a first-class value with a
currency and a fixed minor-unit precision.

### Minor-4:  `@Autowired` field injection / mutable state

Field injection hides dependencies and complicates testing; `snsClient` is a mutable
non-final field. Constructor injection of explicit dependencies is preferred.

## 4. What "preserving business functionality" means here

The optimized version keeps the exact capability — *conditionally debit an account and
emit a withdrawal event* — and the same external contract shape (`POST .../withdraw` with
an account and amount, returning the outcome). It changes only the **how**: making the
debit atomic, making event delivery reliable, returning correct statuses, and isolating
the domain so it can be tested and evolved. No business rule is added or removed.


`@copyright: Tshakule Nemalili  tnemalili@gmail.com +2773-5845-995`
