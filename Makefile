run:
	go run cmd/main.go

fluentd:
	fluent-bit -c /Users/msw/Desktop/Development/Startup_Companies/Arcipelago_International/repo-personal/pba-graphql/fluent-bit.conf

prometheus:
	prometheus --config.file=/Users/msw/Desktop/Development/Startup_Companies/Arcipelago_International/repo-personal/pba-graphql/prometheus.yml

generate:
	rm -rf graph/generated/generated.go
	rm -rf graph/model/models_gen.go
	go run github.com/99designs/gqlgen generate