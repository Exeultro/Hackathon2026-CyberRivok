package texts

import (
	"math/rand"
	"strings"
)

var ruWords = map[int][]string{
	1: {"дом", "кот", "лес", "мир", "сок", "нос", "рот", "сад", "час", "шум", "бег", "вал", "год", "дуб", "ёл", "жар", "зал", "кит", "луг", "мак"},
	2: {"окно", "река", "весна", "город", "книга", "лампа", "море", "снег", "трава", "утро", "хлеб", "цветы", "школа", "ягода", "берег", "ветер", "голос", "дерево", "жизнь", "звезда"},
	3: {"компьютер", "клавиатура", "программа", "разработка", "интерфейс", "алгоритм", "структура", "функция", "переменная", "библиотека", "архитектура", "производительность", "безопасность", "тестирование", "документация", "оптимизация"},
	4: {"информационные", "технологический", "высокопроизводительный", "многопоточность", "распределённый", "масштабируемость", "асинхронность", "децентрализация", "интероперабельность", "конфигурируемость"},
	5: {"электроэнцефалографический", "рентгеноструктурный", "гидроэлектростанция", "криптографический", "нейрокибернетический", "психофизиологический", "сверхпроводимость", "экспериментальный", "фундаментальный", "междисциплинарный"},
}

var enWords = map[int][]string{
	1: {"cat", "dog", "sun", "run", "map", "sky", "box", "key", "ice", "red", "top", "win", "far", "low", "new", "old", "big", "fix", "fly", "joy"},
	2: {"apple", "river", "green", "house", "music", "night", "ocean", "paper", "quick", "stone", "table", "under", "water", "young", "zebra", "brain", "cloud", "dream", "earth", "flame"},
	3: {"computer", "keyboard", "software", "developer", "function", "algorithm", "structure", "variable", "interface", "database", "framework", "performance", "security", "testing", "documentation", "optimization"},
	4: {"asynchronous", "distributed", "concurrency", "scalability", "cryptographic", "multithreaded", "deterministic", "idempotent", "serialization", "interoperable"},
	5: {"electroencephalographic", "cryptographically", "interdisciplinary", "superconductivity", "psychophysiological", "neurocybernetic", "counterrevolutionary", "institutionalization", "internationalization", "characteristically"},
}

var ruSymbols = []string{
	"а", "б", "в", "г", "д", "е", "ё", "ж", "з", "и", "й", "к", "л", "м", "н", "о", "п", "р", "с", "т", "у", "ф", "х", "ц", "ч", "ш", "щ", "ъ", "ы", "ь", "э", "ю", "я",
	"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
	".", ",", ";", ":", "!", "?", "-", "_", "(", ")", "@", "#", "$", "%", "&", "*", "+", "=", "/", "\\",
}

var enSymbols = []string{
	"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
	"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
	".", ",", ";", ":", "!", "?", "-", "_", "(", ")", "@", "#", "$", "%", "&", "*", "+", "=", "/", "\\",
}

func GenerateWords(lang string, difficulty, count int) []Item {
	if difficulty < 1 {
		difficulty = 1
	}
	if difficulty > 5 {
		difficulty = 5
	}
	pool := pickWordPool(lang, difficulty)
	out := make([]Item, 0, count)
	for i := 0; i < count; i++ {
		n := 5
		switch difficulty {
		case 1:
			n = 3
		case 2:
			n = 4
		case 3:
			n = 5
		case 4:
			n = 3
		case 5:
			n = 2
		}
		words := make([]string, 0, n)
		for j := 0; j < n; j++ {
			words = append(words, pool[rand.Intn(len(pool))])
		}
		out = append(out, Item{
			ID:         "gen_" + randomID(),
			Mode:       "words",
			Lang:       lang,
			Difficulty: difficulty,
			Text:       strings.Join(words, " "),
		})
	}
	return out
}

func pickWordPool(lang string, difficulty int) []string {
	if lang == "en" {
		if p, ok := enWords[difficulty]; ok {
			return p
		}
		return enWords[1]
	}
	if p, ok := ruWords[difficulty]; ok {
		return p
	}
	return ruWords[1]
}

func GenerateSymbols(lang string, difficulty, count int) []Item {
	if difficulty < 1 {
		difficulty = 1
	}
	if difficulty > 5 {
		difficulty = 5
	}
	pool := ruSymbols
	if lang == "en" {
		pool = enSymbols
	}
	length := 6 + difficulty*2
	out := make([]Item, 0, count)
	for i := 0; i < count; i++ {
		b := make([]string, 0, length)
		for j := 0; j < length; j++ {
			b = append(b, pool[rand.Intn(len(pool))])
		}
		out = append(out, Item{
			ID:         "gen_" + randomID(),
			Mode:       "symbols",
			Lang:       lang,
			Difficulty: difficulty,
			Text:       strings.Join(b, ""),
		})
	}
	return out
}

const idAlphabet = "abcdefghijklmnopqrstuvwxyz0123456789"

func randomID() string {
	b := make([]byte, 10)
	for i := range b {
		b[i] = idAlphabet[rand.Intn(len(idAlphabet))]
	}
	return string(b)
}
