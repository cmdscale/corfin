# SPDX-FileCopyrightText: 2024 CmdScale GmbH
#
# SPDX-License-Identifier: CC0-1.0

test:
	@$(MAKE) -s test-go-cover
	@$(MAKE) -s test-license

test-go-cover:
	go test -cover -coverprofile cover.all.out ./...
	@echo "Check if total coverage > 80%"
	@go tool cover -func cover.all.out | tail -n1 | tee cover.summary.out
	@grep -qE '([89][0-9](\.[0-9])?%)|100\.0%' cover.summary.out
	@rm cover.*.out

test-license:
	reuse lint

fuzz:
	go test -fuzztime=10s -fuzz=NewIBAN
	go test -fuzztime=10s -fuzz=NewBIC
	go test -fuzztime=10s -fuzz=FromIBAN ./bic
