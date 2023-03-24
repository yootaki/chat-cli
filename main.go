package main

import (
	"fmt"
	"os"
	"strings"

	deepl "./deepl"
	openai "./openai"
)

func main() {
	deeplClient, openaiClient := initClient()

	for {
		input, err := getInputFromUser()
		if err != nil {
			fmt.Println("入力の読み取りに失敗しました:", err)
			continue
		}

		translated, err := deeplClient.Translate([]string{input}, "EN")
		if err != nil {
			fmt.Println("翻訳に失敗しました:", err)
			continue
		}

		response, err := openaiClient.GetChatResponse(translated)
		if err != nil {
			fmt.Println("OpenAI APIに接続できませんでした:", err)
			continue
		}

		if len(response) != 0 {
			result, err := deeplClient.Translate(response, "JA")
			if err != nil {
				fmt.Println("翻訳に失敗しました:", err)
				continue
			}
			fmt.Printf(">> answer   : %s\n", result)
		}
	}
}

func initClient() (*deepl.Client, *openai.Client) {
	deeplAPIKey := os.Getenv("DEEPL_API_KEY")
	if deeplAPIKey == "" {
		fmt.Println("DEEPL_API_KEYが設定されていません。")
		return nil, nil
	}
	deeplClient := deepl.NewClient(deeplAPIKey)

	openaiAPIKey := os.Getenv("CHATGPT_API_KEY")
	if openaiAPIKey == "" {
		fmt.Println("OPENAI_API_KEYが設定されていません。")
		return nil, nil
	}
	openaiClient := openai.NewClient(openaiAPIKey)

	return deeplClient, openaiClient
}

func getInputFromUser() (string, error) {
	fmt.Print(">> question : ")
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}
