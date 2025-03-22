package alignment_rule_test

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/w0rng/protolint-alignment/internal/custom_rules/alignment_rule"
	"github.com/w0rng/protolint-alignment/internal/utils"
	"github.com/yoheimuta/protolint/linter/rule"
)

func newTestIndentData(
	fileName string,
) (utils.TestData, error) {
	return utils.NewTestData(utils.TestDataPath(fileName))
}

func TestIndentRule_Apply_fix(t *testing.T) {
	correctSyntaxPath := must(newTestIndentData("syntax.proto"))

	incorrectSyntaxPath := must(newTestIndentData("incorrect_syntax.proto"))

	correctEnumPath := must(newTestIndentData("enum.proto"))

	incorrectEnumPath := must(newTestIndentData("incorrect_enum.proto"))
	//
	//correctMessagePath := must(newTestIndentData("message.proto"))
	//
	//incorrectMessagePath := must(newTestIndentData("incorrect_message.proto"))
	//
	//correctIssue99Path := must(newTestIndentData("issue_99.proto"))
	//
	//incorrectIssue99Path := must(newTestIndentData("incorrect_issue_99.proto"))
	//
	//incorrectIssue139Path := must(newTestIndentData("incorrect_issue_139.proto"))
	//
	//correctIssue139Path := must(newTestIndentData("issue_139.proto"))
	//
	//correctIssue139InsertPath := must(newTestIndentData("issue_139_insert_linebreaks.proto"))

	tests := []struct {
		name               string
		inputTestData      utils.TestData
		inputInsertNewline bool
		wantCorrectData    utils.TestData
	}{
		{
			name:            "correct syntax",
			inputTestData:   correctSyntaxPath,
			wantCorrectData: correctSyntaxPath,
		},
		{
			name:            "incorrect syntax",
			inputTestData:   incorrectSyntaxPath,
			wantCorrectData: correctSyntaxPath,
		},
		{
			name:            "correct enum",
			inputTestData:   correctEnumPath,
			wantCorrectData: correctEnumPath,
		},
		{
			name:            "incorrect enum",
			inputTestData:   incorrectEnumPath,
			wantCorrectData: correctEnumPath,
		},
		//{
		//	name:            "correct message",
		//	inputTestData:   correctMessagePath,
		//	wantCorrectData: correctMessagePath,
		//},
		//{
		//	name:            "incorrect message",
		//	inputTestData:   incorrectMessagePath,
		//	wantCorrectData: correctMessagePath,
		//},
		//{
		//	name:            "correct issue_99",
		//	inputTestData:   correctIssue99Path,
		//	wantCorrectData: correctIssue99Path,
		//},
		//{
		//	name:            "incorrect issue_99",
		//	inputTestData:   incorrectIssue99Path,
		//	wantCorrectData: correctIssue99Path,
		//},
		//{
		//	name:            "do nothing against inner elements on the same line. Fix https://github.com/yoheimuta/protolint/issues/139",
		//	inputTestData:   incorrectIssue139Path,
		//	wantCorrectData: correctIssue139Path,
		//},
		//{
		//	name:               "insert linebreaks against inner elements on the same line. Fix https://github.com/yoheimuta/protolint/issues/139",
		//	inputTestData:      incorrectIssue139Path,
		//	inputInsertNewline: true,
		//	wantCorrectData:    correctIssue139InsertPath,
		//},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule_to_test := alignment_rule.New(
				true,
				rule.SeverityError,
			)

			proto, err := utils.NewProtoFile(tt.inputTestData.FilePath, tt.inputTestData.FilePath).Parse(false)
			require.NoError(t, err)

			_, err = rule_to_test.Apply(proto)
			require.NoError(t, err)

			got, err := tt.inputTestData.Data()
			assert.Equal(t, string(tt.wantCorrectData.OriginData), string(got))

			// restore the file
			defer func() {
				err = tt.inputTestData.Restore()
				if err != nil {
					t.Errorf("got err %v", err)
				}
			}()

			// check whether the modified content can pass the lint in the end.
			ruleOnlyCheck := alignment_rule.New(
				false,
				rule.SeverityError,
			)
			proto, err = utils.NewProtoFile(tt.inputTestData.FilePath, tt.inputTestData.FilePath).Parse(false)
			require.NoError(t, err)
			gotCheck, err := ruleOnlyCheck.Apply(proto)
			require.NoError(t, err)

			if 0 < len(gotCheck) {
				t.Errorf("got failures %v, but want no failures", gotCheck)
				return
			}
		})
	}
}

func must[T any](res T, err error) T {
	if err != nil {
		panic(err)
	}
	return res
}
