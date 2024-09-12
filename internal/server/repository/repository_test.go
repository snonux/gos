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

	for _, entry := range makeEntries(t) {
		t.Run(entry.ID, func(t *testing.T) {
			_ = repo.persist(entry)
			entGot, err := repo.Get(entry.ID)
			if err != nil {
				t.Error(err)
			}
			if !entGot.Equals(entry) {
				t.Error("expected to get", entry, "but got", entGot)
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
	for _, entry := range entries {
		bytes, _ := entry.JSONMarshal()
		_ = repo.fs.WriteFile(repo.entryPath(entry), bytes)
	}

	// Load entries from VFS into the repo
	if err := repo.load(); err != nil {
		t.Error(err)
	}

	for _, entry := range entries {
		t.Run(entry.ID, func(t *testing.T) {
			entGot, err := repo.Get(entry.ID)
			if err != nil {
				t.Error(err)
			}
			if !entGot.Equals(entry) {
				t.Error("expected to get", entry, "but got", entGot)
			}
		})
	}
}

func TestRepositoryList(t *testing.T) {
	t.Parallel()

	fs := make(vfs.MemoryFS)
	repo := newRepository(server.ServerConfig{DataDir: "./data"}, fs)
	entries := makeEntries(t)

	for _, entry := range entries {
		_ = repo.persist(entry)
	}

	pairs, _ := repo.List()
	if len(entries) != len(pairs) {
		t.Error("expected as many entries as pairs")
	}

	for _, entry := range entries {
		var found bool
		for _, pair := range pairs {
			if entry.ID == pair.ID && entry.Checksum() == pair.Checksum {
				found = true
				t.Log("entry matches pair", entry, pair)
				break
			}
		}
		if !found {
			t.Error("could not find entry", entry, "in", pairs)
		}
	}
}

func TestRepositoryHasSameEntry(t *testing.T) {
	t.Parallel()

	fs := make(vfs.MemoryFS)
	repo := newRepository(server.ServerConfig{DataDir: "./data"}, fs)
	entry, _ := makeAnEntry()
	_ = repo.persist(entry)

	pair := entryPair{entry.ID, entry.Checksum()}
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
	entry1, _ := makeAnEntry()
	_ = repo.persist(entry1)

	entry2, _ := makeAnotherEntry()
	// Need to have the same IDs so that the entries will actually be merged
	entry2.ID = entry1.ID
	// Merge a modified entry2 into the repository.
	entry2.Body = "merged"
	entry2.Epoch = 12345
	_ = repo.Merge(entry2)

	pairs, _ := repo.List()
	// Ensuring the merge didn't add a new entry
	if len(pairs) != 1 {
		t.Error("expected exactly one element in the repo but got", pairs)
	}

	entGot, _ := repo.Get(entry1.ID)
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

	entry1, _ := makeAnEntry()
	_ = repo1.persist(entry1)
	entry2, _ := makeAnotherEntry()
	_ = repo2.persist(entry2)

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

	getEntry := func(ctx context.Context, partner, id string, entry *types.Entry) error {
		var (
			entry_ types.Entry
			err    error
		)

		switch partner {
		case "repo1":
			entry_, err = repo1.Get(id)
		case "repo2":
			entry_, err = repo2.Get(id)
		}

		if err != nil {
			return err
		}
		*entry = entry_

		t.Log("got entry", *entry, "from repo", partner)
		return nil
	}

	// Compare both repos, they should now contain the same entries
	compare := func(repo1, repo2 Repository) error {
		pairs, err := repo1.List()
		if err != nil {
			return err
		}

		for _, pair := range pairs {
			entry1, err := repo1.Get(pair.ID)
			if err != nil {
				return err
			}
			entry2, err := repo2.Get(pair.ID)
			if err != nil {
				return err
			}

			t.Log("comparing entries")
			t.Log("entry1", entry1)
			t.Log("entry2", entry2)

			if !entry1.Equals(entry2) {
				return fmt.Errorf("entries entry1 and entry2 don't equal")
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
		entry, err := repo1.Get(entry1.ID)
		if err != nil {
			t.Error(err)
		}

		// Validate the correct test setup
		if entry.IsShared(types.LinkedIn) {
			t.Error("for the test expected LinkedIn not to be shared")
		}

		// Simulate that the entry was shared to LinkedIn social media!
		linkedIn, ok := entry.Shared[types.LinkedIn]
		if !ok {
			t.Error("expected to have a LinkedIn shared entry")
		}
		linkedIn.Is = true
		entry.Shared[types.LinkedIn] = linkedIn

		if err := repo1.Merge(entry); err != nil {
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

func TestRepositoryNext(t *testing.T) {
	t.Parallel()

	fs := make(vfs.MemoryFS)
	repo := newRepository(server.ServerConfig{
		DataDir: "./data",
		SocialPlatformsEnabled: []types.PlatformName{
			types.LinkedIn, types.Mastodon, types.Textfile,
		},
	}, fs)
	entries := makeEntries(t)

	for _, entry := range entries {
		_ = repo.persist(entry)
	}

	if entry, ok := repo.Next(types.Mastodon); ok {
		t.Error("expected no Mastodon entry to be found", entry)
	}

	if _, ok := repo.Next(types.LinkedIn); !ok {
		t.Error("expected an unshared LinkedIn entry to be found")
	}

	if _, ok := repo.Next(types.Textfile); !ok {
		t.Error("expected an unshared Textfile entry to be found")
	}
}

func makeEntries(t *testing.T) []types.Entry {
	entry1, err := makeAnEntry()
	if err != nil {
		t.Error(err)
	}
	entry2, err := makeAnotherEntry()
	if err != nil {
		t.Error(err)
	}
	return []types.Entry{entry1, entry2}
}

func makeAnEntry() (types.Entry, error) {
	entry := `
		{
			"body": "Body text here",
			"shared": {
				"Mastodon": { "is": true },
				"LinkedIn": { "is": false }
			}
		}
	`
	return types.NewEntry([]byte(entry))
}

func makeAnotherEntry() (types.Entry, error) {
	entry := `
		{
			"body": "Another text here",
			"shared": {
				"Mastodon": { "is": true },
				"LinkedIn": { "is": true },
				"Textfile": { "is": false }
			}
		}
	`
	return types.NewEntry([]byte(entry))
}
