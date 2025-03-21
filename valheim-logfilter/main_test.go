package main

import (
	"os"
	"regexp"
	"strings"
	"testing"
	"time"
	"unicode/utf8"
)

func TestRemoveLogLine(t *testing.T) {
	t.Parallel()
	result := removeLogLine("", "test log line")
	if !result {
		t.Errorf("removeLogLine with empty command should return true")
	}

	result = removeLogLine("echo test", "test log line")
	if result {
		t.Errorf("removeLogLine with non-empty command should return false")
	}
}

func TestRunHook(t *testing.T) {
	t.Parallel()
	tmpfile, err := os.CreateTemp("", "testhook")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	cmd := "cat > " + tmpfile.Name()
	testLine := "test log line for hook"

	runHook(cmd, testLine)

	time.Sleep(100 * time.Millisecond)

	content, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	expected := testLine + "\n"
	if string(content) != expected {
		t.Errorf("runHook did not write expected content. Got: %q, want: %q", string(content), expected)
	}
}

func TestUTF8Filtering(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name           string
		input          string
		expectedOutput string
	}{
		{
			name:           "Valid UTF-8 string",
			input:          "Hello, this is a valid UTF-8 string",
			expectedOutput: "Hello, this is a valid UTF-8 string",
		},
		{
			name:           "Valid UTF-8 with special characters",
			input:          "Special UTF-8 chars: ñáéíóúü汉字🚀",
			expectedOutput: "Special UTF-8 chars: ñáéíóúü汉字🚀",
		},
		{
			name:           "Invalid UTF-8 sequence",
			input:          string([]byte{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0xc3, 0x28}), // Hello with invalid UTF-8
			expectedOutput: "Hello(",
		},
		{
			name:           "Multiple invalid UTF-8 sequences",
			input:          string([]byte{0x41, 0xc3, 0x28, 0x42, 0xed, 0xa0, 0x80, 0x43}), // A<invalid>B<invalid>C
			expectedOutput: "A(BC",
		},
		{
			name:           "Line with only invalid UTF-8 sequence",
			input:          string([]byte{0xc3, 0x28}), // Just an invalid sequence
			expectedOutput: "(",
		},
		{
			name:           "Invalid UTF-8 at beginning of string",
			input:          string([]byte{0xc3, 0x28, 0x56, 0x61, 0x6c, 0x69, 0x64}), // <invalid>Valid
			expectedOutput: "(Valid",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := tc.input

			if !utf8.ValidString(tc.input) {
				v := make([]rune, 0, len(tc.input))
				for i, r := range tc.input {
					if r == utf8.RuneError {
						_, size := utf8.DecodeRuneInString(tc.input[i:])
						if size == 1 {
							continue
						}
					}
					v = append(v, r)
				}
				result = string(v)
			}

			if result != tc.expectedOutput {
				t.Errorf("UTF-8 filtering failed: got %q, want %q", result, tc.expectedOutput)
			}
		})
	}
}

func TestRegexpFiltering(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name        string
		pattern     string
		input       string
		shouldMatch bool
	}{
		{
			name:        "Simple regexp match",
			pattern:     "^test.*$",
			input:       "test string",
			shouldMatch: true,
		},
		{
			name:        "Regexp with numbers",
			pattern:     "player[0-9]+",
			input:       "player123 connected",
			shouldMatch: true,
		},
		{
			name:        "Regexp for player events",
			pattern:     "Got character ZDOID from .*",
			input:       "Got character ZDOID from Thorbjorn",
			shouldMatch: true,
		},
		{
			name:        "Regexp with no match",
			pattern:     "server_shutdown",
			input:       "player connecting",
			shouldMatch: false,
		},
		{
			name:        "Complex regexp",
			pattern:     "(Got|Player) (character|connection) [A-Za-z0-9]+",
			input:       "Got character ABC123",
			shouldMatch: true,
		},
		{
			name:        "Complex regexp no match",
			pattern:     "(Got|Player) (character|connection) [A-Za-z0-9]+",
			input:       "System connection ABC123",
			shouldMatch: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			pattern, err := regexp.Compile(tc.pattern)
			if err != nil {
				t.Fatalf("Failed to compile regexp %q: %v", tc.pattern, err)
			}

			matched := pattern.MatchString(tc.input)

			if tc.shouldMatch && !matched {
				t.Errorf("Expected pattern %q to match input %q, but it didn't", tc.pattern, tc.input)
			} else if !tc.shouldMatch && matched {
				t.Errorf("Expected pattern %q NOT to match input %q, but it did", tc.pattern, tc.input)
			}

			action := RegexpAction{
				Filter: pattern,
				Cmd:    "echo test",
			}

			matched = action.Filter.MatchString(tc.input)

			if tc.shouldMatch && !matched {
				t.Errorf("RegexpAction with pattern %q failed to match input %q", tc.pattern, tc.input)
			} else if !tc.shouldMatch && matched {
				t.Errorf("RegexpAction with pattern %q incorrectly matched input %q", tc.pattern, tc.input)
			}
		})
	}
}

func TestFiltering(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tmpDir, err := os.MkdirTemp("", "valheim-logfilter-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	t.Skip("Skipping integration test - requires binary build")
}

func TestLogFiltering(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		line         string
		shouldFilter bool
		filterType   string
		patternOrCmd string
		commandToRun string
	}{
		{
			name:         "Exact match with filter",
			line:         "exact match test",
			shouldFilter: true,
			filterType:   "match",
			patternOrCmd: "exact match test",
		},
		{
			name:         "Prefix match with filter",
			line:         "prefix: some text",
			shouldFilter: true,
			filterType:   "prefix",
			patternOrCmd: "prefix:",
		},
		{
			name:         "Suffix match with filter",
			line:         "some text with suffix",
			shouldFilter: true,
			filterType:   "suffix",
			patternOrCmd: "with suffix",
		},
		{
			name:         "Contains match with filter",
			line:         "line containing substring in the middle",
			shouldFilter: true,
			filterType:   "contains",
			patternOrCmd: "substring",
		},
		{
			name:         "No match should not filter",
			line:         "this line shouldn't match anything",
			shouldFilter: false,
		},
		{
			name:         "Empty line with empty filter",
			line:         "",
			shouldFilter: true,
			filterType:   "empty",
		},
		{
			name:         "Line with command hook",
			line:         "player connected",
			shouldFilter: false,
			filterType:   "contains",
			patternOrCmd: "connected",
			commandToRun: "echo command executed",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var result bool

			switch tc.filterType {
			case "match":
				action := PatternAction{Filter: tc.patternOrCmd, Cmd: tc.commandToRun}
				if tc.line == action.Filter {
					result = removeLogLine(action.Cmd, tc.line)
				}
			case "prefix":
				action := PatternAction{Filter: tc.patternOrCmd, Cmd: tc.commandToRun}
				if strings.HasPrefix(tc.line, action.Filter) {
					result = removeLogLine(action.Cmd, tc.line)
				}
			case "suffix":
				action := PatternAction{Filter: tc.patternOrCmd, Cmd: tc.commandToRun}
				if strings.HasSuffix(tc.line, action.Filter) {
					result = removeLogLine(action.Cmd, tc.line)
				}
			case "contains":
				action := PatternAction{Filter: tc.patternOrCmd, Cmd: tc.commandToRun}
				if strings.Contains(tc.line, action.Filter) {
					result = removeLogLine(action.Cmd, tc.line)
				}
			case "empty":
				if tc.line == "" {
					result = true
				}
			default:
				result = false
			}

			if tc.shouldFilter && !result {
				t.Errorf("Expected line %q to be filtered with %s filter", tc.line, tc.filterType)
			} else if !tc.shouldFilter && result {
				t.Errorf("Expected line %q NOT to be filtered with %s filter", tc.line, tc.filterType)
			}
		})
	}
}
