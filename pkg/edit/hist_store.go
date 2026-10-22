package edit

import (
	"sync"

	"src.elv.sh/pkg/cli/histutil"
	"src.elv.sh/pkg/store/storedefs"
)

// A wrapper of histutil.Store that is concurrency-safe and supports an
// additional FastForward method.
type histStore struct {
	m  sync.Mutex
	db storedefs.Store
	hs histutil.Store

	// Deduplicated command texts, oldest to newest, used for autosuggestion.
	// Loaded once and kept up to date by AddCmd, so that autosuggestion
	// doesn't fetch the entire history from the daemon on every keystroke.
	suggest       []string
	suggestLoaded bool
}

func newHistStore(db storedefs.Store) (*histStore, error) {
	hs, err := histutil.NewHybridStore(db)
	return &histStore{db: db, hs: hs}, err
}

func (s *histStore) AddCmd(cmd storedefs.Cmd) (int, error) {
	s.m.Lock()
	defer s.m.Unlock()
	if s.suggestLoaded {
		s.suggest = append(s.suggest, cmd.Text)
	}
	return s.hs.AddCmd(cmd)
}

// SuggestCandidates returns deduplicated command texts, oldest to newest. The
// returned slice must not be modified.
func (s *histStore) SuggestCandidates() []string {
	s.m.Lock()
	defer s.m.Unlock()
	if !s.suggestLoaded {
		cmds, err := s.hs.AllCmds()
		if err != nil {
			return nil
		}
		seen := make(map[string]struct{}, len(cmds)/4)
		uniq := make([]string, 0, len(cmds)/4)
		for i := len(cmds) - 1; i >= 0; i-- {
			if _, ok := seen[cmds[i].Text]; !ok {
				seen[cmds[i].Text] = struct{}{}
				uniq = append(uniq, cmds[i].Text)
			}
		}
		for i, j := 0, len(uniq)-1; i < j; i, j = i+1, j-1 {
			uniq[i], uniq[j] = uniq[j], uniq[i]
		}
		s.suggest, s.suggestLoaded = uniq, true
	}
	return s.suggest[:len(s.suggest):len(s.suggest)]
}

// AllCmds returns a slice of all interactive commands in oldest to newest order.
func (s *histStore) AllCmds() ([]storedefs.Cmd, error) {
	s.m.Lock()
	defer s.m.Unlock()
	return s.hs.AllCmds()
}

func (s *histStore) Cursor(prefix string) histutil.Cursor {
	s.m.Lock()
	defer s.m.Unlock()
	return cursor{&s.m, histutil.NewDedupCursor(s.hs.Cursor(prefix))}
}

func (s *histStore) FastForward() error {
	s.m.Lock()
	defer s.m.Unlock()
	hs, err := histutil.NewHybridStore(s.db)
	s.hs = hs
	s.suggest, s.suggestLoaded = nil, false
	return err
}

type cursor struct {
	m *sync.Mutex
	c histutil.Cursor
}

func (c cursor) Prev() {
	c.m.Lock()
	defer c.m.Unlock()
	c.c.Prev()
}

func (c cursor) Next() {
	c.m.Lock()
	defer c.m.Unlock()
	c.c.Next()
}

func (c cursor) Get() (storedefs.Cmd, error) {
	c.m.Lock()
	defer c.m.Unlock()
	return c.c.Get()
}
