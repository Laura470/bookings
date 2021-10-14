#!/bin/bash

go build -o bookings cmd/web/*.go
./bookings -dbname=bookings -dbuser=postgres -dbpass=laura -cache=false -production=false 