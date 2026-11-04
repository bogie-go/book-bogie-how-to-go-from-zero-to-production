# Chapter 14: Send Email in the Background

Code for this chapter: the app as it stands at the end of Chapter 14.

Bogie version: v0.6.8 (created with `bogie new depot`, module name `depot`).

## Run it

You need Go 1.26+, Docker with Compose v2, and Bogie v0.6.8 (`go install github.com/bogie-go/bogie@v0.6.8`).

```sh
cd ch14-send-email-in-the-background
make up                              # postgres on :5440
bogie db:prepare                     # create, migrate, seed
bogie s                              # http://localhost:8080
```

Run the tests:

```sh
DEPOT_ENV=test bogie db:prepare      # once, creates the test database
make test
```

## Notes

- This folder is a complete Bogie app, a snapshot of Bogie Depot at the end of this chapter. The generated `README.md` of the app is kept as `README-bogie.md`.
- Credentials are not committed: no `config/master.key` and no `config/credentials.yml.enc`. Without them the app reads its settings from environment variables and its built-in development defaults, so it runs as is. Run `make credentials` to create your own encrypted settings.
- Confirmation email: checkout queues a `send_order_confirmation` River job inside `PlaceOrder`'s own transaction (`AfterOrder`), so the order and its job commit or roll back together. `bogie worker` works the queue; `app/services/order_confirmation` writes the receipt and `app/services/mailer` sends it over SMTP.
- Run the worker beside the server in a second terminal: `bogie worker`. `make up` now also starts Mailpit (SMTP on :1025, inbox at http://localhost:8025).
- SMTP settings come from `config/smtp.go`. The book puts Mailpit's address in the encrypted development credentials; this repo commits no credentials, so the development and test defaults are set in `SmtpConfig.defaults` instead. Production sets `SMTP_ADDRESS`, `SMTP_PORT`, `SMTP_USERNAME`, `SMTP_PASSWORD` and `SMTP_FROM` (nothing real is committed).
- `bogie g job` also adds River to the app (client in `app/application.go`, River's own schema applied with the migrations, `bogie worker`, extra `bin/ci` steps).
- The worker test (a fake `Confirm` that records the order id) is described in the chapter but not printed; it is written for this repo.
- Checked on the box with a real Mailpit: the email arrives after the worker starts; with Mailpit stopped the job goes `retryable` with the connection error and is delivered once Mailpit is back.
- `bin/ci` passes.
