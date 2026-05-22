package problem_statement

import (
	"context"
	"io"
	"strings"
	"testing"
	"time"

	"contestr/internal/repository"
	"contestr/pkg/regatta"
)

type memRepo struct {
	docs []repository.ProblemStatement
}

func (m *memRepo) Upsert(ctx context.Context, doc repository.ProblemStatement) error {
	for i, d := range m.docs {
		if d.ContestID == doc.ContestID && d.ProblemCode == doc.ProblemCode {
			m.docs[i] = doc
			return nil
		}
	}
	m.docs = append(m.docs, doc)
	return nil
}

func (m *memRepo) Get(ctx context.Context, contestID int, problemCode string) (*repository.ProblemStatement, error) {
	for _, d := range m.docs {
		if d.ContestID == contestID && d.ProblemCode == problemCode {
			return &d, nil
		}
	}
	return nil, nil
}

func (m *memRepo) ListByContest(ctx context.Context, contestID int) ([]repository.ProblemStatement, error) {
	var out []repository.ProblemStatement
	for _, d := range m.docs {
		if d.ContestID == contestID {
			out = append(out, d)
		}
	}
	return out, nil
}

func (m *memRepo) Delete(ctx context.Context, contestID int, problemCode string) error {
	filtered := m.docs[:0]
	for _, d := range m.docs {
		if d.ContestID == contestID && d.ProblemCode == problemCode {
			continue
		}
		filtered = append(filtered, d)
	}
	m.docs = filtered
	return nil
}

func (m *memRepo) DeleteByContest(ctx context.Context, contestID int) error {
	filtered := m.docs[:0]
	for _, d := range m.docs {
		if d.ContestID == contestID {
			continue
		}
		filtered = append(filtered, d)
	}
	m.docs = filtered
	return nil
}

type fakeStore struct {
	urlPrefix string
}

func (f *fakeStore) PutObject(ctx context.Context, objectKey string, body io.Reader, size int64) error {
	return nil
}

func (f *fakeStore) DeleteObject(ctx context.Context, objectKey string) error {
	return nil
}

func (f *fakeStore) PublicURL(objectKey string) string {
	return f.urlPrefix + "/" + objectKey
}

type fakeTours struct {
	tours []regatta.Tour
}

func (f *fakeTours) FindByContestID(ctx context.Context, contestID int) ([]regatta.Tour, error) {
	return f.tours, nil
}

func TestListPublicStatements_onlyStartedRounds(t *testing.T) {
	repo := &memRepo{docs: []repository.ProblemStatement{
		{ContestID: 1, ProblemCode: "1A", ObjectKey: "k1", UploadedAt: time.Now()},
		{ContestID: 1, ProblemCode: "2A", ObjectKey: "k2", UploadedAt: time.Now()},
	}}
	svc := NewService(&fakeStore{urlPrefix: "https://storage.example/bucket"}, repo, &fakeTours{
		tours: []regatta.Tour{{ContestID: 1, Round: 1, IsPause: false}},
	})

	out, err := svc.ListPublicStatements(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 1 {
		t.Fatalf("statements = %+v, want only 1A", out)
	}
	if !strings.Contains(out["1A"], "https://") {
		t.Fatalf("url = %q", out["1A"])
	}
	if _, ok := out["2A"]; ok {
		t.Fatal("2A must not be released before tour 2")
	}
}
