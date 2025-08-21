package listmonkgo_test

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	godotenv.Load()
	code := m.Run()
	os.Exit(code)
}

// func createClient() *listmonkgo.Client {
// 	return listmonkgo.New(
// 		listmonkgo.WithBaseURL(os.Getenv("LISTMONK_URL")),
// 		listmonkgo.WithAPIUser(os.Getenv("API_USER")),
// 		listmonkgo.WithToken(os.Getenv("API_TOKEN")),
// 	)
// }
