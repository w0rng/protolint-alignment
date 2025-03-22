package main

import (
	"github.com/w0rng/protolint-alignment/internal/custom_rules/alignmentrule"
	"github.com/yoheimuta/protolint/linter/rule"
	"github.com/yoheimuta/protolint/plugin"
)

func main() {
	plugin.RegisterCustomRules(
		plugin.RuleGen(func(
			_ bool,
			fixMode bool,
		) rule.Rule {
			return alignmentrule.New(fixMode, rule.SeverityError)
		}),
	)
}
