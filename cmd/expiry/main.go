package main

import(
	"github.com/zvdv/ECSS-Lockers/internal/email"

)

func init() {
	if err := godotenv.Load(); err != nil {
		logger.Error.Fatal(err)
	}
	dbURL := fmt.Sprintf(
		"%s?authToken=%s",
		env.MustEnv("DATABASE_URL"),
		env.MustEnv("DATABASE_AUTH_TOKEN"))
	database.Connect(dbURL)
	email.Initialize()
}

func main() {
	
}