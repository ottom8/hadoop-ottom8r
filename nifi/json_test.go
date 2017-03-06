package nifi

import (
	"testing"
	"io/ioutil"

	"github.com/onsi/gomega"
)

const (
	testGetProcessGroupFlowFile = "../resources/example-ProcessGroupFlow.json"
	expectedGetProcessGroupFlowFile = "../resources/expected-ProcessGroupFlow.json"
	testProcessSnippetResponse = "../resources/test-ProcessSnippetResponse.json"
)

func TestProcessGetProcessGroupFlow(t *testing.T) {
	expectedResult, _ := ioutil.ReadFile(expectedGetProcessGroupFlowFile)

	rawJson, _ := ioutil.ReadFile(testGetProcessGroupFlowFile)
	actualResult := ProcessGetProcessGroupFlow(rawJson)
	gomega.Expect(actualResult).To(gomega.MatchJSON(expectedResult))
	//if !reflect.DeepEqual(actualResult, expectedResult) {
	//	t.Fatalf("Expected %v but got %v",expectedResult,actualResult)
	//}
}

func TestProcessSnippetResponse(t *testing.T) {
	expectedResult := "test123-id"
	rawJson, _ := ioutil.ReadFile(testProcessSnippetResponse)
	actualResult := ProcessSnippetResponse(rawJson)
	gomega.Expect(actualResult).To(gomega.Equal(expectedResult))
}