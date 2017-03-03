package nifi

import (
	"testing"
	"reflect"
	"io/ioutil"
)

const (
	testGetProcessGroupFlowFile = "../resources/example-ProcessGroupFlow.json"
)

func TestProcessGetProcessGroupFlow(t *testing.T) {
	expectedResult := []string{"connections","funnels","inputPorts","labels",
		"outputPorts","processGroups","processors","remoteProcessGroups"}

	rawJson, _ := ioutil.ReadFile(testGetProcessGroupFlowFile)
	actualResult1 := ProcessGetProcessGroupFlow(rawJson)
	actualResult1 = actualResult1
	actualResult := []string{"test","123"}
	//if reflect.TypeOf(actualResult) != reflect.TypeOf(expectedResult) {
	//	t.Fatalf("Expected true but got false")
	//}
	if !reflect.DeepEqual(actualResult, expectedResult) {
		t.Fatalf("Expected %v but got %v",expectedResult,actualResult)
	}
}
