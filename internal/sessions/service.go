package sessions

import (
	"context"
	"errors"
	"math"
	"strings"
	"time"

	"typing/internal/gamification"
	"typing/internal/maps"
	"typing/internal/progress"
)

const (
	maxKeyStats    = 500
	maxBigramStats = 1000
)

var ErrValidation = errors.New("sessions: validation")

type NodeResult struct {
	NodeID     string  `json:"node_id"`
	Passed     bool    `json:"passed"`
	Stars      int     `json:"stars"`
	BestStars  int     `json:"best_stars"`
	NextNodeID *string `json:"next_node_id"`
}

type KeyStatInput struct {
	Key          string   `json:"key"`
	Hits         int      `json:"hits"`
	Errors       int      `json:"errors"`
	AvgLatencyMS *float64 `json:"avg_latency_ms,omitempty"`
}

type CreateRequest struct {
	Mode               string         `json:"mode"`
	Lang               string         `json:"lang"`
	CodeLanguage       *string        `json:"code_language,omitempty"`
	NodeID             *string        `json:"node_id,omitempty"`
	Difficulty         int            `json:"difficulty"`
	TextIDs            []string       `json:"text_ids,omitempty"`
	StartedAt          time.Time      `json:"started_at"`
	DurationMS         int            `json:"duration_ms"`
	CharsTyped         int            `json:"chars_typed"`
	CharsCorrect       int            `json:"chars_correct"`
	Errors             int            `json:"errors"`
	MaxCombo           int            `json:"max_combo"`
	LongestCleanStreak int            `json:"longest_clean_streak,omitempty"`
	ClientScore        *int           `json:"client_score,omitempty"`
	BonusesUsed        []string       `json:"bonuses_used,omitempty"`
	KeyStats           []KeyStatInput `json:"key_stats"`
	BigramStats        []KeyStatInput `json:"bigram_stats,omitempty"`
}

type Result struct {
	Session         Session
	LevelBefore     int
	Progress        progress.Progress
	NodeResult      *NodeResult
	NewAchievements []progress.Achievement
	Unlocked        []progress.Unlock
	IsPersonalBest  bool
}

type Service struct {
	repo        *Repo
	maps        *maps.Service
	progressSvc *progress.Service
}

func NewService(repo *Repo, progressSvc *progress.Service, mapsSvc *maps.Service) *Service {
	return &Service{repo: repo, progressSvc: progressSvc, maps: mapsSvc}
}

func (s *Service) Create(ctx context.Context, userID string, req CreateRequest) (*Result, error) {
	if err := validate(&req); err != nil {
		return nil, err
	}

	durMin := float64(req.DurationMS) / 60000.0
	cpm := 0.0
	wpm := 0.0
	if durMin > 0 {
		cpm = float64(req.CharsCorrect) / durMin
		wpm = (float64(req.CharsCorrect) / 5.0) / durMin
	}
	acc := 0.0
	if req.CharsTyped > 0 {
		acc = float64(req.CharsCorrect) / float64(req.CharsTyped) * 100.0
	}

	xp := gamification.XPForSession(req.CharsCorrect, req.Difficulty, acc)
	score := gamification.ScoreForSession(req.CharsCorrect, req.MaxCombo, acc)

	bestCPM, err := s.repo.BestCPM(ctx, userID)
	if err != nil {
		return nil, err
	}
	if cpm > bestCPM {
		bestCPM = cpm
	}
	bestWPM := bestCPM / 5.0

	streak, err := s.repo.MaxCleanStreak(ctx, userID)
	if err != nil {
		return nil, err
	}
	if req.LongestCleanStreak > streak {
		streak = req.LongestCleanStreak
	}

	maxCombo, err := s.repo.MaxCombo(ctx, userID)
	if err != nil {
		return nil, err
	}
	if req.MaxCombo > maxCombo {
		maxCombo = req.MaxCombo
	}

	totalChars, err := s.repo.TotalChars(ctx, userID)
	if err != nil {
		return nil, err
	}
	totalChars += int64(req.CharsCorrect)

	count, err := s.repo.SessionsCount(ctx, userID)
	if err != nil {
		return nil, err
	}
	count++

	codeCount, err := s.repo.CodeSessionsCount(ctx, userID)
	if err != nil {
		return nil, err
	}
	if req.Mode == "code" {
		codeCount++
	}

	langs, err := s.repo.LangsWithCurrent(ctx, userID, req.Lang)
	if err != nil {
		return nil, err
	}

	perfect, err := s.repo.PerfectSessionsCount(ctx, userID)
	if err != nil {
		return nil, err
	}
	if req.CharsTyped >= 20 && acc == 100 {
		perfect++
	}

	streakDays, err := s.repo.StreakDays(ctx, userID)
	if err != nil {
		return nil, err
	}

	var nodeResult *NodeResult
	if req.NodeID != nil && *req.NodeID != "" {
		node, _, next, ok := maps.FindNode(*req.NodeID)
		if !ok {
			return nil, ErrValidation
		}
		stars := maps.StarsForAttempt(node, wpm, acc)
		best, err := s.maps.UpsertBest(ctx, userID, node.ID, stars, wpm, acc)
		if err != nil {
			return nil, err
		}
		nodeResult = &NodeResult{
			NodeID:     node.ID,
			Passed:     stars >= 1,
			Stars:      stars,
			BestStars:  best.Stars,
			NextNodeID: next,
		}
	}

	starsTotal, err := s.maps.StarsTotal(ctx, userID)
	if err != nil {
		return nil, err
	}

	metrics := map[string]float64{
		"best_cpm":        bestCPM,
		"best_wpm":        bestWPM,
		"clean_streak":    float64(streak),
		"total_chars":     float64(totalChars),
		"sessions_count":  float64(count),
		"max_combo":       float64(maxCombo),
		"code_sessions":   float64(codeCount),
		"languages_count": float64(langs),
		"perfect_session": float64(perfect),
		"streak_days":     float64(streakDays),
		"stars_total":     float64(starsTotal),
	}

	bestBefore, err := s.repo.BestWPMBefore(ctx, userID, req.StartedAt)
	if err != nil {
		return nil, err
	}
	isPB := wpm > bestBefore && wpm > 0

	sess := Session{
		UserID:             userID,
		Mode:               req.Mode,
		Lang:               req.Lang,
		CodeLanguage:       req.CodeLanguage,
		Difficulty:         req.Difficulty,
		StartedAt:          req.StartedAt,
		DurationMS:         req.DurationMS,
		CharsTyped:         req.CharsTyped,
		CharsCorrect:       req.CharsCorrect,
		Errors:             req.Errors,
		MaxCombo:           req.MaxCombo,
		LongestCleanStreak: req.LongestCleanStreak,
		ClientScore:        req.ClientScore,
		WPM:                round2(wpm),
		CPM:                round2(cpm),
		Accuracy:           round2(acc),
		Score:              score,
		XPGained:           xp,
	}

	var krows []KeyStatRow
	for _, k := range req.KeyStats {
		if !validKey(k.Key, 1) {
			continue
		}
		krows = append(krows, KeyStatRow{
			Kind: "key", Key: k.Key, Lang: req.Lang,
			Hits: k.Hits, Errors: k.Errors, AvgLatencyMS: k.AvgLatencyMS,
		})
	}
	for _, k := range req.BigramStats {
		if !validKey(k.Key, 2) {
			continue
		}
		krows = append(krows, KeyStatRow{
			Kind: "bigram", Key: k.Key, Lang: req.Lang,
			Hits: k.Hits, Errors: k.Errors, AvgLatencyMS: k.AvgLatencyMS,
		})
	}

	saved, err := s.repo.Create(ctx, CreateInput{Session: sess, KeyStats: krows})
	if err != nil {
		return nil, err
	}
	sess = *saved

	pg, err := s.progressSvc.ApplySession(ctx, userID, xp, metrics)
	if err != nil {
		return nil, err
	}

	return &Result{
		Session:         sess,
		LevelBefore:     pg.LevelBefore,
		Progress:        pg.Progress,
		NodeResult:      nodeResult,
		NewAchievements: pg.NewAchievements,
		Unlocked:        pg.Unlocked,
		IsPersonalBest:  isPB,
	}, nil
}

func (s *Service) List(ctx context.Context, userID string, limit, offset int) ([]Session, int, error) {
	return s.repo.List(ctx, userID, limit, offset)
}

func validate(r *CreateRequest) error {
	switch r.Mode {
	case "words", "sentences", "code", "symbols":
	default:
		return ErrValidation
	}
	if r.Lang != "ru" && r.Lang != "en" {
		return ErrValidation
	}
	if r.Difficulty < 1 || r.Difficulty > 5 {
		return ErrValidation
	}
	if r.DurationMS < 1000 {
		return ErrValidation
	}
	if r.CharsTyped < 0 || r.CharsCorrect < 0 || r.Errors < 0 {
		return ErrValidation
	}
	if r.CharsCorrect > r.CharsTyped {
		return ErrValidation
	}
	if r.MaxCombo < 0 || r.LongestCleanStreak < 0 {
		return ErrValidation
	}
	if r.Mode == "code" && r.CodeLanguage == nil {
		return ErrValidation
	}
	if r.CodeLanguage != nil {
		switch *r.CodeLanguage {
		case "python", "javascript", "html", "css":
		default:
			return ErrValidation
		}
	}
	if len(r.KeyStats) > maxKeyStats || len(r.BigramStats) > maxBigramStats {
		return ErrValidation
	}
	if r.StartedAt.IsZero() || r.StartedAt.After(time.Now().Add(time.Minute)) {
		return ErrValidation
	}
	return nil
}

func validKey(k string, wantRunes int) bool {
	k = strings.TrimSpace(k)
	if k == "" {
		return false
	}
	return len([]rune(k)) == wantRunes
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}
