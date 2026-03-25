package yamlx

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveScalarPreservesSpacingAndComments(t *testing.T) {
	t.Parallel()

	const original = "# comment about root\nroot:\n  child: old-value\n\n  sibling:\n    nested: keep\n"
	const expected = "# comment about root\nroot:\n  child: new-value\n\n  sibling:\n    nested: keep\n"

	path := filepath.Join(t.TempDir(), "test.yaml")
	if err := os.WriteFile(path, []byte(original), 0644); err != nil {
		t.Fatalf("write test yaml: %v", err)
	}

	root, err := Load(path)
	if err != nil {
		t.Fatalf("load yaml: %v", err)
	}

	target, err := Set(root, "root.child", "new-value")
	if err != nil {
		t.Fatalf("set yaml value: %v", err)
	}

	if err := SaveScalar(path, target); err != nil {
		t.Fatalf("save yaml value: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read updated yaml: %v", err)
	}

	if string(got) != expected {
		t.Fatalf("unexpected yaml after update:\n--- got ---\n%s\n--- want ---\n%s", got, expected)
	}
}

func TestSaveScalarPreservesQuotedStyle(t *testing.T) {
	t.Parallel()

	const original = "root:\n  child: 'old value'\n  sibling: keep\n"
	const expected = "root:\n  child: 'it''s new'\n  sibling: keep\n"

	path := filepath.Join(t.TempDir(), "test.yaml")
	if err := os.WriteFile(path, []byte(original), 0644); err != nil {
		t.Fatalf("write test yaml: %v", err)
	}

	root, err := Load(path)
	if err != nil {
		t.Fatalf("load yaml: %v", err)
	}

	target, err := Set(root, "root.child", "it's new")
	if err != nil {
		t.Fatalf("set yaml value: %v", err)
	}

	if err := SaveScalar(path, target); err != nil {
		t.Fatalf("save yaml value: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read updated yaml: %v", err)
	}

	if string(got) != expected {
		t.Fatalf("unexpected yaml after update:\n--- got ---\n%s\n--- want ---\n%s", got, expected)
	}
}

func TestSaveScalarDoesNotAddQuotesToPlainStyle(t *testing.T) {
	t.Parallel()

	const original = "root:\n  child: old-value\n  sibling: keep\n"
	const expected = "root:\n  child: true\n  sibling: keep\n"

	path := filepath.Join(t.TempDir(), "test.yaml")
	if err := os.WriteFile(path, []byte(original), 0644); err != nil {
		t.Fatalf("write test yaml: %v", err)
	}

	root, err := Load(path)
	if err != nil {
		t.Fatalf("load yaml: %v", err)
	}

	target, err := Set(root, "root.child", "true")
	if err != nil {
		t.Fatalf("set yaml value: %v", err)
	}

	if err := SaveScalar(path, target); err != nil {
		t.Fatalf("save yaml value: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read updated yaml: %v", err)
	}

	if string(got) != expected {
		t.Fatalf("unexpected yaml after update:\n--- got ---\n%s\n--- want ---\n%s", got, expected)
	}
}

func TestGetValueFromArrayPath(t *testing.T) {
	t.Parallel()

	const original = "items:\n  - name: first\n  - name: second\n"

	path := filepath.Join(t.TempDir(), "test.yaml")
	if err := os.WriteFile(path, []byte(original), 0644); err != nil {
		t.Fatalf("write test yaml: %v", err)
	}

	root, err := Load(path)
	if err != nil {
		t.Fatalf("load yaml: %v", err)
	}

	got, err := Get(root, "items[1].name")
	if err != nil {
		t.Fatalf("get yaml value: %v", err)
	}

	if got != "second" {
		t.Fatalf("unexpected value: got %q want %q", got, "second")
	}
}

func TestSaveScalarUpdatesArrayItem(t *testing.T) {
	t.Parallel()

	const original = "items:\n  - name: first\n  - name: second\n"
	const expected = "items:\n  - name: first\n  - name: updated\n"

	path := filepath.Join(t.TempDir(), "test.yaml")
	if err := os.WriteFile(path, []byte(original), 0644); err != nil {
		t.Fatalf("write test yaml: %v", err)
	}

	root, err := Load(path)
	if err != nil {
		t.Fatalf("load yaml: %v", err)
	}

	target, err := Set(root, "items[1].name", "updated")
	if err != nil {
		t.Fatalf("set yaml value: %v", err)
	}

	if err := SaveScalar(path, target); err != nil {
		t.Fatalf("save yaml value: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read updated yaml: %v", err)
	}

	if string(got) != expected {
		t.Fatalf("unexpected yaml after update:\n--- got ---\n%s\n--- want ---\n%s", got, expected)
	}
}
