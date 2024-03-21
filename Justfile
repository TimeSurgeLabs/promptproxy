set dotenv-load

dev:
  go run main.go serve --dir pb_data --http=0.0.0.0:8090

build:
  mkdir -p bin
  go build -o bin/proxy main.go

build-docker:
  docker build -t timesurgelabs/promptproxy .

run-docker:
  docker run -p 8090:8090 -v $(PWD)/pb_data:/pb_data timesurgelabs/promptproxy
