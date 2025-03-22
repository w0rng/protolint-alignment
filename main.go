package main

import (
	"github.com/w0rng/protolint-alignment/internal/rules"
	"github.com/yoheimuta/protolint/linter/rule"
	"github.com/yoheimuta/protolint/plugin"
)

func main() {
	plugin.RegisterCustomRules(
		plugin.RuleGen(func(
			_ bool,
			fixMode bool,
		) rule.Rule {
			return rules.NewAlignmentRule(fixMode)
		}),
	)
}
