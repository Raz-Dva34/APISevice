package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Item struct {
	Caption string  `json:"caption"`
	Weight  float32 `json:"weight"`
	Number  int     `json:"number"`
}

type Store struct {
	items []Item
}

func main() {
	store := &Store{items: []Item{}}

	http.HandleFunc("/item", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handlePostItem(w, r, store)
		} else {
			http.Error(w, "Метод не поддерживается!", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/item/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handleGetItem(w, r, store)
		} else {
			http.Error(w, "Метод не поддерживается!", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func handlePostItem(w http.ResponseWriter, r *http.Request, store *Store) {
	var newItem Item

	// Парсим JSON из тела запроса
	if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
		http.Error(w, "Некорректные данные JSON", http.StatusBadRequest)
		return
	}

	// Проверяем, что все поля заполнены
	if newItem.Caption == "" || newItem.Weight < 0 || newItem.Number < 0 {
		http.Error(w, "Все поля обязательны и должны быть положительными", http.StatusBadRequest)
		return
	}

	// Проверяем, что caption уникален
	for _, item := range store.items {
		if item.Caption == newItem.Caption {
			http.Error(w, "Caption должен быть уникальным", http.StatusBadRequest)
			return
		}
	}

	// Добавляем новый элемент в store.items
	store.items = append(store.items, newItem)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s,%.2f,%d\n", newItem.Caption, newItem.Weight, newItem.Number)
}

func handleGetItem(w http.ResponseWriter, r *http.Request, store *Store) {
	caption := strings.TrimPrefix(r.URL.Path, "/item/")
	if caption == "" {
		http.Error(w, "Caption обязателен", http.StatusBadRequest)
		return
	}

	for _, item := range store.items {
		if item.Caption == caption {
			fmt.Fprintf(w, "%s,%.2f,%d\n", item.Caption, item.Weight, item.Number)
			return
		}
	}

	http.Error(w, "Элемент не найден", http.StatusNotFound)
}
