package repository

import (
	"context"
	"fmt"
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
			entGot, err := repo.Get(ent.ID)
			if err != nil {
				t.Error(err)
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
			entGot, err := repo.Get(ent.ID)
			if err != nil {
				t.Error(err)
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
		var (
			pairs_ []entryPair
			err    error
		)

		switch partner {
		case "repo1":
			pairs_, err = repo1.List()
		case "repo2":
			pairs_, err = repo2.List()
		}

		if err != nil {
			return err
		}
		*pairs = pairs_

		t.Log("got pairs", *pairs, "from repo", partner)
		return nil
	}

	getEntry := func(ctx context.Context, partner, id string, ent *types.Entry) error {
		var (
			ent_ types.Entry
			err  error
		)

		switch partner {
		case "repo1":
			ent_, err = repo1.Get(id)
		case "repo2":
			ent_, err = repo2.Get(id)
		}

		if err != nil {
			return err
		}
		*ent = ent_

		t.Log("got entry", *ent, "from repo", partner)
		return nil
	}

	// Compare both repos, they should now contain the same entries
	compare := func(repo1, repo2 Repository) error {
		pairs, err := repo1.List()
		if err != nil {
			return err
		}

		for _, pair := range pairs {
			ent1, err := repo1.Get(pair.ID)
			if err != nil {
				return err
			}
			ent2, err := repo2.Get(pair.ID)
			if err != nil {
				return err
			}

			t.Log("comparing entries")
			t.Log("ent1", ent1)
			t.Log("ent2", ent2)

			if !ent1.Equals(ent2) {
				return fmt.Errorf("entries ent1 and ent2 don't equal")
			}
		}

		return nil
	}

	t.Run("Merge entries from repo2 into repo1", func(t *testing.T) {
		if err := repo1.mergeFromPartner(context.Background(), "repo2", getPair, getEntry); err != nil {
			t.Error(err)
		}
		if err := compare(repo2, repo1); err != nil {
			t.Error(err)
		}
	})

	t.Run("Merge entries from repo1 into repo2", func(t *testing.T) {
		if err := repo2.mergeFromPartner(context.Background(), "repo1", getPair, getEntry); err != nil {
			t.Error(err)
		}
		if err := compare(repo1, repo2); err != nil {
			t.Error(err)
		}
	})

	t.Run("Change shared flag and merge to partner", func(t *testing.T) {
		ent, err := repo1.Get(ent1.ID)
		if err != nil {
			t.Error(err)
		}

		// Validate the correct test setup
		if ent.Shared[1].Name != "LinkedIn" || ent.Shared[1].Is != false {
			t.Error("for the test expected LinkedIn not to be shared", ent.Shared[1])
		}

		// Simulate that the entry was shared to LinkedIn social media!
		ent.Shared[1].Is = true
		if err := repo1.Update(ent); err != nil {
			t.Error(err)
		}

		// Before merging, repos should be out of sync.
		if err := compare(repo1, repo2); err == nil {
			t.Log("as expected repos are out of sync", err)
		}

		// Partner is merging the repo.
		if err := repo1.mergeFromPartner(context.Background(), "repo2", getPair, getEntry); err != nil {
			t.Error(err)
		}

		// Still out of sync, as we merged the repos the wrong direction.
		if err := compare(repo1, repo2); err == nil {
			t.Log("as expected repos are out of sync", err)
		}

		// Partner is merging the repo the right direction.
		if err := repo2.mergeFromPartner(context.Background(), "repo1", getPair, getEntry); err != nil {
			t.Error(err)
		}

		// Now, partners should be in sync.
		if err := compare(repo1, repo2); err != nil {
			t.Error(err)
		}
	})
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
				{ "Name": "Mastodon", "Is": true },
				{ "Name": "LinkedIn", "Is": false }
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
				{ "Name": "Mastodon", "Is": true },
				{ "Name": "LinkedIn", "Is": true },
				{ "Name": "foo.zone", "Is": false }
			]
		}
	`
	return types.NewEntry([]byte(entry))
}
