package template

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

const (
	netBase = "Net"
)

var netFuncs = dictionary{
	"httpGet": httpGet,
	"httpDoc": httpDocument,
}

var netFuncsArgs = arguments{
	"httpGet": {"url"},
	"httpDoc": {"url"},
}

var netFuncsAliases = aliases{
	"httpDoc": {"httpDocument", "curl"},
}

var netFuncsHelp = descriptions{
	"httpGet": "Returns http get response from supplied URL.",
	"httpDoc": "Returns http document returned by supplied URL.",
}

func (t *Template) addNetFuncs() {
	t.AddFunctions(netFuncs, netBase, FuncOptions{
		FuncHelp:    netFuncsHelp,
		FuncArgs:    netFuncsArgs,
		FuncAliases: netFuncsAliases,
	})
}

func httpGet(url interface{}) (*http.Response, error) {
	return http.Get(fmt.Sprint(url))
}

func httpDocument(url interface{}) (interface{}, error) {
	response, err := httpGet(url)
	if err != nil {
		return response, err
	}
	return goquery.NewDocumentFromReader(response.Body)
}
