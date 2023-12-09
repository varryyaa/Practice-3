package hashtable

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// SIZE определяет размер хэш-таблицы
const SIZE = 500

// HashTableItem представляет элемент хэш-таблицы
type HashTableItem struct {
	Key    string         // Ключ элемента
	Data   string         // Данные элемента
	Next   *HashTableItem // Следующий элемент в случае коллизии
	IsUsed bool           // Флаг, указывающий на использование элемента
}

// HashTable представляет хэш-таблицу
type HashTable struct {
	Items [SIZE]*HashTableItem // Массив элементов хэш-таблицы
	Count int                  // Количество элементов в хэш-таблице
}

// NewHashTable создает и возвращает новую пустую хэш-таблицу
func NewHashTable() *HashTable {
	h := &HashTable{Count: 0}
	for i := 0; i < SIZE; i++ {
		h.Items[i] = &HashTableItem{}
	}
	return h
}

// HashFun вычисляет хэш для заданного ключа
func (h *HashTable) HashFun(key string) int {
	const prime = 53
	hash := 0
	for _, letter := range key {
		hash = (hash*prime + int(letter)) % SIZE
	}
	return hash
}

// CreateItem создает новый элемент хэш-таблицы с заданным ключом и данными
func (h *HashTable) CreateItem(key, data string) *HashTableItem {
	return &HashTableItem{Key: key, Data: data, Next: nil, IsUsed: true}
}

// Push добавляет элемент в хэш-таблицу с указанным ключом и данными
func (h *HashTable) Push(key, data string) string {
	if h.Count >= SIZE {
		return "Таблица полна"
	}

	index := h.HashFun(key)
	item := h.CreateItem(key, data)

	if !h.Items[index].IsUsed {
		h.Items[index] = item
		h.Count++
		return "Элемент добавлен"
	} else {
		current := h.Items[index]
		for current.Next != nil {
			if current.Key == key {
				current.Data = data
				return "Ключи совпали, значение по ключу перезаписано"
			}
			current = current.Next
		}
		if current.Key == key {
			current.Data = data
			return "Ключи совпали, значение по ключу перезаписано"
		} else {
			current.Next = item
			h.Count++
			return "Элемент добавлен"
		}
	}
}

// Search ищет элемент в хэш-таблице по заданному ключу и возвращает его данные
func (h *HashTable) Search(key string) (string, bool) {
	if h.Count == 0 {
		return "", false
	}

	index := h.HashFun(key)
	current := h.Items[index]
	if current == nil || !current.IsUsed {
		return "", false
	} else {
		for current.Next != nil {
			if current.Key == key {
				return current.Data, true
			}
			current = current.Next
		}

		if current.Key == key {
			return current.Data, true
		} else {
			return "", false
		}
	}
}

// Pop удаляет элемент из хэш-таблицы по заданному ключу и возвращает его данные
func (h *HashTable) Pop(key string) string {
	index := h.HashFun(key)
	current := h.Items[index]
	var prev *HashTableItem

	for current != nil {
		if current.Key == key {
			if prev != nil {
				prev.Next = current.Next
			} else {
				h.Items[index] = current.Next
			}
			return "Удаленный элемент: " + current.Data
		}
		prev = current
		current = current.Next
	}
	return "Элемент не найден"
}

// WriteToFile записывает содержимое хэш-таблицы в файл
func (h *HashTable) WriteToFile(filename string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for _, item := range h.Items {
		current := item
		for current != nil {
			if current.IsUsed {
				_, err := writer.WriteString(fmt.Sprintf("%s %s\n", current.Key, current.Data))
				if err != nil {
					return err
				}
			}
			current = current.Next
		}
	}

	return nil
}

// ReadFromFile читает содержимое файла и добавляет элементы в хэш-таблицу
func (h *HashTable) ReadFromFile(filename string) error {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) == 2 {
			key := fields[0]
			data := fields[1]
			h.Push(key, data)
		}
	}

	return scanner.Err()
}
