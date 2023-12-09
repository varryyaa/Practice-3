package main

import (
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
)

func main() {
	for {
		var choice string
		fmt.Print("Выберите действие (1-сократить/2-перейти/3-получить полную ссылку/0-выход): ")
		fmt.Scanln(&choice)

		switch choice {
		case "1":
			var originalURL string
			fmt.Print("Введите URL для сокращения: ")
			fmt.Scanln(&originalURL)
			shortenURL(originalURL)

		case "2":
			var shortURL string
			fmt.Print("Введите сокращенный URL для перехода: ")
			fmt.Scanln(&shortURL)
			redirectURL(shortURL)

		case "3":
			var shortURL string
			fmt.Print("Введите сокращенный URL для получения полной ссылки: ")
			fmt.Scanln(&shortURL)
			getFullURL(shortURL)

		case "0":
			return

		default:
			fmt.Println("Неверный выбор.")
		}
	}
}

// Функция для получения полной ссылки по сокращенному URL
func getFullURL(shortURL string) {
	// Адрес сервера для получения полной ссылки
	var url string
	var shorturl1 string
	if strings.Contains(shortURL, "http://localhost:8080/") {
		shorturl1 = strings.TrimPrefix(shortURL, "http://localhost:8080/")
		url = "http://localhost:8080/getFullURL/" + shorturl1
	} else {
		url = "http://localhost:8080/getFullURL/" + shortURL
	}

	// Выполнение HTTP-запроса методом GET
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer res.Body.Close()

	// Проверка кода статуса
	if res.StatusCode == http.StatusNotFound {
		fmt.Println("Error: Short URL not found")
		return
	} else if res.StatusCode != http.StatusOK {
		fmt.Printf("Error: Status %d\n", res.StatusCode)
		return
	}

	// Чтение и вывод тела ответа
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Вывод полной ссылки
	fmt.Println("Full URL:", string(body))
}

// Функция для сокращения URL
func shortenURL(originalURL string) {
	// Адрес сервера для сокращения
	url := "http://localhost:8080/"

	// Подготовка тела запроса с использованием переданного URL
	payload := strings.NewReader("link=" + originalURL)

	// Создание HTTP-запроса методом POST
	req, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Установка заголовка Content-Type
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Выполнение запроса к серверу сокращения URL
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer res.Body.Close()

	// Чтение и вывод ответа от сервера сокращения
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	fmt.Println("Shortened URL:", "http://localhost:8080/getFullURL/"+string(body))
}

// Функция для перехода по сокращенному URL
func redirectURL(shortURL string) {
	var url string
	if strings.Contains(shortURL, "http://localhost:8080/") {
		url = shortURL
	} else {
		url = "http://localhost:8080/" + shortURL
	}

	err := openURL(url)
	if err != nil {
		fmt.Println("Error opening URL:", err)
	}
}

// Функция для открытия URL в браузере
func openURL(url string) error {
	// Запуск команды с URL в качестве аргумента
	return exec.Command("cmd", "/c", "start", url).Start()
}
