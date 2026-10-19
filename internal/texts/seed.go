package texts

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type seedText struct {
	ID           string
	Mode         string
	Lang         string
	Difficulty   int
	CodeLanguage *string
	Text         string
	Source       *string
}

func Seed(ctx context.Context, pool *pgxpool.Pool) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, t := range seedTexts() {
		if _, err := tx.Exec(ctx,
			`INSERT INTO texts(id, mode, lang, difficulty, code_language, text, source)
			 VALUES ($1,$2,$3,$4,$5,$6,$7)
			 ON CONFLICT (id) DO UPDATE
			   SET mode=EXCLUDED.mode, lang=EXCLUDED.lang,
			       difficulty=EXCLUDED.difficulty,
			       code_language=EXCLUDED.code_language,
			       text=EXCLUDED.text, source=EXCLUDED.source`,
			t.ID, t.Mode, t.Lang, t.Difficulty, t.CodeLanguage, t.Text, t.Source,
		); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func sp(s string) *string { return &s }

func seedTexts() []seedText {
	py := "python"
	js := "javascript"
	html := "html"
	css := "css"

	return []seedText{
		{"s_ru_1", "sentences", "ru", 1, nil,
			"Скорость света в вакууме равна приблизительно 299 792 458 метров в секунду.",
			sp("Физический справочник")},
		{"s_ru_2", "sentences", "ru", 1, nil,
			"Вода закипает при ста градусах Цельсия на уровне моря, а замерзает при нуле.",
			sp("Техническая энциклопедия")},
		{"s_ru_3", "sentences", "ru", 2, nil,
			"Научное познание мира начинается с наблюдения и опыта, а завершается проверяемой гипотезой.",
			sp("Научная проза")},
		{"s_ru_4", "sentences", "ru", 3, nil,
			"Все счастливые семьи похожи друг на друга, каждая несчастливая семья несчастлива по-своему.",
			sp("Л. Н. Толстой, «Анна Каренина»")},
		{"s_ru_5", "sentences", "ru", 3, nil,
			"Мороз и солнце; день чудесный! Ещё ты дремлешь, друг прелестный — пора, красавица, проснись.",
			sp("А. С. Пушкин, «Зимнее утро»")},
		{"s_ru_6", "sentences", "ru", 4, nil,
			"Белеет парус одинокой\nВ тумане моря голубом!..\nЧто ищет он в стране далёкой?\nЧто кинул он в краю родном?..",
			sp("М. Ю. Лермонтов, «Парус»")},
		{"s_ru_7", "sentences", "ru", 5, nil,
			"О, как убийственно мы любим,\nКак в буйной слепоте страстей\nМы то всего вернее губим,\nЧто нашей нежности милей!",
			sp("Ф. И. Тютчев")},

		{"s_en_1", "sentences", "en", 1, nil,
			"Water boils at one hundred degrees Celsius at sea level.",
			sp("Technical reference")},
		{"s_en_2", "sentences", "en", 2, nil,
			"Science begins with observation and ends with a testable hypothesis.",
			sp("Scientific prose")},
		{"s_en_3", "sentences", "en", 3, nil,
			"It was the best of times, it was the worst of times, it was the age of wisdom, it was the age of foolishness.",
			sp("Charles Dickens, A Tale of Two Cities")},
		{"s_en_4", "sentences", "en", 3, nil,
			"All happy families are alike; each unhappy family is unhappy in its own way.",
			sp("Leo Tolstoy, Anna Karenina")},
		{"s_en_5", "sentences", "en", 4, nil,
			"The wind was a torrent of darkness among the gusty trees,\nThe moon was a ghostly galleon tossed upon cloudy seas.",
			sp("Alfred Noyes, The Highwayman")},
		{"s_en_6", "sentences", "en", 5, nil,
			"Because I could not stop for Death,\nHe kindly stopped for me;\nThe carriage held but just ourselves\nAnd immortality.",
			sp("Emily Dickinson")},

		// code python
		{"c_py_1", "code", "en", 1, &py, "def add(a, b):\n    return a + b\n", nil},
		{"c_py_2", "code", "en", 2, &py, "def factorial(n):\n    if n <= 1:\n        return 1\n    return n * factorial(n - 1)\n", nil},
		{"c_py_3", "code", "en", 3, &py, "squares = [x * x for x in range(10) if x % 2 == 0]\n", nil},
		{"c_py_4", "code", "en", 4, &py, "async def fetch(url):\n    async with session.get(url) as resp:\n        return await resp.json()\n", nil},

		// code javascript
		{"c_js_1", "code", "en", 1, &js, "function add(a, b) {\n    return a + b;\n}\n", nil},
		{"c_js_2", "code", "en", 2, &js, "const squares = arr.filter(x => x % 2 === 0).map(x => x * x);\n", nil},
		{"c_js_3", "code", "en", 3, &js, "async function fetchData(url) {\n    const res = await fetch(url);\n    return res.json();\n}\n", nil},
		{"c_js_4", "code", "en", 4, &js, "const greet = (name) => `Hello, ${name}!`;\n", nil},

		// code html
		{"c_html_1", "code", "en", 1, &html, "<div class=\"container\">\n    <h1>Hello</h1>\n</div>\n", nil},
		{"c_html_2", "code", "en", 2, &html, "<button type=\"button\" disabled>Click me</button>\n", nil},

		// code css
		{"c_css_1", "code", "en", 1, &css, ".container {\n    display: flex;\n    gap: 8px;\n}\n", nil},
		{"c_css_2", "code", "en", 2, &css, "button:hover {\n    background: #333;\n    color: #fff;\n}\n", nil},
	}
}
