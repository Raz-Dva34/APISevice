package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Item struct {
	Caption string  `json:"caption"`
	Weight  float32 `json:"weight"`
	Number  int     `json:"number"`
}

func generateItems(count int) []Item {
	rand.Seed(time.Now().UnixNano())
	items := make([]Item, count)
	for i := 0; i < count; i++ {
		items[i] = Item{
			Caption: fmt.Sprintf("Item%d-%d", i, rand.Intn(1000000)),
			Weight:  rand.Float32() * 100,
			Number:  rand.Intn(1000),
		}
	}
	return items
}

func sendItem(item Item) error {
	jsonData, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("не удалось сериализовать элемент: %v", err)
	}

	resp, err := http.Post("http://localhost:8080/item", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("не удалось отправить запрос: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("не удалось добавить элемент, статус: %d, сообщение: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	return nil
}

func getItem(caption string) (Item, error) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/item/%s", caption))
	if err != nil {
		return Item{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return Item{}, fmt.Errorf("элемент не найден или ошибка, статус: %d, сообщение: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Item{}, fmt.Errorf("не удалось прочитать ответ: %v", err)
	}

	parts := strings.Split(strings.TrimSpace(string(body)), ",")
	if len(parts) != 3 {
		return Item{}, fmt.Errorf("неожиданный формат ответа")
	}

	weight, _ := strconv.ParseFloat(parts[1], 32)
	number, _ := strconv.Atoi(parts[2])

	return Item{Caption: parts[0], Weight: float32(weight), Number: number}, nil
}

func main() {
	items := generateItems(5)
	for _, item := range items {
		if err := sendItem(item); err != nil {
			fmt.Printf("Ошибка отправки: %v\n", err)
		} else {
			fmt.Printf("Элемент отправлен: %+v\n", item)
		}
	}

	if len(items) > 0 {
		if item, err := getItem(items[0].Caption); err != nil {
			fmt.Printf("Ошибка получения: %v\n", err)
		} else {
			fmt.Printf("Получен элемент: %+v\n", item)
		}
	}
}
