package alignmentrule

import (
	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/rule"
	"github.com/yoheimuta/protolint/linter/visitor"
)

type AlignmentRule struct {
	fixMode  bool
	severity rule.Severity
}

func New(fixMode bool, severity rule.Severity) AlignmentRule {
	return AlignmentRule{
		fixMode:  fixMode,
		severity: severity,
	}
}

func (r AlignmentRule) ID() string {
	return "alignment rule"
}

func (r AlignmentRule) Purpose() string {
	return "Aligment rule"
}

func (r AlignmentRule) IsOfficial() bool {
	return false
}

func (r AlignmentRule) Apply(
	proto *parser.Proto,
) ([]report.Failure, error) {
	base, err := visitor.NewBaseFixableVisitor(r.ID(), true, proto, string(r.Severity()))
	if err != nil {
		return nil, err
	}

	v := &alignmentVisitor{
		BaseFixableVisitor: base,
		fixMode:            true,
	}
	return visitor.RunVisitor(v, proto, r.ID())
}

func (r AlignmentRule) Severity() rule.Severity {
	return r.severity
}
