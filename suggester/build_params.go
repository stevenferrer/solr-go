package suggester

import (
	"net/url"
	"strconv"
)

func buildParams(p Params) string {
	params := url.Values{}

	params.Add("suggest", "true")
	params.Add("suggest.q", p.Query)

	for _, dict := range p.Dictionaries {
		params.Add("suggest.dictionary", dict)
	}

	if p.Count > 0 {
		params.Add("suggest.count", strconv.Itoa(p.Count))
	}

	if p.Cfq != "" {
		params.Add("suggest.cfg", p.Cfq)
	}

	if p.Build {
		params.Add("suggest.build", "true")
	}

	if p.Reload {
		params.Add("suggest.reload", "true")
	}

	if p.BuildAll {
		params.Add("suggest.buildAll", "true")
	}

	if p.ReloadAll {
		params.Add("suggest.reloadAll", "true")
	}

	return params.Encode()
}
