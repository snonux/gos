package repository

import (
	"context"
	"testing"

	"codeberg.org/snonux/gos/internal/config/server"
	"codeberg.org/snonux/gos/internal/types"
	"codeberg.org/snonux/gos/internal/vfs"
)

func TestRepositoryPutGet(t *testing.T) {
	t.Parallel()

	fs := make(vfs.MemoryFS)
	repo := newRepository(server.ServerConfig{DataDir: "./data"}, fs)

	for _, ent := range makeEntries(t) {
		t.Run(ent.ID, func(t *testing.T) {
			_ = repo.put(ent)
			entGot, ok := repo.Get(ent.ID)
			if !ok {
				t.Errorf("could not find entry with id %s in repo", ent.ID)
			}
			if !entGot.Equals(ent) {
				t.Error("expected to get", ent, "but got", entGot)
			}
		})
	}
}

func TestRepositoryLoad(t *testing.T) {
	t.Parallel()

	fs := make(vfs.MemoryFS)
	repo := newRepository(server.ServerConfig{DataDir: "./data"}, fs)
	entries := makeEntries(t)

	// Write entries into the VFS
	for _, ent := range entries {
		bytes, _ := ent.JSONMarshal()
		_ = repo.fs.WriteFile(repo.entryPath(ent), bytes)
	}

	// Load entries from VFS into the repo
	if err := repo.load(); err != nil {
		t.Error(err)
	}

	for _, ent := range entries {
		t.Run(ent.ID, func(t *testing.T) {
			entGot, ok := repo.Get(ent.ID)
			if !ok {
				t.Errorf("could not find entry with id %s in repo", ent.ID)
			}
			if !entGot.Equals(ent) {
				t.Error("expected to get", ent, "but got", entGot)
			}
		})
	}
}

func TestRepositoryList(t *testing.T) {
	t.Parallel()

	fs := make(vfs.MemoryFS)
	repo := newRepository(server.ServerConfig{DataDir: "./data"}, fs)
	entries := makeEntries(t)

	for _, ent := range entries {
		_ = repo.put(ent)
	}

	pairs, _ := repo.List()
	if len(entries) != len(pairs) {
		t.Error("expected as many entries as pairs")
	}

	for _, ent := range entries {
		var found bool
		for _, pair := range pairs {
			if ent.ID == pair.ID && ent.Checksum() == pair.Checksum {
				found = true
				t.Log("entry matches pair", ent, pair)
				break
			}
		}
		if !found {
			t.Error("could not find entry", ent, "in", pairs)
		}
	}
}

func TestRepositoryHasSameEntry(t *testing.T) {
	t.Parallel()

	fs := make(vfs.MemoryFS)
	repo := newRepository(server.ServerConfig{DataDir: "./data"}, fs)
	ent, _ := makeAnEntry()
	_ = repo.put(ent)

	pair := entryPair{ent.ID, ent.Checksum()}
	if !repo.hasSameEntry(pair) {
		t.Error("repo does not contain entry corresponding to pair", pair)
	}

	pair = entryPair{"nonexistent", "nonexistent"}
	if repo.hasSameEntry(pair) {
		t.Error("repo does contain entry corresponding to pair", pair, "but that should not be")
	}
}

func TestRepositoryMerge(t *testing.T) {
	t.Parallel()

	fs := make(vfs.MemoryFS)
	repo := newRepository(server.ServerConfig{DataDir: "./data"}, fs)
	ent1, _ := makeAnEntry()
	_ = repo.put(ent1)

	ent2, _ := makeAnotherEntry()
	// Need to have the same IDs so that the entries will actually be merged
	ent2.ID = ent1.ID
	// Merge a modified ent2 into the repository.
	ent2.Body = "merged"
	ent2.Epoch = 12345
	_ = repo.Merge(ent2)

	pairs, _ := repo.List()
	// Ensuring the merge didn't add a new entry
	if len(pairs) != 1 {
		t.Error("expected exactly one element in the repo but got", pairs)
	}

	entGot, _ := repo.Get(ent1.ID)
	if entGot.Body != "merged" {
		t.Error("unexpected body", entGot.Body)
	}
	if entGot.Epoch != 12345 {
		t.Error("unexpected epoch", entGot.Epoch)
	}
}

// TODO: Finish implementing this test
func TestRepositoryMergeFromPartner(t *testing.T) {
	fs1 := make(vfs.MemoryFS)
	repo1 := newRepository(server.ServerConfig{DataDir: "./data1"}, fs1)
	fs2 := make(vfs.MemoryFS)
	repo2 := newRepository(server.ServerConfig{DataDir: "./data2"}, fs2)

	ent1, _ := makeAnEntry()
	_ = repo1.put(ent1)
	ent2, _ := makeAnotherEntry()
	_ = repo2.put(ent2)

	getPair := func(ctx context.Context, partner string, pairs *[]entryPair) error {
		var partnerRepo Repository

		switch partner {
		case "repo1":
			partnerRepo = repo2
		case "repo2":
			partnerRepo = repo1
		}

		pairs_, err := partnerRepo.List()
		if err != nil {
			return err
		}
		*pairs = pairs_
		t.Log("got pairs", *pairs, "from repo", partner)

		return nil
	}

	getEntry := func(ctx context.Context, partner, id string, ent *types.Entry) error {
		return nil
		/*
			uri := fmt.Sprintf("%s/get?id=%s", partner, id)
			return easyhttp.GetData(ctx, uri, r.conf.APIKey, ent)
		*/
	}

	if err := repo1.mergeFromPartner(context.Background(), "repo2", getPair, getEntry); err != nil {
		t.Error(err)
	}
	if err := repo2.mergeFromPartner(context.Background(), "repo1", getPair, getEntry); err != nil {
		t.Error(err)
	}

}

func makeEntries(t *testing.T) []types.Entry {
	ent1, err := makeAnEntry()
	if err != nil {
		t.Error(err)
	}
	ent2, err := makeAnotherEntry()
	if err != nil {
		t.Error(err)
	}
	return []types.Entry{ent1, ent2}
}

func makeAnEntry() (types.Entry, error) {
	entry := `
		{
			"Body": "Body text here",
			"Shared": [
				{ "Name": "Foo", "Is": true },
				{ "Name": "Bar", "Is": false }
			]
		}
	`
	return types.NewEntry([]byte(entry))
}

func makeAnotherEntry() (types.Entry, error) {
	entry := `
		{
			"Body": "Another text here",
			"Shared": [
				{ "Name": "Foo", "Is": true },
				{ "Name": "Bar", "Is": true },
				{ "Name": "Baz", "Is": false }
			]
		}
	`
	return types.NewEntry([]byte(entry))
}
