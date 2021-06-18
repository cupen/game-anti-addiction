package auth

import (
	"net/http"
	"sort"
	"strings"
)

// Headers ...
type Headers struct {
	AppID     string `json:"appId"`
	BizID     string `json:"bizId"`
	Timestamp string `json:"timestamps"`
	Sign      string `json:"sign"`
}

// SetHTTPHeaders ...
func (h *Headers) SetHTTPHeaders(header http.Header) {
	// NOTE: Do not use header.Set, it will rename the key to align CanonicalMIMEHeader
	// header.Set("appId", h.AppID)
	// header.Set("bizId", h.BizID)
	// header.Set("timestamps", h.Timestamp)
	// header.Set("sign", h.Sign)

	header["appId"] = []string{h.AppID}
	header["bizId"] = []string{h.BizID}
	header["timestamps"] = []string{h.Timestamp}
	header["sign"] = []string{h.Sign}
}

func (h *Headers) Dump() string {
	if h == nil {
		return "null"
	}
	lines := []string{
		"appId=", h.AppID, "\n",
		"bizId=", h.BizID, "\n",
		"timestamps=", h.Timestamp, "\n",
		"sign=", h.Sign, "\n",
	}
	return strings.Join(lines, "")
}

func generateSortedKeys(m map[string]string) []string {
	keys := make([]string, len(m))
	i := 0
	for k, _ := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}
