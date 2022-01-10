package definitions

import "testing"

func TestSchemaEmbedding(t *testing.T) {
	version := "3.9.x"
	schema, err := LoadSchema(version)
	if err != nil {
		t.Fatal(err)
	}

	if schema.Version != version {
		t.Fatalf("schema version '%s' doesn not match requested  version '%s'", schema.Version, version)
	}

	if len(schema.Endpoints) < 1 {
		t.Fatal("Expected at least some endpoints in the loaded schema")
	}
}
func TestFeatureEmbedding(t *testing.T) {
	version := "3.9.x"
	set, err := LoadFeatures(version)
	if err != nil {
		t.Fatal(err)
	}

	if set.Version != version {
		t.Fatalf("feature set version '%s' doesn not match requested  version '%s'", set.Version, version)
	}

	if len(set.Features) < 1 {
		t.Fatal("Expected at least some features in the loaded featureset")
	}
}

func TestAvailableSchemaVersions(t *testing.T) {
	versions := AvailableSchemaVersions()

	if len(versions) < 2 {
		t.Errorf("expected at least 2 schema versions")
	}

}

func BenchmarkIsValidVersionSpec(b *testing.B) {
	testcases := []struct {
		version  string
		expected bool
	}{
		{"", false},
		{".", false},
		{"1.", false},
		{".1", false},
		{"1.2.", false},
		{".1.2.", false},
		{"1.x.4", false},
		{"x.1.4", false},
		{"xx", false},
		{"1", true},
		{"112", true},
		{"x", true},
		{"1.1", true},
		{"1.2.3.4", true},
		{"1.2.3.4.5", true},
		{"1.2.3", true},
		{"1.4.x", true},
		{"1.x.x", true},
	}

	n := len(testcases)
	b.ResetTimer()

	for i := 0; i < b.N; i += n {
		for _, tc := range testcases {
			if tc.expected != IsValidVersionSpec(tc.version) {
				b.Errorf("expected validVersion('%s') to return %v", tc.version, tc.expected)
			}
		}
	}
}

func BenchmarkMatchVersion(b *testing.B) {
	testcases := []struct {
		spec     string
		version  string
		expected bool
	}{
		{"1.2.x", "1.1.b", false},
		{"", ".", false},
		{"x.x.x", "x.x.", false},
		{"x.x.x", "x.x", false},
		{"x.x", "x.x.x", false},
		{"x.x.x", "1.2.b", false},
		{"1.2.3", "1.2.4", false},
		{"4.2.3", "1.2.3", false},
		{"x.x.x", "x.2.x", false},
		{"x.x", "x.2", false},
		{"x.x", "2.x", true},
		{"x", "x", true},
		{"x", "2", true},
		{"", "", false},
		{"x.x.x", "1.2.3", true},
		{"x.x.x", "2.2.x", true},
		{"x.x.x", "2.x.x", true},
		{"x.x.x", "x.x.x", true},
		{"1.x.x", "1.2.3", true},
		{"1.2.x", "1.2.3", true},
		{"1.2.3", "1.2.3", true},
	}

	n := len(testcases)
	b.ResetTimer()
	for i := 0; i < b.N; i += n {
		for _, tc := range testcases {
			if tc.expected != MatchVersion(tc.spec, tc.version) {
				b.Errorf("expected matchVersion('%s', '%s') to return %v", tc.spec, tc.version, tc.expected)
			}
		}
	}

}
