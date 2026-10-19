package gamification

import "math"

func XPForSession(charsCorrect, difficulty int, accuracy float64) int {
	if charsCorrect <= 0 {
		return 0
	}
	v := float64(charsCorrect) * (accuracy / 100.0) * (1.0 + float64(difficulty)*0.1)
	return int(math.Round(v))
}

func ScoreForSession(charsCorrect, maxCombo int, accuracy float64) int {
	if charsCorrect <= 0 {
		return 0
	}
	base := float64(charsCorrect) * 10
	comboMul := 1.0 + float64(maxCombo)*0.01
	accMul := accuracy / 100.0
	return int(math.Floor(base * comboMul * accMul))
}

func LevelXPThreshold(level int) int64 {
	if level <= 1 {
		return 0
	}
	n := int64(level - 1)
	return 100*n + 25*n*(n-1)
}

func LevelFromXP(xp int64) int {
	lo, hi := 1, 100
	for lo < hi {
		mid := (lo + hi + 1) / 2
		if LevelXPThreshold(mid) <= xp {
			lo = mid
		} else {
			hi = mid - 1
		}
	}
	return lo
}

func XPInLevel(xp int64) int64 {
	lvl := LevelFromXP(xp)
	return xp - LevelXPThreshold(lvl)
}

func XPToNext(xp int64) int64 {
	lvl := LevelFromXP(xp)
	return LevelXPThreshold(lvl+1) - LevelXPThreshold(lvl)
}
