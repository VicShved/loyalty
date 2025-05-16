#!/bin/bash

gophermarttest \
-test.v -test.run=^TestGophermart$ \
-gophermart-binary-path=cmd/gophermart/gophermart \
-gophermart-host=localhost \
-gophermart-port=8080 \
-gophermart-database-uri="postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable" \
-accrual-binary-path=cmd/accrual/accrual_linux_amd64 \
-accrual-host=localhost \
-accrual-port=8088 \
-accrual-database-uri=postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable
