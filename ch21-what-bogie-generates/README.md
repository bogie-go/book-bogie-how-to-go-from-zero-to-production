# Chapter 21: What Bogie Generates, and Why

Chapter 21 has no code. It is a tour of the code Bogie generates, by layer, read with the finished `depot` app open beside it ("There's no new code in this chapter, and nothing to type"). This folder is a README only.

The finished app is in [Chapter 20](../ch20-healthy/): Bogie Depot as it stands at the end of the book's code. Bogie version: v0.6.8.

Following the tour in that folder:

- the layers are `app/controllers/routes.go` and `app/middlewares`, `app/controllers`, `app/views`, `app/models`, `app/services` and `app/jobs`, `db/adapter`, `db/queries` and `db/sqlc`, and `db/migrate`;
- the `_extra` files (for example `db/queries/categories_extra.sql` and `db/adapter/orders_extra.go`) are the ones no generator writes;
- `bogie doctor` checks the markers (`bogie:wire`, `bogie:routes`, `bogie:jobs`, and the rest) and `bogie routes` lists the mounted routes.
