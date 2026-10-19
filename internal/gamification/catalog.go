package gamification

type Unlock struct {
	Level int
	Type  string // mode | difficulty | theme
	Value string
}

func Unlocks() []Unlock {
	return []Unlock{
		{1, "mode", "words"},
		{1, "difficulty", "1"},
		{1, "theme", "dark"},
		{2, "mode", "sentences"},
		{2, "difficulty", "2"},
		{3, "mode", "symbols"},
		{3, "difficulty", "3"},
		{4, "theme", "light"},
		{5, "mode", "code"},
		{5, "difficulty", "4"},
		{8, "difficulty", "5"},
	}
}

type AchievementDef struct {
	Code        string
	Title       string
	Description string
	Icon        string
	Target      float64
	Metric      string // best_wpm | best_cpm | clean_streak | total_chars | sessions_count | accuracy_avg
}

func Achievements() []AchievementDef {
	return []AchievementDef{
		{"first_steps", "Первые шаги", "Заверши первую тренировку", "keyboard", 1, "sessions_count"},
		{"speed_200", "Разгон", "Достичь 200 символов/мин", "bolt", 200, "best_cpm"},
		{"speed_300", "Молния", "Достичь 300 символов/мин", "bolt", 300, "best_cpm"},
		{"genius", "Гений", "Достичь 100 слов/мин", "trophy", 100, "best_wpm"},
		{"clean_100", "Без единой ошибки", "100 символов подряд без ошибок", "shield", 100, "clean_streak"},
		{"steel_finger", "Стальной палец", "250 символов подряд без ошибок", "shield", 250, "clean_streak"},
		{"perfect", "Идеально", "Тренировка со 100% точностью (не короче 20 символов)", "target", 1, "perfect_session"},
		{"thousand", "Тысячник", "1000 символов суммарно", "clock", 1000, "total_chars"},
		{"marathon", "Марафонец", "10 000 символов суммарно", "flame", 10000, "total_chars"},
		{"streak_7", "Неделя подряд", "7 дней подряд с тренировкой", "flame", 7, "streak_days"},
		{"stars_10", "Звёздный охотник", "10 звёзд на картах", "star", 10, "stars_total"},
		{"stars_30", "Звёздный коллекционер", "30 звёзд на картах", "star", 30, "stars_total"},
		{"combo_50", "Комбо", "Комбо из 50 верных нажатий", "flame", 50, "max_combo"},
		{"coder", "Кодер", "Тренировка в режиме code", "code", 1, "code_sessions"},
		{"polyglot", "Полиглот", "Тренировки на ru и en", "globe", 2, "languages_count"},
	}
}
