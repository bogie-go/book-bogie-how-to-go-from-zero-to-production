# Chapter 17: Keep the Shop Healthy

Code for this chapter: the app as it stands at the end of Chapter 17.

Bogie version: v0.6.8 (created with `bogie new depot`, module name `depot`).

## Run it

You need Go 1.26+, Docker with Compose v2, and Bogie v0.6.8 (`go install github.com/bogie-go/bogie@v0.6.8`).

```sh
cd ch17-keep-the-shop-healthy
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
- `GET /status` checks the database (`Store.Ping`) and the job queue (`QueueHealth`, a hand-written query on River's `river_job` table) with a two-second limit, answering 200 or 503. It is kept out of the OpenAPI contract, like `/healthz`: `test/requests/contract_test.go` skips it.
- `ops/backup-depot.sh` is the nightly backup script the chapter has you create on the server (it is not used in development).
- The chapter's operations steps (uptime monitor, stopping the production worker, the cron job, `rclone`) need a live server and were not run. On the box I ran the same `pg_dump --format=custom` / `pg_restore --no-owner` rehearsal against a local Postgres: the restored `orders` count matched, and the scratch database was dropped. I also ran the app with `DEPOT_ENV=production` and saw the JSON request log lines the chapter shows (`path` is the route pattern, `took` is in nanoseconds).
- `bin/ci` passes.
