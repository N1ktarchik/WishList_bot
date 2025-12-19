package interaction

//for beta-test
import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var testerPassword string

func InitPassword() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found, using environment variables")
	}

	testerPassword = os.Getenv("TESTERPASSWORD")

}

func CheckPassword(txt string) bool {

	mas := strings.Split(txt, " ")

	if len(mas) != 2 {
		return false
	}

	pass := mas[1]

	if pass == testerPassword {
		return true
	}

	return false
}
