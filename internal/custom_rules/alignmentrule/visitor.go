package alignmentrule

import (
	"strings"
	"unicode"

	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/protolint/linter/visitor"
)

type alignFix struct {
	left, right string
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

func (v alignmentVisitor) VisitEnum(_ *parser.Enum) bool             { return false }
func (v alignmentVisitor) VisitEnumField(_ *parser.EnumField) bool   { return false }
func (v alignmentVisitor) VisitExtend(_ *parser.Extend) bool         { return false }
func (v alignmentVisitor) VisitField(_ *parser.Field) bool           { return false }
func (v alignmentVisitor) VisitGroupField(_ *parser.GroupField) bool { return false }
func (v alignmentVisitor) VisitImport(_ *parser.Import) bool         { return false }
func (v alignmentVisitor) VisitMapField(_ *parser.MapField) bool     { return false }
func (v alignmentVisitor) VisitMessage(_ *parser.Message) bool       { return false }
func (v alignmentVisitor) VisitOneof(_ *parser.Oneof) bool           { return false }
func (v alignmentVisitor) VisitOneofField(_ *parser.OneofField) bool { return false }
func (v alignmentVisitor) VisitOption(_ *parser.Option) bool         { return false }
func (v alignmentVisitor) VisitPackage(_ *parser.Package) bool       { return false }
func (v alignmentVisitor) VisitReserved(_ *parser.Reserved) bool     { return false }
func (v alignmentVisitor) VisitRPC(_ *parser.RPC) bool               { return false }
func (v alignmentVisitor) VisitService(_ *parser.Service) bool       { return false }
func (v alignmentVisitor) VisitSyntax(_ *parser.Syntax) bool         { return false }
func (v alignmentVisitor) VisitEdition(_ *parser.Edition) bool       { return false }

func (v alignmentVisitor) fix() error {
	v.Fixer.ReplaceAll(func(lines []string) []string {
		result := make([]string, 0, len(lines))

		needFix := make([]alignFix, 0)
		for _, line := range lines {
			if !strings.Contains(line, "=") {
				result = append(result, calcLines(needFix)...)
				needFix = make([]alignFix, 0)
				result = append(result, line)
				continue
			}

			split := strings.Split(line, "=")
			left, right := split[0], strings.Join(split[1:], "=")
			needFix = append(needFix, alignFix{trimRightSpace(left), trimLeftSpace(right)})
		}

		return result
	})

	return v.BaseFixableVisitor.Finally()
}

func calcLines(lines []alignFix) []string {
	result := make([]string, 0, len(lines))

	maxLeftLen := 0
	for _, line := range lines {
		if len(line.left) > maxLeftLen {
			maxLeftLen = len(line.left)
		}
	}
	for _, line := range lines {
		needLen := maxLeftLen - len(line.left)
		result = append(result, line.left+strings.Repeat(" ", needLen)+" = "+line.right)
	}

	return result
}

func trimRightSpace(s string) string {
	count := 0
	for i := len(s) - 1; i >= 0; i-- {
		if !unicode.IsSpace(rune(s[i])) {
			break
		}
		count++
	}
	return s[:len(s)-count]
}

func trimLeftSpace(s string) string {
	count := 0
	for i := range len(s) {
		if !unicode.IsSpace(rune(s[i])) {
			break
		}
		count++
	}
	return s[count:]
}
