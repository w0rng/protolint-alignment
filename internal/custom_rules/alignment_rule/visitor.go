package alignment_rule

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/go-protoparser/v4/parser/meta"
	"github.com/yoheimuta/protolint/linter/visitor"
)

type alignFix struct {
	currentChars int
	replacement  string
	pos          meta.Position
	isLast       bool
}

type alignmentVisitor struct {
	*visitor.BaseFixableVisitor
	fixMode          bool
	alignFixes       map[int][]alignFix
	notInsertNewline bool
}

func (v alignmentVisitor) Finally() error {
	if v.fixMode {
		return v.fix()
	}
	return nil
}

func (v alignmentVisitor) VisitEnum(e *parser.Enum) (next bool) {
	v.validateAlignLeading(e.Meta.Pos)
	defer func() { v.validateAlignLast(e.Meta.LastPos) }()
	for _, comment := range e.Comments {
		v.validateAlignLeading(comment.Meta.Pos)
	}

	for _, body := range e.EnumBody {
		body.Accept(v)
	}
	return false
}

func (v alignmentVisitor) VisitEnumField(f *parser.EnumField) (next bool) {
	v.validateAlignLeading(f.Meta.Pos)
	for _, comment := range f.Comments {
		v.validateAlignLeading(comment.Meta.Pos)
	}
	return false
}

func (v alignmentVisitor) VisitExtend(e *parser.Extend) (next bool) {
	v.validateAlignLeading(e.Meta.Pos)
	defer func() { v.validateAlignLast(e.Meta.LastPos) }()
	for _, comment := range e.Comments {
		v.validateAlignLeading(comment.Meta.Pos)
	}

	for _, body := range e.ExtendBody {
		body.Accept(v)
	}
	return false
}

func (v alignmentVisitor) VisitField(f *parser.Field) (next bool) {
	v.validateAlignLeading(f.Meta.Pos)
	for _, comment := range f.Comments {
		v.validateAlignLeading(comment.Meta.Pos)
	}
	return false
}

func (v alignmentVisitor) VisitGroupField(f *parser.GroupField) (next bool) {
	v.validateAlignLeading(f.Meta.Pos)
	defer func() { v.validateAlignLast(f.Meta.LastPos) }()
	for _, comment := range f.Comments {
		v.validateAlignLeading(comment.Meta.Pos)
	}

	for _, body := range f.MessageBody {
		body.Accept(v)
	}
	return false
}

func (v alignmentVisitor) VisitImport(i *parser.Import) (next bool) {
	v.validateAlignLeading(i.Meta.Pos)
	for _, comment := range i.Comments {
		v.validateAlignLeading(comment.Meta.Pos)
	}
	return false
}

func (v alignmentVisitor) VisitMapField(m *parser.MapField) (next bool) {
	v.validateAlignLeading(m.Meta.Pos)
	for _, comment := range m.Comments {
		v.validateAlignLeading(comment.Meta.Pos)
	}
	return false
}

func (v alignmentVisitor) VisitMessage(m *parser.Message) (next bool) {
	v.validateAlignLeading(m.Meta.Pos)
	defer func() { v.validateAlignLast(m.Meta.LastPos) }()
	for _, comment := range m.Comments {
		v.validateAlignLeading(comment.Meta.Pos)
	}

	for _, body := range m.MessageBody {
		body.Accept(v)
	}
	return false
}

func (v alignmentVisitor) VisitOneof(o *parser.Oneof) (next bool) {
	v.validateAlignLeading(o.Meta.Pos)
	defer func() { v.validateAlignLast(o.Meta.LastPos) }()
	for _, comment := range o.Comments {
		v.validateAlignLeading(comment.Meta.Pos)
	}

	for _, field := range o.OneofFields {
		field.Accept(v)
	}
	return false
}

func (v alignmentVisitor) VisitOneofField(f *parser.OneofField) (next bool) {
	v.validateAlignLeading(f.Meta.Pos)
	for _, comment := range f.Comments {
		v.validateAlignLeading(comment.Meta.Pos)
	}
	return false
}

func (v alignmentVisitor) VisitOption(o *parser.Option) (next bool) {
	v.validateAlignLeading(o.Meta.Pos)
	for _, comment := range o.Comments {
		v.validateAlignLeading(comment.Meta.Pos)
	}
	return false
}

func (v alignmentVisitor) VisitPackage(p *parser.Package) (next bool) {
	v.validateAlignLeading(p.Meta.Pos)
	for _, comment := range p.Comments {
		v.validateAlignLeading(comment.Meta.Pos)
	}
	return false
}

func (v alignmentVisitor) VisitReserved(r *parser.Reserved) (next bool) {
	v.validateAlignLeading(r.Meta.Pos)
	for _, comment := range r.Comments {
		v.validateAlignLeading(comment.Meta.Pos)
	}
	return false
}

func (v alignmentVisitor) VisitRPC(r *parser.RPC) (next bool) {
	v.validateAlignLeading(r.Meta.Pos)
	defer func() {
		line := v.Fixer.Lines()[r.Meta.LastPos.Line-1]
		runes := []rune(line)
		for i := r.Meta.LastPos.Column - 2; 0 < i; i-- {
			r := runes[i]
			if r == '{' || r == ')' {
				// skip validating the alignment when the line ends with {}, {};, or );
				return
			}
			if r == '}' || unicode.IsSpace(r) {
				continue
			}
			break
		}
		v.validateAlignLast(r.Meta.LastPos)
	}()
	for _, comment := range r.Comments {
		v.validateAlignLeading(comment.Meta.Pos)
	}

	for _, body := range r.Options {
		body.Accept(v)
	}
	return false
}

func (v alignmentVisitor) VisitService(s *parser.Service) (next bool) {
	v.validateAlignLeading(s.Meta.Pos)
	defer func() { v.validateAlignLast(s.Meta.LastPos) }()
	for _, comment := range s.Comments {
		v.validateAlignLeading(comment.Meta.Pos)
	}

	for _, body := range s.ServiceBody {
		body.Accept(v)
	}
	return false
}

func (v alignmentVisitor) VisitSyntax(s *parser.Syntax) (next bool) {
	v.validateAlignLeading(s.Meta.Pos)
	for _, comment := range s.Comments {
		v.validateAlignLeading(comment.Meta.Pos)
	}
	return false
}

func (v alignmentVisitor) VisitEdition(s *parser.Edition) (next bool) {
	v.validateAlignLeading(s.Meta.Pos)
	for _, comment := range s.Comments {
		v.validateAlignLeading(comment.Meta.Pos)
	}
	return false
}

func (v alignmentVisitor) validateAlignLeading(
	pos meta.Position,
) {
	v.validateAlign(pos, false)
}

func (v alignmentVisitor) validateAlignLast(
	pos meta.Position,
) {
	v.validateAlign(pos, true)
}

func (v alignmentVisitor) validateAlign(
	pos meta.Position,
	isLast bool,
) {
	line := v.Fixer.Lines()[pos.Line-1]
	equalsIndex := strings.Index(line, "=")
	if equalsIndex == -1 {
		return
	}

	leading := line[:equalsIndex]
	alignment := strings.Repeat(" ", equalsIndex)

	v.alignFixes[pos.Line-1] = append(v.alignFixes[pos.Line-1], alignFix{
		currentChars: len(leading),
		replacement:  alignment,
		pos:          pos,
		isLast:       isLast,
	})

	if leading == alignment {
		return
	}
	if 1 < len(v.alignFixes[pos.Line-1]) && v.notInsertNewline {
		return
	}
	if len(v.alignFixes[pos.Line-1]) == 1 {
		v.AddFailuref(
			pos,
			`Found an incorrect alignment style "%s". "%s" is correct.`,
			leading,
			alignment,
		)
	} else {
		v.AddFailuref(
			pos,
			`Found a possible incorrect alignment style. Inserting a new line is recommended.`,
		)
	}
}

func (v alignmentVisitor) fix() error {
	var shouldFixed bool

	v.Fixer.ReplaceAll(func(lines []string) []string {
		var fixedLines []string
		for i, line := range lines {
			fmt.Printf("line: %+v\n", line)
			lines := []string{line}
			if fixes, ok := v.alignFixes[i]; ok {
				lines[0] = fixes[0].replacement + line[fixes[0].currentChars:]
				shouldFixed = true

				if 1 < len(fixes) && !v.notInsertNewline {
					// compose multiple lines in reverse order from right to left on one line.
					var rlines []string
					for j := len(fixes) - 1; 0 <= j; j-- {
						alignment := strings.Repeat(" ", fixes[j].pos.Column-1)
						if fixes[j].isLast {
							// deal with last position followed by ';'. See https://github.com/yoheimuta/protolint/issues/99
							for line[fixes[j].pos.Column-1] == ';' {
								fixes[j].pos.Column--
							}
						}

						endColumn := len(line)
						if j < len(fixes)-1 {
							endColumn = fixes[j+1].pos.Column - 1
						}
						text := line[fixes[j].pos.Column-1 : endColumn]
						text = strings.TrimRightFunc(text, func(r rune) bool {
							// removing right spaces is a possible side effect that users do not expect,
							// but it's probably acceptable and usually recommended.
							return unicode.IsSpace(r)
						})

						rlines = append(rlines, alignment+text)
					}

					// sort the multiple lines in order
					lines = []string{}
					for j := len(rlines) - 1; 0 <= j; j-- {
						lines = append(lines, rlines[j])
					}
				}
			}
			fixedLines = append(fixedLines, lines...)
		}
		return fixedLines
	})

	if !shouldFixed {
		return nil
	}
	return v.BaseFixableVisitor.Finally()
}
