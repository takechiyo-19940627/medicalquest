// Define the environments
env "dev" {
  src = "ent://infrastructure/ent/schema"
  url = "postgres://postgres:postgres@localhost:5432/medicalquest?sslmode=disable"
  migration {
    dir = "file://migrations"
  }
}
