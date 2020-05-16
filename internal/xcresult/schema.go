package xcresult

import (
	"encoding/json"

	"github.com/tidwall/gjson"
)

type TestFailureSummary struct {
	Type            Type      `json:"_type"`
	Doc             Document  `json:"documentLocationInCreatingWorkspace"`
	IssueType       IssueType `json:"issueType"`
	Message         IssueType `json:"message"`
	ProducingTarget IssueType `json:"producingTarget"`
	TestCaseName    IssueType `json:"testCaseName"`
}

type Document struct {
	Type             SupertypeClass `json:"_type"`
	ConcreteTypeName IssueType      `json:"concreteTypeName"`
	URL              IssueType      `json:"url"`
}

type IssueType struct {
	Type  SupertypeClass `json:"_type"`
	Value string         `json:"_value"`
}

type SupertypeClass struct {
	Name string `json:"_name"`
}

type Type struct {
	Name      string         `json:"_name"`
	Supertype SupertypeClass `json:"_supertype"`
}

func ParseIssues(data []byte) ([]TestFailureSummary, error) {
	// Retrieve values
	res := gjson.GetBytes(data, "issues.testFailureSummaries._values")

	var ts []TestFailureSummary
	if err := json.Unmarshal([]byte(res.Raw), &ts); err != nil {
		return nil, err
	}

	return ts, nil
}
