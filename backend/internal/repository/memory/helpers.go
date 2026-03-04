package memory

import "poke/backend/internal/model"

func cloneUser(in *model.User) *model.User {
	if in == nil {
		return nil
	}
	out := *in
	if in.RecentMatchIDs != nil {
		out.RecentMatchIDs = append([]string(nil), in.RecentMatchIDs...)
	}
	return &out
}

func cloneRoom(in *model.Room) *model.Room {
	if in == nil {
		return nil
	}
	out := *in
	if in.MemberUserIDs != nil {
		out.MemberUserIDs = append([]string(nil), in.MemberUserIDs...)
	}
	return &out
}

func cloneMatch(in *model.Match) *model.Match {
	if in == nil {
		return nil
	}
	out := *in
	if in.Players != nil {
		out.Players = append([]model.MatchPlayer(nil), in.Players...)
	}
	if in.ScoreEntries != nil {
		out.ScoreEntries = append([]model.MatchScoreEntry(nil), in.ScoreEntries...)
	}
	if in.TotalScores != nil {
		out.TotalScores = make(map[string]int, len(in.TotalScores))
		for k, v := range in.TotalScores {
			out.TotalScores[k] = v
		}
	}
	return &out
}

