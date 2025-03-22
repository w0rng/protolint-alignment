package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/w0rng/protolint-alignment/internal/rules"
	"github.com/w0rng/protolint-alignment/internal/utils"
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
	correctMessagePath := must(newTestIndentData("message.proto"))
	incorrectMessagePath := must(newTestIndentData("incorrect_message.proto"))
	correctIssue99Path := must(newTestIndentData("issue_99.proto"))
	incorrectIssue99Path := must(newTestIndentData("incorrect_issue_99.proto"))

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
		{
			name:            "correct message",
			inputTestData:   correctMessagePath,
			wantCorrectData: correctMessagePath,
		},
		{
			name:            "incorrect message",
			inputTestData:   incorrectMessagePath,
			wantCorrectData: correctMessagePath,
		},
		{
			name:            "correct issue_99",
			inputTestData:   correctIssue99Path,
			wantCorrectData: correctIssue99Path,
		},
		{
			name:            "incorrect issue_99",
			inputTestData:   incorrectIssue99Path,
			wantCorrectData: correctIssue99Path,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ruleToTest := rules.NewAlignmentRule(true)

			proto, err := utils.NewProtoFile(tt.inputTestData.FilePath, tt.inputTestData.FilePath).Parse(false)
			require.NoError(t, err)

			_, err = ruleToTest.Apply(proto)
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
			ruleOnlyCheck := rules.NewAlignmentRule(false)
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
