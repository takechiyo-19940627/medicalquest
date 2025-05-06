// Define the environments
env "dev" {
  src = "ent://internal/ent/schema"
  url = "postgres://postgres:postgres@localhost:5432/medicalquest?sslmode=disable"
  migration {
    dir = "file://migrations"
  }
}
