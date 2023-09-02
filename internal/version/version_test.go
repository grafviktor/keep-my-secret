package version

import (
	"testing"
)

func TestBuildVersion(t *testing.T) {
	type testValues struct {
		buildVersion string
		buildDate    string
		buildCommit  string
	}

	tests := []struct {
		name string
		arg  testValues
		want testValues
	}{
		{
			name: "Set() does not accept empty values",
			arg: testValues{
				buildVersion: "",
				buildDate:    "",
				buildCommit:  "",
			},
			want: testValues{
				buildVersion: "N/A",
				buildDate:    "N/A",
				buildCommit:  "N/A",
			},
		},
		{
			name: "Set() Build Version",
			arg: testValues{
				buildVersion: "v0.9.9",
				buildDate:    "",
				buildCommit:  "",
			},
			want: testValues{
				buildVersion: "v0.9.9",
				buildDate:    "N/A",
				buildCommit:  "N/A",
			},
		},
		{
			name: "Set() Build Date",
			arg: testValues{
				buildVersion: "v0.9.9",
				buildDate:    "2099-01-01",
				buildCommit:  "",
			},
			want: testValues{
				buildVersion: "v0.9.9",
				buildDate:    "2099-01-01",
				buildCommit:  "N/A",
			},
		},
		{
			name: "Set() Build Commit",
			arg: testValues{
				buildVersion: "v0.9.9",
				buildDate:    "2099-01-01",
				buildCommit:  "f005ba11",
			},
			want: testValues{
				buildVersion: "v0.9.9",
				buildDate:    "2099-01-01",
				buildCommit:  "f005ba11",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Set(tt.arg.buildVersion, tt.arg.buildDate, tt.arg.buildCommit)

			if bi.buildNumber != BuildVersion() ||
				bi.buildDate != BuildDate() ||
				bi.buildCommit != BuildCommit() {
				t.Errorf("Actual args %v, expected %v", bi, tt.want)
			}
		})
	}
}
