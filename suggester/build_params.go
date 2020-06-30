package suggester

import (
	"fmt"
	"net/url"
)

func buildParams(p Params) []string {
	params := []string{}

	enc := func(s string) string {
		return url.QueryEscape(s)
	}

	params = append(params, "suggest=true",
		fmt.Sprintf("suggest.q=%s", enc(p.Query)))

	for _, dict := range p.Dictionaries {
		params = append(params, fmt.Sprintf("suggest.dictionary=%s", enc(dict)))
	}

	if p.Count > 0 {
		params = append(params, fmt.Sprintf("suggest.count=%d", p.Count))
	}

	if p.Cfq != "" {
		params = append(params, fmt.Sprintf("suggest.cfg=%s", enc(p.Cfq)))
	}

	if p.Build {
		params = append(params, "suggest.build=true")
	}

	if p.Reload {
		params = append(params, "suggest.reload=true")
	}

	if p.BuildAll {
		params = append(params, "suggest.buildAll=true")
	}

	if p.ReloadAll {
		params = append(params, "suggest.reloadAll=true")
	}

	return params
}
