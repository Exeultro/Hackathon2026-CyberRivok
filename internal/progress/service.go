package progress

import (
	"context"
	"time"

	"typing/internal/gamification"
)

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{repo: repo}
}

func (s *Service) Get(ctx context.Context, userID string) (Progress, error) {
	xp, _, err := s.repo.GetXP(ctx, userID)
	if err != nil {
		return Progress{}, err
	}
	return s.repo.BuildProgress(ctx, userID, xp)
}

func (s *Service) Catalog(ctx context.Context, userID string) ([]Achievement, error) {
	m, err := s.repo.AchievementProgress(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]Achievement, 0, len(m))
	for _, a := range m {
		out = append(out, a)
	}
	sortAchievements(out)
	return out, nil
}

func (s *Service) ApplySession(ctx context.Context, userID string, xpGained int, metrics map[string]float64) (*ApplyResult, error) {
	if err := s.repo.EnsureRow(ctx, userID); err != nil {
		return nil, err
	}

	xp, levelBefore, err := s.repo.GetXP(ctx, userID)
	if err != nil {
		return nil, err
	}

	newXP := xp + int64(xpGained)
	newLevel := gamification.LevelFromXP(newXP)

	if err := s.repo.SetXPAndLevel(ctx, userID, newXP, newLevel); err != nil {
		return nil, err
	}

	prog, err := s.repo.BuildProgress(ctx, userID, newXP)
	if err != nil {
		return nil, err
	}

	unlocked, err := s.repo.UnlocksInRange(ctx, levelBefore, newLevel)
	if err != nil {
		return nil, err
	}

	newAch, err := s.updateAchievements(ctx, userID, metrics)
	if err != nil {
		return nil, err
	}

	return &ApplyResult{
		LevelBefore:     levelBefore,
		Progress:        prog,
		NewAchievements: newAch,
		Unlocked:        unlocked,
	}, nil
}

func (s *Service) updateAchievements(ctx context.Context, userID string, metrics map[string]float64) ([]Achievement, error) {
	catalog := gamification.Achievements()

	cur, err := s.repo.AchievementProgress(ctx, userID)
	if err != nil {
		return nil, err
	}

	var newOnes []Achievement
	for _, def := range catalog {
		val := metrics[def.Metric]
		alreadyUnlocked := false
		if existing, ok := cur[def.Code]; ok && existing.UnlockedAt != nil {
			alreadyUnlocked = true
		}
		if val < def.Target {
			if err := s.repo.UpsertAchievement(ctx, userID, def.Code, val, false); err != nil {
				return nil, err
			}
			continue
		}
		if !alreadyUnlocked {
			if err := s.repo.UpsertAchievement(ctx, userID, def.Code, val, true); err != nil {
				return nil, err
			}
			now := timeNowRFC3339()
			newOnes = append(newOnes, Achievement{
				Code: def.Code, Title: def.Title, Description: def.Description,
				Icon: def.Icon, Target: def.Target, Progress: val, UnlockedAt: &now,
			})
		} else {
			if err := s.repo.UpsertAchievement(ctx, userID, def.Code, val, false); err != nil {
				return nil, err
			}
		}
	}
	return newOnes, nil
}

func sortAchievements(a []Achievement) {
	for i := 1; i < len(a); i++ {
		for j := i; j > 0 && a[j].Code < a[j-1].Code; j-- {
			a[j], a[j-1] = a[j-1], a[j]
		}
	}
}
func timeNowRFC3339() string { return time.Now().UTC().Format(time.RFC3339) }

func (s *Service) GetProgress(ctx context.Context, userID string) (Progress, error) {
	return s.Get(ctx, userID)
}

func (s *Service) ListAchievements(ctx context.Context, userID string) ([]Achievement, error) {
	return s.Catalog(ctx, userID)
}

func (s *Service) ListAchievementsAnonymous(ctx context.Context) ([]Achievement, error) {
	return s.repo.CatalogAll(ctx)
}

func (s *Service) UnlockedModesAndDifficulty(ctx context.Context, userID string) ([]string, int, error) {
	p, err := s.Get(ctx, userID)
	if err != nil {
		return nil, 0, err
	}
	return p.UnlockedModes, p.MaxDifficulty, nil
}
func (s *Service) Reset(ctx context.Context, userID string) error {
	return s.repo.Reset(ctx, userID)
}

func (s *Service) Level(ctx context.Context, userID string) (int, error) {
	p, err := s.Get(ctx, userID)
	if err != nil {
		return 1, err
	}
	return p.Level, nil
}
