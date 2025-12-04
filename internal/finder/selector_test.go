package finder

import (
	"testing"
	"time"
)

func TestFilterItems_EmptyQuery(t *testing.T) {
	items := []entry{
		{displayLabel: "2025-12-04  memo.md  [repo1]", absPath: "/path/memo.md"},
		{displayLabel: "2025-12-03  notes.txt  [repo2]", absPath: "/path/notes.txt"},
	}

	m := newSelector(items)
	m.input.SetValue("")
	m.filterItems()

	if len(m.filteredItems) != len(items) {
		t.Errorf("Expected %d items, got %d", len(items), len(m.filteredItems))
	}
}

func TestFilterItems_SingleKeyword(t *testing.T) {
	items := []entry{
		{displayLabel: "2025-12-04  memo.md  [repo1]", absPath: "/path/memo.md"},
		{displayLabel: "2025-12-03  notes.txt  [repo2]", absPath: "/path/notes.txt"},
		{displayLabel: "2025-12-02  script.sh  [repo3]", absPath: "/path/script.sh"},
	}

	m := newSelector(items)
	m.input.SetValue("memo")
	m.filterItems()

	if len(m.filteredItems) != 1 {
		t.Errorf("Expected 1 item, got %d", len(m.filteredItems))
	}

	if len(m.filteredItems) > 0 && m.filteredItems[0].absPath != "/path/memo.md" {
		t.Errorf("Expected /path/memo.md, got %s", m.filteredItems[0].absPath)
	}
}

func TestFilterItems_ANDSearch(t *testing.T) {
	items := []entry{
		{displayLabel: "2025-12-04  memo.md  [my-project]", absPath: "/path/memo.md"},
		{displayLabel: "2025-12-03  notes.txt  [my-notes]", absPath: "/path/notes.txt"},
		{displayLabel: "2025-12-02  readme.md  [other-project]", absPath: "/path/readme.md"},
		{displayLabel: "2025-12-01  todo.md  [my-project]", absPath: "/path/todo.md"},
	}

	m := newSelector(items)
	m.input.SetValue("md my")
	m.filterItems()

	// Should match items containing both "md" AND "my"
	// Expected: memo.md (has "md" and "my-project")
	//           todo.md (has "md" and "my-project")
	if len(m.filteredItems) != 2 {
		t.Errorf("Expected 2 items, got %d", len(m.filteredItems))
		for _, item := range m.filteredItems {
			t.Logf("Matched: %s", item.displayLabel)
		}
	}
}

func TestFilterItems_CaseInsensitive(t *testing.T) {
	items := []entry{
		{displayLabel: "2025-12-04  README.md  [repo1]", absPath: "/path/README.md"},
		{displayLabel: "2025-12-03  notes.txt  [repo2]", absPath: "/path/notes.txt"},
	}

	m := newSelector(items)
	m.input.SetValue("readme")
	m.filterItems()

	if len(m.filteredItems) != 1 {
		t.Errorf("Expected 1 item, got %d", len(m.filteredItems))
	}

	if len(m.filteredItems) > 0 && m.filteredItems[0].absPath != "/path/README.md" {
		t.Errorf("Expected /path/README.md, got %s", m.filteredItems[0].absPath)
	}
}

func TestFilterItems_NoMatch(t *testing.T) {
	items := []entry{
		{displayLabel: "2025-12-04  memo.md  [repo1]", absPath: "/path/memo.md"},
		{displayLabel: "2025-12-03  notes.txt  [repo2]", absPath: "/path/notes.txt"},
	}

	m := newSelector(items)
	m.input.SetValue("nonexistent")
	m.filterItems()

	if len(m.filteredItems) != 0 {
		t.Errorf("Expected 0 items, got %d", len(m.filteredItems))
	}
}

func TestFilterItems_PreservesOrder(t *testing.T) {
	// Create items with specific timestamps
	now := time.Now()
	items := []entry{
		{displayLabel: "2025-12-04  file3.txt  [repo]", absPath: "/path/file3.txt", modTime: now},
		{displayLabel: "2025-12-03  file2.txt  [repo]", absPath: "/path/file2.txt", modTime: now.Add(-24 * time.Hour)},
		{displayLabel: "2025-12-02  file1.txt  [repo]", absPath: "/path/file1.txt", modTime: now.Add(-48 * time.Hour)},
	}

	m := newSelector(items)
	m.input.SetValue("file")
	m.filterItems()

	// Should preserve the original order (sorted by modTime descending)
	if len(m.filteredItems) != 3 {
		t.Errorf("Expected 3 items, got %d", len(m.filteredItems))
	}

	if len(m.filteredItems) >= 3 {
		if m.filteredItems[0].absPath != "/path/file3.txt" {
			t.Errorf("Expected first item to be file3.txt, got %s", m.filteredItems[0].absPath)
		}
		if m.filteredItems[1].absPath != "/path/file2.txt" {
			t.Errorf("Expected second item to be file2.txt, got %s", m.filteredItems[1].absPath)
		}
		if m.filteredItems[2].absPath != "/path/file1.txt" {
			t.Errorf("Expected third item to be file1.txt, got %s", m.filteredItems[2].absPath)
		}
	}
}
