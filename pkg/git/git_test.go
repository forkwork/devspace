package git

import (
	"testing"

	"gotest.tools/assert"
	"gotest.tools/assert/cmp"
)

type testCaseNormalizeRepository struct {
	in                  string
	expectedPRReference string
	expectedRepo        string
	expectedBranch      string
	expectedCommit      string
	expectedSubpath     string
}

type testCaseGetBranchNameForPR struct {
	in             string
	expectedBranch string
}

func TestNormalizeRepository(t *testing.T) {
	testCases := []testCaseNormalizeRepository{
		{
			in:                  "ssh://dev.khulnasoft.com.git",
			expectedRepo:        "ssh://dev.khulnasoft.com.git",
			expectedPRReference: "",
			expectedBranch:      "",
			expectedCommit:      "",
			expectedSubpath:     "",
		},
		{
			in:                  "ssh://git@dev.khulnasoft.com.git",
			expectedRepo:        "ssh://git@dev.khulnasoft.com.git",
			expectedPRReference: "",
			expectedBranch:      "",
			expectedCommit:      "",
			expectedSubpath:     "",
		},
		{
			in:                  "git@github.com:khulnasoftdevspace-without-branch.git",
			expectedRepo:        "git@github.com:khulnasoftdevspace-without-branch.git",
			expectedPRReference: "",
			expectedBranch:      "",
			expectedCommit:      "",
			expectedSubpath:     "",
		},
		{
			in:                  "https://dev.khulnasoft.com.git",
			expectedRepo:        "https://dev.khulnasoft.com.git",
			expectedPRReference: "",
			expectedBranch:      "",
			expectedCommit:      "",
			expectedSubpath:     "",
		},
		{
			in:                  "dev.khulnasoft.com.git",
			expectedRepo:        "https://dev.khulnasoft.com.git",
			expectedPRReference: "",
			expectedBranch:      "",
			expectedCommit:      "",
			expectedSubpath:     "",
		},
		{
			in:                  "dev.khulnasoft.com.git@test-branch",
			expectedRepo:        "https://dev.khulnasoft.com.git",
			expectedPRReference: "",
			expectedBranch:      "test-branch",
			expectedCommit:      "",
			expectedSubpath:     "",
		},
		{
			in:                  "git@github.com:khulnasoftdevspace-with-branch.git@test-branch",
			expectedRepo:        "git@github.com:khulnasoftdevspace-with-branch.git",
			expectedPRReference: "",
			expectedBranch:      "test-branch",
			expectedCommit:      "",
			expectedSubpath:     "",
		},
		{
			in:                  "git@github.com:khulnasoftdevspace-with-branch.git@test_branch",
			expectedRepo:        "git@github.com:khulnasoftdevspace-with-branch.git",
			expectedPRReference: "",
			expectedBranch:      "test_branch",
			expectedCommit:      "",
			expectedSubpath:     "",
		},
		{
			in:                  "ssh://git@github.com:khulnasoftdevspace.git@test_branch",
			expectedRepo:        "ssh://git@github.com:khulnasoftdevspace.git",
			expectedPRReference: "",
			expectedBranch:      "test_branch",
			expectedCommit:      "",
			expectedSubpath:     "",
		},
		{
			in:                  "dev.khulnasoft.com-without-protocol-with-slash.git@user/branch",
			expectedRepo:        "https://dev.khulnasoft.com-without-protocol-with-slash.git",
			expectedPRReference: "",
			expectedBranch:      "user/branch",
			expectedCommit:      "",
			expectedSubpath:     "",
		},
		{
			in:                  "git@github.com:khulnasoftdevspace-with-slash.git@user/branch",
			expectedRepo:        "git@github.com:khulnasoftdevspace-with-slash.git",
			expectedPRReference: "",
			expectedBranch:      "user/branch",
			expectedCommit:      "",
			expectedSubpath:     "",
		},
		{
			in:                  "dev.khulnasoft.com.git@sha256:905ffb0",
			expectedRepo:        "https://dev.khulnasoft.com.git",
			expectedPRReference: "",
			expectedBranch:      "",
			expectedCommit:      "905ffb0",
			expectedSubpath:     "",
		},
		{
			in:                  "git@github.com:khulnasoftdevspace.git@sha256:905ffb0",
			expectedRepo:        "git@github.com:khulnasoftdevspace.git",
			expectedPRReference: "",
			expectedBranch:      "",
			expectedCommit:      "905ffb0",
			expectedSubpath:     "",
		},
		{
			in:                  "dev.khulnasoft.com.git@pull/996/head",
			expectedRepo:        "https://dev.khulnasoft.com.git",
			expectedPRReference: "pull/996/head",
			expectedBranch:      "",
			expectedCommit:      "",
			expectedSubpath:     "",
		},
		{
			in:                  "git@github.com:khulnasoftdevspace.git@pull/996/head",
			expectedRepo:        "git@github.com:khulnasoftdevspace.git",
			expectedPRReference: "pull/996/head",
			expectedBranch:      "",
			expectedCommit:      "",
			expectedSubpath:     "",
		},
		{
			in:                  "dev.khulnasoft.com-without-protocol-with-slash.git@subpath:/test/path",
			expectedRepo:        "https://dev.khulnasoft.com-without-protocol-with-slash.git",
			expectedPRReference: "",
			expectedBranch:      "",
			expectedCommit:      "",
			expectedSubpath:     "/test/path",
		},
		{
			in:                  "dev.khulnasoft.com-without-protocol-with-slash.git@subpath:/test/path/",
			expectedRepo:        "https://dev.khulnasoft.com-without-protocol-with-slash.git",
			expectedPRReference: "",
			expectedBranch:      "",
			expectedCommit:      "",
			expectedSubpath:     "/test/path",
		},
		{
			in:                  "https://my_prefix@dev.khulnasoft.com.git@test-branch",
			expectedRepo:        "https://my_prefix@dev.khulnasoft.com.git",
			expectedPRReference: "",
			expectedBranch:      "test-branch",
			expectedCommit:      "",
			expectedSubpath:     "",
		},
		{
			in:                  "https://test@dev.azure.com/org/project/_git/repo@dev",
			expectedRepo:        "https://test@dev.azure.com/org/project/_git/repo",
			expectedPRReference: "",
			expectedBranch:      "dev",
			expectedCommit:      "",
			expectedSubpath:     "",
		},
		{
			in:                  "https://test@dev.azure.com/org/project/_git/repo@sha256:905ffb0",
			expectedRepo:        "https://test@dev.azure.com/org/project/_git/repo",
			expectedPRReference: "",
			expectedBranch:      "",
			expectedCommit:      "905ffb0",
			expectedSubpath:     "",
		},
		{
			in:                  "git@ssh.dev.azure.com:v3/org/project/repo@dev",
			expectedRepo:        "git@ssh.dev.azure.com:v3/org/project/repo",
			expectedPRReference: "",
			expectedBranch:      "dev",
			expectedCommit:      "",
			expectedSubpath:     "",
		},
		{
			in:                  "file:///workspace/projects/project",
			expectedRepo:        "file:///workspace/projects/project",
			expectedPRReference: "",
			expectedBranch:      "",
			expectedCommit:      "",
			expectedSubpath:     "",
		},
		{
			in:                  "file:///workspace/projects/project@dev",
			expectedRepo:        "file:///workspace/projects/project",
			expectedPRReference: "",
			expectedBranch:      "dev",
			expectedCommit:      "",
			expectedSubpath:     "",
		},
		{
			in:                  "file:///workspace/projects/project@sha256:905ffb0",
			expectedRepo:        "file:///workspace/projects/project",
			expectedPRReference: "",
			expectedBranch:      "",
			expectedCommit:      "905ffb0",
			expectedSubpath:     "",
		},
		{
			in:                  "file:///workspace/projects/project@subpath:/test/path",
			expectedRepo:        "file:///workspace/projects/project",
			expectedPRReference: "",
			expectedBranch:      "",
			expectedCommit:      "",
			expectedSubpath:     "/test/path",
		},
	}

	for _, testCase := range testCases {
		outRepo, outPRReference, outBranch, outCommit, outSubpath := NormalizeRepository(testCase.in)
		assert.Check(t, cmp.Equal(testCase.expectedRepo, outRepo))
		assert.Check(t, cmp.Equal(testCase.expectedPRReference, outPRReference))
		assert.Check(t, cmp.Equal(testCase.expectedBranch, outBranch))
		assert.Check(t, cmp.Equal(testCase.expectedCommit, outCommit))
		assert.Check(t, cmp.Equal(testCase.expectedSubpath, outSubpath))
	}
}

func TestGetBranchNameForPRReference(t *testing.T) {
	testCases := []testCaseGetBranchNameForPR{
		{
			in:             "pull/996/head",
			expectedBranch: "PR996",
		},
		{
			in:             "pull/abc/head",
			expectedBranch: "pull/abc/head",
		},
	}

	for _, testCase := range testCases {
		outBranch := GetBranchNameForPR(testCase.in)
		assert.Check(t, cmp.Equal(testCase.expectedBranch, outBranch))
	}
}
