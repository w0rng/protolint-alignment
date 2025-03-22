package main

import (
	"github.com/w0rng/protolint-alignment/internal/custom_rules/alignment_rule"
	"github.com/yoheimuta/protolint/linter/rule"
	"github.com/yoheimuta/protolint/plugin"
)

func main() {
	plugin.RegisterCustomRules(
		plugin.RuleGen(func(
			verbose bool,
			fixMode bool,
		) rule.Rule {
			return alignment_rule.New(fixMode, rule.SeverityError)
		}),
	)
}
