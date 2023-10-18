package internal

import (
	"fmt"
	"log"
	"os"

	"github.com/islu/bard-sdk-go/bard"

	_ "github.com/joho/godotenv/autoload"
)

func NewChat() {
	bardApiKey := os.Getenv("BARD_API_KEY")

	bot, err := bard.NewChatbot(bardApiKey)
	if err != nil {
		fmt.Println("new error")
		log.Fatalln(err)
		return
	}

	resp, err := bot.Ask("明天是中秋節")
	if err != nil {
		fmt.Println("ask error")
		log.Fatalln(err)
		return
	}

	fmt.Println(resp.ResponseID)
	fmt.Println(resp.ConversationID)
	fmt.Println(resp.Content)
}
