# Makefile atualizado
.PHONY: build test lint validate

build:
	go build -o bin/flip cli_engine.go

test:
	# Executa testes LogLine
	./bin/flip test --file tests/ui_schema_validation.logline
	./bin/flip test --file tests/whatsapp_test.logline

lint:
	# Linter de YAML/JSON e shell scripts
	yamllint -c .yamllint ./**/*.logline
	shellcheck validate_logline.sh

validate:
	# Valida sintaxe de todos os arquivos .logline e .jsonl
	bash validate_logline.sh