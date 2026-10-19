package settings

type Settings struct {
	Lang         string `json:"lang"`
	Theme        string `json:"theme"`
	Sound        bool   `json:"sound"`
	ShowKeyboard bool   `json:"show_keyboard"`
	DisplayUnit  string `json:"display_unit"`
	RealMode     bool   `json:"real_mode"`
	AutoPair     bool   `json:"auto_pair"`
}

type Patch struct {
	Lang         *string `json:"lang,omitempty"`
	Theme        *string `json:"theme,omitempty"`
	Sound        *bool   `json:"sound,omitempty"`
	ShowKeyboard *bool   `json:"show_keyboard,omitempty"`
	DisplayUnit  *string `json:"display_unit,omitempty"`
	RealMode     *bool   `json:"real_mode,omitempty"`
	AutoPair     *bool   `json:"auto_pair,omitempty"`
}

func Default() Settings {
	return Settings{
		Lang:         "ru",
		Theme:        "dark",
		Sound:        true,
		ShowKeyboard: true,
		DisplayUnit:  "word",
		RealMode:     false,
		AutoPair:     true,
	}
}

func Apply(s Settings, p Patch) Settings {
	if p.Lang != nil {
		s.Lang = *p.Lang
	}
	if p.Theme != nil {
		s.Theme = *p.Theme
	}
	if p.Sound != nil {
		s.Sound = *p.Sound
	}
	if p.ShowKeyboard != nil {
		s.ShowKeyboard = *p.ShowKeyboard
	}
	if p.DisplayUnit != nil {
		s.DisplayUnit = *p.DisplayUnit
	}
	if p.RealMode != nil {
		s.RealMode = *p.RealMode
	}
	if p.AutoPair != nil {
		s.AutoPair = *p.AutoPair
	}
	return s
}
