#!/bin/bash
go build github.com/doctori/OpenBicycleDatabase/models
go install github.com/doctori/OpenBicycleDatabase/models
go build
go install
~/dev/bin/OpenBicycleDatabase > OBD.log 2>&1