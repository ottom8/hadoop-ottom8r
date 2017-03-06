package nifi

import (
	"github.com/Jeffail/gabs"
	"github.com/ottom8/hadoop-ottom8r/logger"
	"encoding/json"
	"fmt"
)

func processGetProcessGroup(respBody []byte) {
	jsonParsed, _ := gabs.ParseJSON(respBody)
	logger.Debug(fmt.Sprintf("%+v", jsonParsed))
}

func getMapKeys(inMap map[string]*gabs.Container) []string {
	keys := make([]string, len(inMap))
	i := 0
	for k := range inMap {
		keys[i] = k
		i++
	}
	return keys
}

func getSnippetItem(item *gabs.Container) (map[string]interface{}, string) {
	itemId := item.Path("id").Data().(string)
	itemVersion := item.Path("revision.version").Data().(float64)
	snippetItem := map[string]interface{} {"version": itemVersion}
	if itemVersion > 0 {
		snippetItem["clientId"] = item.Path("revision.clientId").Data().(string)
	}
	//snippetItem := map[string]interface{} {itemId: innerItem}
	return snippetItem, itemId
}

func getSnippetGroup(groupName string, jsonParsed *gabs.Container) map[string]interface{} {
	itemMap, _ := jsonParsed.Path("processGroupFlow.flow." + groupName).Children()
	snippetGroup := map[string]interface{} {}
	if len(itemMap) > 0 {
		for _, item := range itemMap {
			snippetItem, itemId := getSnippetItem(item)
			snippetGroup[itemId] = snippetItem
		}
	}
	//snippetGroup := map[string]interface{} {groupName: innerGroup}
	return snippetGroup
}

func getSnippet(jsonParsed *gabs.Container) map[string]interface{} {
	parentGroupId := jsonParsed.Path("processGroupFlow.id").Data().(string)
	subSnippet := map[string]interface{} {"parentGroupId": parentGroupId}
	groupMap, _ := jsonParsed.Path("processGroupFlow.flow").ChildrenMap()
	groups := getMapKeys(groupMap)
	//logger.Debug(fmt.Sprintf("%+v", groups))
	for _, group := range groups {
		snippetGroup := getSnippetGroup(group, jsonParsed)
		subSnippet[group] = snippetGroup
	}
	snippet := map[string]interface{} {"snippet": subSnippet}
	return snippet
}

func getJson(structured interface{}) []byte {
	jsonRaw, err := json.Marshal(structured)
	if err != nil {
		logger.Fatal(err.Error())
	}
	return jsonRaw
}

func getSnippetId(jsonParsed *gabs.Container) string {
	return jsonParsed.Path("snippet.id").Data().(string)
}

// ProcessGetProcessGroupFlow takes the response body from flow/ProcessGroup
// and returns a JSON payload to create a snippet
func ProcessGetProcessGroupFlow(respBody []byte) []byte {
	jsonParsed, _ := gabs.ParseJSON(respBody)
	snippet := getSnippet(jsonParsed)
	snippetJson := getJson(snippet)
	//logger.Debug(fmt.Sprintf("%s",snippetJson))
	return snippetJson
}

// ProcessSnippetResponse takes the response body from snippets post
// and returns snippetId.
func ProcessSnippetResponse(respBody []byte) string {
	jsonParsed, _ := gabs.ParseJSON(respBody)
	snippetId := getSnippetId(jsonParsed)
	return snippetId
}

// ProcessTemplateRequest builds a request body for a templates post
func ProcessTemplateRequest(tmplBody map[string]interface{}) []byte {
	return getJson(tmplBody)
}