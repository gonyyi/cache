package cache

import "fmt"

// =====================================================================================================================
// CACHE ITEM
// =====================================================================================================================
type cacheItem struct {
	Request struct {
		Method      string            `json:"method"`
		URL         string            `json:"url"`
		ID          string            `json:"id"`
		Passwd      string            `json:"passwd"`
		Body        string            `json:"body"`
		ContentType string            `json:"content_type"`
		SubVar      map[string]string `json:"sub_var"`
	} `json:"request"`
	Response struct {
		LastUpdated string `json:"last_updated"` // 2020-10-10 24:10:10
		Size        int    `json:"size"`
	} `json:"response"`
}

func (ci *cacheItem) AddSubVarItem(key, value string) {
	ci.Request.SubVar[key] = value
}
func (ci *cacheItem) RemoveSubVarItem(key, value string) {
	delete(ci.Request.SubVar, key)
}
func (ci *cacheItem) GetSubVarItem(key string) (string, error) {
	if v, ok := ci.Request.SubVar[key]; ok {
		return v, nil
	} else {
		return "", fmt.Errorf("subVar <%s> not exist", key)
	}

}
func (ci *cacheItem) ClearSubVar() {
	ci.Request.SubVar = make(map[string]string)
}
func newCacheItem(method, url, id, passwd string) *cacheItem {
	r := cacheItem{}
	r.Request.Method = method
	r.Request.URL = url
	r.Request.ID = id
	r.Request.Passwd = passwd
	r.Request.SubVar = make(map[string]string)
	r.Response.LastUpdated = "0000-00-00 00:00:00"
	return &r
}
