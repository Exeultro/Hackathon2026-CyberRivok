package maps

type StarThreshold struct {
	MinWPM      float64 `json:"min_wpm"`
	MinAccuracy float64 `json:"min_accuracy"`
}

type Node struct {
	ID             string           `json:"id"`
	Order          int              `json:"order"`
	Title          string           `json:"title"`
	Mode           string           `json:"mode"`
	Difficulty     int              `json:"difficulty"`
	CodeLanguage   *string          `json:"code_language,omitempty"`
	TextCount      int              `json:"text_count"`
	StarThresholds [3]StarThreshold `json:"star_thresholds"`
}

type MapDef struct {
	ID               string `json:"id"`
	Title            string `json:"title"`
	Description      string `json:"description"`
	Color            string `json:"color"`
	RecommendedLevel int    `json:"recommended_level"`
	Nodes            []Node `json:"nodes"`
}

func sp(s string) *string { return &s }

func Catalog() []MapDef {
	py := "python"
	js := "javascript"

	return []MapDef{
		{
			ID:               "journey",
			Title:            "Путешествие",
			Description:      "Спокойный старт: слова и короткие предложения",
			Color:            "#22c55e",
			RecommendedLevel: 1,
			Nodes: []Node{
				{"journey-1", 1, "Первые слова", "words", 1, nil, 1,
					[3]StarThreshold{{30, 90}, {50, 95}, {70, 98}}},
				{"journey-2", 2, "Слова посложнее", "words", 2, nil, 1,
					[3]StarThreshold{{40, 90}, {60, 95}, {80, 98}}},
				{"journey-3", 3, "Первые предложения", "sentences", 1, nil, 1,
					[3]StarThreshold{{35, 90}, {55, 95}, {75, 98}}},
				{"journey-4", 4, "Символы", "symbols", 1, nil, 1,
					[3]StarThreshold{{30, 85}, {50, 92}, {70, 96}}},
			},
		},
		{
			ID:               "chase",
			Title:            "Погоня",
			Description:      "Скорость важнее точности",
			Color:            "#f59e0b",
			RecommendedLevel: 3,
			Nodes: []Node{
				{"chase-1", 1, "Быстрые слова", "words", 3, nil, 1,
					[3]StarThreshold{{60, 90}, {80, 93}, {100, 96}}},
				{"chase-2", 2, "Предложения на скорость", "sentences", 2, nil, 1,
					[3]StarThreshold{{60, 90}, {85, 93}, {110, 96}}},
				{"chase-3", 3, "Символы на скорость", "symbols", 2, nil, 1,
					[3]StarThreshold{{50, 88}, {70, 92}, {90, 95}}},
			},
		},
		{
			ID:               "search",
			Title:            "Поиск",
			Description:      "Точность и сложные тексты",
			Color:            "#3b82f6",
			RecommendedLevel: 5,
			Nodes: []Node{
				{"search-1", 1, "Точные слова", "words", 4, nil, 1,
					[3]StarThreshold{{50, 95}, {70, 97}, {90, 99}}},
				{"search-2", 2, "Классическая проза", "sentences", 3, nil, 1,
					[3]StarThreshold{{50, 94}, {70, 96}, {90, 98}}},
				{"search-3", 3, "Код на python", "code", 2, &py, 1,
					[3]StarThreshold{{30, 92}, {45, 95}, {60, 98}}},
			},
		},
		{
			ID:               "siege",
			Title:            "Осада",
			Description:      "Финальный экзамен: код, символы, сложные тексты",
			Color:            "#ef4444",
			RecommendedLevel: 8,
			Nodes: []Node{
				{"siege-1", 1, "Сложные символы", "symbols", 4, nil, 1,
					[3]StarThreshold{{40, 90}, {60, 94}, {80, 97}}},
				{"siege-2", 2, "Код на javascript", "code", 3, &js, 1,
					[3]StarThreshold{{35, 93}, {50, 96}, {70, 98}}},
				{"siege-3", 3, "Сложная проза", "sentences", 5, nil, 1,
					[3]StarThreshold{{55, 95}, {75, 97}, {95, 99}}},
			},
		},
	}
}

func FindNode(nodeID string) (node Node, mapID string, nextNodeID *string, ok bool) {
	for _, m := range Catalog() {
		for i, n := range m.Nodes {
			if n.ID == nodeID {
				var next *string
				if i+1 < len(m.Nodes) {
					v := m.Nodes[i+1].ID
					next = &v
				}
				return n, m.ID, next, true
			}
		}
	}
	return Node{}, "", nil, false
}
