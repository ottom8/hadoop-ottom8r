package nifi

import (
	"github.com/Jeffail/gabs"
	"github.com/ottom8/hadoop-ottom8r/logger"
	"encoding/json"
	"fmt"
)

type Snippet struct {
	Id				string `json:"id"`
	URI				string `json:"uri"`
	ParentGroupId	string `json:"parentGroupId"`
	SnippetGroups 	[]SnippetGroup `json:"snippetGroups,omit"`
}

type SnippetGroup struct {
	Name			string `json:"name"`
	SnippetItems 	[]SnippetItem `json:"snippetItems,omit"`
}

type SnippetItem struct {
	Name 			string `json:"name"`
	ClientId 		string `json:"clientId"`
	Version 		float64 `json:"version"`
	LastModifier 	string `json:"lastModifier"`
}

func (s *Snippet) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Id            string `json:"id"`
		URI           string `json:"uri"`
		ParentGroupId string `json:"parentGroupId"`
		SnippetGroups []SnippetGroup `json:""`
	}{
		Id: s.Id,
		URI: s.URI,
		ParentGroupId: s.ParentGroupId,
		SnippetGroups: s.SnippetGroups,
	})
}

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

func getSnippetItem(item *gabs.Container) *SnippetItem {
	snippetItem := SnippetItem{LastModifier:"", ClientId:""}
	snippetItem.Name = item.Path("id").Data().(string)
	snippetItem.Version = item.Path("revision.version").Data().(float64)
	if snippetItem.Version > 0 {
		snippetItem.ClientId = item.Path("revision.clientId").Data().(string)
	}
	return &snippetItem
}

func getSnippetGroup(groupName string, jsonParsed *gabs.Container) *SnippetGroup {
	snippetGroup := SnippetGroup{Name:groupName}
	itemMap, _ := jsonParsed.
		Path(fmt.Sprintf("processGroupFlow.flow.%s", groupName)).Children()
	if len(itemMap) > 0 {
		for _, item := range itemMap {
			snippetItem := getSnippetItem(item)
			snippetGroup.SnippetItems = append(snippetGroup.SnippetItems, *snippetItem)
		}
	}
	return &snippetGroup
}

func getSnippet(jsonParsed *gabs.Container) *Snippet {
	snippet := Snippet{Id:"", URI:""}
	snippet.ParentGroupId = jsonParsed.Path("processGroupFlow.id").Data().(string)
	groupMap, _ := jsonParsed.Path("processGroupFlow.flow").ChildrenMap()
	groups := getMapKeys(groupMap)
	logger.Debug(fmt.Sprintf("%+v", groups))
	for _, group := range groups {
		snippetGroup := getSnippetGroup(group, jsonParsed)
		snippet.SnippetGroups = append(snippet.SnippetGroups, *snippetGroup)
	}
	return &snippet
}

func getSnippetJson(snippet *Snippet) []byte {
	jsonRaw, err := json.Marshal(snippet)
	if err != nil {
		logger.Fatal(err.Error())
	}
	return jsonRaw
}

func ProcessGetProcessGroupFlow(respBody []byte) []byte {
	jsonParsed, _ := gabs.ParseJSON(respBody)
	snippet := getSnippet(jsonParsed)
	//logger.Debug(fmt.Sprintf("%+v", snippet))
	snippetJson := getSnippetJson(snippet)
	logger.Debug(fmt.Sprintf("%s",snippetJson))
	return snippetJson
}

func ProcessGetProcessGroupFlow1(respBody []byte) {
	jsonParsed, _ := gabs.ParseJSON(respBody)
	groupIds, _ := jsonParsed.
		Path("processGroupFlow.flow.remoteProcessGroups.id").Children()
	//for key, child := range children {
	//	logger.Debug(fmt.Sprintf("key: %v, value: %v", key, child.Data().(string)))
	//}
	logger.Debug(fmt.Sprintf("%+v", groupIds))
	revisions, _ := jsonParsed.
		Path("processGroupFlow.flow.remoteProcessGroups.revision").Children()
	logger.Debug(fmt.Sprintf("%+v", revisions))
}
