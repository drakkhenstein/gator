import (
    "fmt"
    "log"

	"github.com/drakkhenstein/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v"err)
	}
	fmt.Println("Read config: %+v\n", cfg.)

	err = cfg.SetUser("christian")
	if err != nil {
		log.Fatalf("error setting user: %v", err)
	}

	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Println("Read config after update: %+v\n", cfg)
}
