package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Item struct {
	Caption string
	Weight  float32
	Number  int
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
	resp, err := http.PostForm("http://localhost:8080/item", url.Values{
		"caption": {item.Caption},
		"weight":  {fmt.Sprintf("%.2f", item.Weight)},
		"number":  {strconv.Itoa(item.Number)},
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := bufio.NewReader(resp.Body).ReadString('\n')
		return fmt.Errorf("не удалось добавить элемент, статус: %d, сообщение: %s", resp.StatusCode, strings.TrimSpace(body))
	}
	return nil
}

func getItem(caption string) (Item, error) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/item/%s", url.PathEscape(caption)))
	if err != nil {
		return Item{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := bufio.NewReader(resp.Body).ReadString('\n')
		return Item{}, fmt.Errorf("элемент не найден или ошибка, статус: %d, сообщение: %s", resp.StatusCode, strings.TrimSpace(body))
	}

	line, _ := bufio.NewReader(resp.Body).ReadString('\n')
	parts := strings.Split(strings.TrimSpace(line), ",")
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
