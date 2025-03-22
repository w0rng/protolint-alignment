package rules

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/rule"
	"github.com/yoheimuta/protolint/linter/visitor"
)

type AlignmentRule struct {
	fixMode bool
}

func NewAlignmentRule(fixMode bool) AlignmentRule {
	return AlignmentRule{
		fixMode: fixMode,
	}
}

func (r AlignmentRule) ID() string {
	return "ALIGNMENT_RULE"
}

func (r AlignmentRule) Purpose() string {
	return "Alignment by sign is equal to (like go)."
}

func (r AlignmentRule) IsOfficial() bool {
	return false
}

func (r AlignmentRule) Severity() rule.Severity {
	return rule.SeverityNote
}

func (r AlignmentRule) Apply(
	proto *parser.Proto,
) ([]report.Failure, error) {
	base, err := visitor.NewBaseFixableVisitor(r.ID(), r.fixMode, proto, string(r.Severity()))
	if err != nil {
		return nil, fmt.Errorf("failed to create base fixable visitor: %w", err)
	}

	v := &alignmentVisitor{
		BaseFixableVisitor: base,
		fixMode:            r.fixMode,
	}
	return visitor.RunVisitor(v, proto, r.ID())
}

type alignFix struct {
	num                 int
	left, right, result string
}

type alignmentVisitor struct {
	*visitor.BaseFixableVisitor
	fixMode bool
}

func (v alignmentVisitor) Finally() error {
	if v.fixMode {
		return v.fix()
	}
	return nil
}

func (v alignmentVisitor) fix() error {
	v.Fixer.ReplaceAll(func(lines []string) []string {
		result := make([]string, len(lines))

		needFix := make([]alignFix, 0)
		for i, line := range lines {
			if strings.HasPrefix(strings.TrimSpace(line), "//") {
				result[i] = line
				continue
			}
			if !strings.Contains(line, "=") {
				for _, fix := range calcLines(needFix) {
					result[fix.num] = fix.result
				}
				needFix = make([]alignFix, 0)
				result[i] = line
				continue
			}

			split := strings.Split(line, "=")
			left, right := split[0], strings.Join(split[1:], "=")
			needFix = append(needFix,
				alignFix{
					i,
					strings.TrimRightFunc(left, unicode.IsSpace),
					strings.TrimLeftFunc(right, unicode.IsSpace),
					"",
				})
		}

		return result
	})

	return v.BaseFixableVisitor.Finally()
}

func calcLines(lines []alignFix) []alignFix {
	maxLeftLen := 0
	for _, line := range lines {
		if len(line.left) > maxLeftLen {
			maxLeftLen = len(line.left)
		}
	}
	for i, line := range lines {
		needLen := maxLeftLen - len(line.left)
		lines[i].result = line.left + strings.Repeat(" ", needLen) + " = " + line.right
	}

	return lines
}
