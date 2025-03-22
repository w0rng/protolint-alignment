package alignment_rule

import (
	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/rule"
	"github.com/yoheimuta/protolint/linter/visitor"
)

type alignmentRule struct {
	fixMode  bool
	severity rule.Severity
}

func New(fixMode bool, severity rule.Severity) alignmentRule {
	return alignmentRule{
		fixMode:  fixMode,
		severity: severity,
	}
}

func (r alignmentRule) ID() string {
	return "alignment rule"
}

func (r alignmentRule) Purpose() string {
	return "Aligment rule"
}

func (r alignmentRule) IsOfficial() bool {
	return false
}

func (r alignmentRule) Apply(
	proto *parser.Proto,
) ([]report.Failure, error) {
	base, err := visitor.NewBaseFixableVisitor(r.ID(), true, proto, string(r.Severity()))
	if err != nil {
		return nil, err
	}

	v := &alignmentVisitor{
		BaseFixableVisitor: base,
		fixMode:            true,
		alignFixes:         make(map[int][]alignFix),
	}
	return visitor.RunVisitor(v, proto, r.ID())
}

func (r alignmentRule) Severity() rule.Severity {
	return r.severity
}
