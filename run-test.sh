#!/bin/bash

go test ./... -cover -bench=. -test.benchtime=3s;