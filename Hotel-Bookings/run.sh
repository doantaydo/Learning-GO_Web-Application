#!/bin/bash

go build -o bookings /cmd/web/*.go
./bookings -dbName=bookings -dbUser=postgres -dbPass=24072001do -cache=false -production=false