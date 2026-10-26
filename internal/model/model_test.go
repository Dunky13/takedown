package model

import "testing"

func TestNamespaceBindingAndUnsafePaths(t *testing.T) {
	p := Part{ExportID: "first", Index: 1, Filename: "archive.zip", Size: 10}
	p.RelativePath = NamespacePath(p.ExportID, p.Filename)
	if err := p.Validate(); err != nil {
		t.Fatal(err)
	}
	if p.RelativePath == NamespacePath("second", p.Filename) {
		t.Fatal("different exports share output namespace")
	}
	for _, bad := range []string{"../archive.zip", "exports/arbitrary/archive.zip", NamespacePath("second", p.Filename)} {
		p.RelativePath = bad
		if p.Validate() == nil {
			t.Fatalf("unbound path accepted: %s", bad)
		}
	}
}
