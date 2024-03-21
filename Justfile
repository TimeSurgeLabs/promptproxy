set dotenv-load

dev:
  go run main.go serve --dir pb_data --http=0.0.0.0:8090
