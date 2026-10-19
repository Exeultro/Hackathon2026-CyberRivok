package maps

import (
	"context"
)

type NodeView struct {
	Node
	Stars        int      `json:"stars"`
	BestWPM      *float64 `json:"best_wpm"`
	BestAccuracy *float64 `json:"best_accuracy"`
}

type MapView struct {
	ID               string     `json:"id"`
	Title            string     `json:"title"`
	Description      string     `json:"description"`
	Color            string     `json:"color"`
	RecommendedLevel int        `json:"recommended_level"`
	Ready            bool       `json:"ready"`
	StarsTotal       int        `json:"stars_total"`
	StarsMax         int        `json:"stars_max"`
	Nodes            []NodeView `json:"nodes"`
}

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context, userID string, userLevel int) ([]MapView, error) {
	var results map[string]NodeResult
	if userID != "" {
		var err error
		results, err = s.repo.ResultsForUser(ctx, userID)
		if err != nil {
			return nil, err
		}
	}

	out := make([]MapView, 0, 4)
	for _, m := range Catalog() {
		view := MapView{
			ID:               m.ID,
			Title:            m.Title,
			Description:      m.Description,
			Color:            m.Color,
			RecommendedLevel: m.RecommendedLevel,
			Ready:            userLevel >= m.RecommendedLevel,
			StarsMax:         len(m.Nodes) * 3,
		}
		for _, n := range m.Nodes {
			nv := NodeView{Node: n}
			if r, ok := results[n.ID]; ok {
				nv.Stars = r.Stars
				w := r.BestWPM
				a := r.BestAccuracy
				nv.BestWPM = &w
				nv.BestAccuracy = &a
				view.StarsTotal += r.Stars
			}
			view.Nodes = append(view.Nodes, nv)
		}
		out = append(out, view)
	}
	return out, nil
}

func StarsForAttempt(n Node, wpm, accuracy float64) int {
	stars := 0
	for _, th := range n.StarThresholds {
		if wpm >= th.MinWPM && accuracy >= th.MinAccuracy {
			stars++
		}
	}
	return stars
}

func (s *Service) UpsertBest(ctx context.Context, userID, nodeID string, stars int, wpm, accuracy float64) (NodeResult, error) {
	return s.repo.UpsertBest(ctx, userID, nodeID, stars, wpm, accuracy)
}

func (s *Service) StarsTotal(ctx context.Context, userID string) (int, error) {
	return s.repo.StarsTotal(ctx, userID)
}
