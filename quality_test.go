package budoux

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"
)

const qualitySeparator = "▁"

func TestDefaultJapaneseQuality(t *testing.T) {
	t.Helper()

	file, err := os.Open("testdata/quality_ja.tsv")
	if err != nil {
		t.Fatalf("failed to open quality fixture: %v", err)
	}
	defer file.Close()

	parser := NewDefaultJapaneseParser()
	var failures []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		cols := strings.SplitN(line, "\t", 2)
		if len(cols) != 2 {
			continue
		}

		expected := strings.TrimSpace(cols[1])
		sentence := strings.ReplaceAll(expected, qualitySeparator, "")
		actual := strings.Join(parser.Parse(sentence), qualitySeparator)
		if actual != expected {
			failures = append(failures, fmt.Sprintf("expected:%s\tactual:%s", expected, actual))
		}
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("failed to read quality fixture: %v", err)
	}
	if len(failures) > 0 {
		t.Fatalf("quality regression detected:\n%s", strings.Join(failures, "\n"))
	}
}
