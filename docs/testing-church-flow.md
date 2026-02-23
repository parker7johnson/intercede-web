# Testing the Church Registration Flow Locally

## Prerequisites

- [Stripe CLI](https://stripe.com/docs/stripe-cli) installed (`brew install stripe/stripe-cli/stripe`)
- A Stripe account with a test-mode subscription price created
- PostgreSQL running with the Intercede database accessible
- `.env` file at the project root

---

## 1. Run the Database Migration

If you haven't already, create the `churches` table:

```sql
CREATE TABLE churches (
    id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name                  TEXT NOT NULL,
    admin_email           TEXT NOT NULL,
    church_code           TEXT UNIQUE,
    stripe_session_id     TEXT UNIQUE NOT NULL,
    stripe_customer_id    TEXT,
    stripe_subscription_id TEXT,
    status                TEXT NOT NULL DEFAULT 'pending',
    created_at            TIMESTAMPTZ DEFAULT now()
);
```

---

## 2. Configure Environment Variables

Add these to your `.env` file:

```
STRIPE_SECRET_KEY=sk_test_...
STRIPE_WEBHOOK_SECRET=whsec_...   # filled in step 4
STRIPE_PRICE_ID=price_...
BASE_URL=http://localhost:3000
DB_URL=...
SUPABASE_URL=...
SUPABASE_KEY=...
```

> Get `STRIPE_SECRET_KEY` and `STRIPE_PRICE_ID` from the [Stripe dashboard](https://dashboard.stripe.com/test/apikeys) (make sure you're in **test mode**).

---

## 3. Start the Server

```bash
make build   # or: go run ./cmd/web
```

Confirm it's running at `http://localhost:3000`.

---

## 4. Start the Stripe Webhook Listener

In a second terminal, log in to the Stripe CLI and forward events to your local server:

```bash
stripe login
stripe listen --forward-to localhost:3000/webhooks/stripe --events checkout.session.completed
```

Copy the **webhook signing secret** printed by the CLI (starts with `whsec_`) and set it as `STRIPE_WEBHOOK_SECRET` in your `.env`, then restart the server.

---

## 5. Walk Through the UI

1. Open `http://localhost:3000`
2. Click **Register Your Church**
3. Fill in:
   - **Church Name** — e.g. `Grace Community`
   - **Admin Email** — any valid email
4. Click **Continue to Payment** — you should be redirected to Stripe Checkout
5. On the Stripe-hosted page, use the test card:
   - **Card number:** `4242 4242 4242 4242`
   - **Expiry:** any future date (e.g. `12/34`)
   - **CVC:** any 3 digits (e.g. `123`)
6. Click **Subscribe** — Stripe processes the payment and redirects you to `/church/success`
7. The page displays your generated church code (e.g. `GRAC-7X3K`)

---

## 6. Verify the Webhook Fired

In the Stripe CLI terminal you should see:

```
--> checkout.session.completed [evt_...]
<-- [POST /webhooks/stripe] 200 OK
```

If you see a non-200 response, check the server logs for the error.

---

## 7. Verify the Database

```sql
SELECT name, admin_email, church_code, status
FROM churches
ORDER BY created_at DESC
LIMIT 5;
```

A successful flow shows:
- `status = 'active'`
- `church_code` set to a value like `GRAC-7X3K`

---

## Troubleshooting

| Symptom | Likely cause |
|---|---|
| Redirected to `/church/cancel` | You clicked **Back** on the Stripe page — try again |
| Webhook returns 400 "Invalid webhook signature" | `STRIPE_WEBHOOK_SECRET` doesn't match the CLI secret — restart server after updating `.env` |
| `/church/success` shows "Church not found" | Webhook hasn't fired yet or failed — check CLI terminal |
| `stripe listen` exits immediately | Run `stripe login` first |
| Church code missing from success page | Check DB: `status` may still be `'pending'` if webhook failed |
