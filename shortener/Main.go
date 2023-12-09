package main

import (
	"fmt"
	"net/http"
	"shortener/hashtable"
	"strings"
	"sync"

	"github.com/gorilla/mux"
)

// DBLock представляет структуру для блокировки хэш-таблицы
type DBLock struct {
	hashTableLock sync.Mutex           // Мьютекс для операций с хэш-таблицей
	hashTable     *hashtable.HashTable // Хэш-таблица, содержащая соответствия сокращенных и оригинальных ссылок
}

// Глобальная переменная для управления блокировками
var dbLock DBLock

func main() {
	// Инициализация хэш-таблицы и чтение из файла
	dbLock.hashTable = hashtable.NewHashTable()
	dbLock.hashTable.ReadFromFile("hashtable.txt")

	// Создание маршрутизатора с использованием gorilla/mux
	router := mux.NewRouter()

	// Назначение обработчиков для конечных точек
	router.HandleFunc("/", shortenHandler).Methods(http.MethodPost)
	router.HandleFunc("/{shortURL}", redirectHandler).Methods(http.MethodGet)
	router.HandleFunc("/getFullURL/{shortURL}", getFullURLHandler).Methods(http.MethodGet) // Добавлен новый обработчик для получения полной ссылки

	port := 8080
	fmt.Printf("Сервер запущен на порту %d\n", port)
	// Запуск сервера с указанием маршрутизатора
	http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}

// shortenHandler обрабатывает запросы на сокращение ссылок
func shortenHandler(w http.ResponseWriter, r *http.Request) {
	// Проверка метода запроса
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Чтение тела запроса
	link := r.FormValue("link")

	// Генерация уникального ключа для сокращенной ссылки
	shortURL := generateShortURL(link)

	// Получение блокировки для безопасных операций с хэш-таблицей
	dbLock.hashTableLock.Lock()
	defer dbLock.hashTableLock.Unlock()

	// Добавление в хэш-таблицу соответствия сокращенной и оригинальной ссылоки
	dbLock.hashTable.Push(shortURL, link)

	// Запись обновленной хэш-таблицы в файл
	dbLock.hashTable.WriteToFile("hashtable.txt")

	// Отправка сокращенной ссылки в ответе
	w.Write([]byte(shortURL))
}

// redirectHandler обрабатывает запросы на перенаправление по сокращенным ссылкам
func redirectHandler(w http.ResponseWriter, r *http.Request) {
	// Проверка метода запроса
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Извлечение сокращенной части URL
	shortURL := strings.TrimPrefix(r.URL.Path, "/")

	// Получение блокировки для безопасных операций с хэш-таблицей
	dbLock.hashTableLock.Lock()
	defer dbLock.hashTableLock.Unlock()

	// Получение оригинальной ссылки из хэш-таблицы
	link, found := dbLock.hashTable.Search(shortURL)
	if !found {
		http.Error(w, "Ссылка не найдена", http.StatusNotFound)
		return
	}

	// Перенаправление пользователя на оригинальную ссылку
	http.Redirect(w, r, link, http.StatusSeeOther)
}

// getFullURLHandler обрабатывает запросы на получение полной ссылки
func getFullURLHandler(w http.ResponseWriter, r *http.Request) {
	// Проверка метода запроса
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Извлечение сокращенной части URL
	shortURL := strings.TrimPrefix(r.URL.Path, "/getFullURL/")

	// Получение блокировки для безопасных операций с хэш-таблицей
	dbLock.hashTableLock.Lock()
	defer dbLock.hashTableLock.Unlock()

	// Получение оригинальной ссылки из хэш-таблицы
	link, found := dbLock.hashTable.Search(shortURL)
	if !found {
		http.Error(w, "Ссылка не найдена", http.StatusNotFound)
		return
	}

	// Отправка полной ссылки в теле ответ
	w.Write([]byte(link))
}

// generateShortURL генерирует уникальный сокращенный URL на основе хэша оригинальной ссылки
func generateShortURL(originalURL string) string {
	hash := 0
	for _, char := range originalURL {
		// Пример другой хеш-функции: умножение на 17 и сложение с ASCII кодом символа
		hash = 5 * (hash*17/int(len(originalURL)) + int(char) + 17) % 10000000
	}
	// Преобразование полученного хэша в строку
	return fmt.Sprintf("%x", hash)
}
