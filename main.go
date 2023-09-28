package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	bardApiKey := os.Getenv("BARD_API_KEY")

	bot, err := NewChatbot(bardApiKey)
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
