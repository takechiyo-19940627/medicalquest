// Define the environments
env "dev" {
  url = "postgres://postgres:postgres@localhost:5432/medicalquest?sslmode=disable"
  dev = true
  migration {
    dir = "file://migrations"
  }
}

env "prod" {
  url = "${env:DATABASE_URL}"
  migration {
    dir = "file://migrations"
  }
}