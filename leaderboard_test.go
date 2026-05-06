package main

import "testing"

func TestRecordAddsAndSorts(t *testing.T) {
	lb := &Leaderboard{}

	lb.record(Beginner, 30)
	lb.record(Beginner, 10)
	lb.record(Beginner, 20)

	got := lb.entries(Beginner)
	if len(got) != 3 {
		t.Fatalf("want 3 entries, got %d", len(got))
	}
	if got[0].Seconds != 10 || got[1].Seconds != 20 || got[2].Seconds != 30 {
		t.Errorf("entries not sorted ascending: %+v", got)
	}
}

func TestRecordTrimsToMax(t *testing.T) {
	lb := &Leaderboard{}
	for i := 1; i <= maxScoresPerDifficulty+3; i++ {
		lb.record(Expert, i*10)
	}

	got := lb.entries(Expert)
	if len(got) != maxScoresPerDifficulty {
		t.Errorf("want %d entries, got %d", maxScoresPerDifficulty, len(got))
	}
	if got[0].Seconds != 10 {
		t.Errorf("fastest should still be 10, got %d", got[0].Seconds)
	}
	if got[len(got)-1].Seconds != maxScoresPerDifficulty*10 {
		t.Errorf("last kept should be %d, got %d", maxScoresPerDifficulty*10, got[len(got)-1].Seconds)
	}
}

func TestRecordReturnsRankAndBest(t *testing.T) {
	lb := &Leaderboard{}

	rank, isBest := lb.record(Intermediate, 50)
	if rank != 1 || !isBest {
		t.Errorf("first record: rank=%d isBest=%v, want 1 true", rank, isBest)
	}

	rank, isBest = lb.record(Intermediate, 30)
	if rank != 1 || !isBest {
		t.Errorf("faster record: rank=%d isBest=%v, want 1 true", rank, isBest)
	}

	rank, isBest = lb.record(Intermediate, 70)
	if rank != 3 || isBest {
		t.Errorf("slower record: rank=%d isBest=%v, want 3 false", rank, isBest)
	}
}

func TestBestEmpty(t *testing.T) {
	lb := &Leaderboard{}

	got, ok := lb.best(Beginner)
	if ok {
		t.Error("empty leaderboard should report ok=false")
	}
	if got != 0 {
		t.Errorf("empty leaderboard should return 0, got %d", got)
	}
}

func TestBestAfterRecords(t *testing.T) {
	lb := &Leaderboard{}
	lb.record(Beginner, 100)
	lb.record(Beginner, 50)
	lb.record(Beginner, 75)

	got, ok := lb.best(Beginner)
	if !ok {
		t.Error("populated leaderboard should report ok=true")
	}
	if got != 50 {
		t.Errorf("best should be 50, got %d", got)
	}
}

func TestEntriesIsolatedPerDifficulty(t *testing.T) {
	lb := &Leaderboard{}
	lb.record(Beginner, 10)
	lb.record(Intermediate, 100)
	lb.record(Expert, 500)

	if len(lb.entries(Beginner)) != 1 || lb.entries(Beginner)[0].Seconds != 10 {
		t.Error("Beginner should only have its own record")
	}
	if len(lb.entries(Intermediate)) != 1 || lb.entries(Intermediate)[0].Seconds != 100 {
		t.Error("Intermediate should only have its own record")
	}
	if len(lb.entries(Expert)) != 1 || lb.entries(Expert)[0].Seconds != 500 {
		t.Error("Expert should only have its own record")
	}
}
