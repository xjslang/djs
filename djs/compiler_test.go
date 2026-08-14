package djs_test

import (
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xjslang/djs/djs"
)

func TestCompiler(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "consecutive defers",
			input: `{
				defer console.log('a')
				defer console.log('b')
			}`,
			expected: "b\na\n",
		},
		{
			name: "separate blocks",
			input: `
			{ defer console.log('a') }
			{ defer console.log('b') }`,
			expected: "a\nb\n",
		},
		{
			name: "nested blocks",
			input: `
			{
				defer console.log('a')
				{ defer console.log('b') }
			}`,
			expected: "b\na\n",
		},
		{
			name: "mixed",
			input: `
			{
				defer console.log('a')
				{ defer console.log('b') }
				defer console.log('c')
			}
			{ defer console.log('d') }`,
			expected: "b\nc\na\nd\n",
		},
	}
	requireNode16OrGreater(t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := djs.Parse(test.input)
			require.NoError(t, err)
			jsCode, err := djs.Compile(result)
			require.NoError(t, err)
			// execute JS code
			cmd := exec.Command("node")
			cmd.Stdin = strings.NewReader(jsCode)
			output, err := cmd.Output()
			require.NoError(t, err)
			require.Equal(t, test.expected, string(output))
		})
	}
}

func TestCompiler_misusedDefer(t *testing.T) {
	input := "defer console.log('aaa')"
	result, err := djs.Parse(input)
	require.NoError(t, err)
	_, err = djs.Compile(result)
	require.Error(t, err, "should panic when using defer outside a block")
}

func requireNode16OrGreater(t *testing.T) {
	t.Helper()
	// test that node is installed
	nodePath, err := exec.LookPath("node")
	require.NoError(t, err)
	// test node version
	cmd := exec.Command(nodePath, "--version")
	output, err := cmd.Output()
	require.NoError(t, err)
	r, err := regexp.Compile(`^v(\d+)`)
	require.NoError(t, err)
	matches := r.FindSubmatch(output)
	require.Len(t, matches, 2)
	version, _ := strconv.Atoi(string(matches[1]))
	if version < 16 {
		t.Fatalf("node v16 or greater is required. Installed version: %s", output)
	}
}
