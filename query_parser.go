package solr

import (
	"fmt"
	"strings"
)

// QueryParser is an abstraction of a query parser
// e.g. standard (lucene), dismax, edismax, boost, block join etc.
type QueryParser interface {
	// BuildParser builds the query from the specified parameters
	BuildParser() string
}

// StandardQueryParser is a standard query parser (lucene)
type StandardQueryParser struct {
	// standard q parser params
	// reference: https://lucene.apache.org/solr/guide/8_7/the-standard-q-parser.html
	q   string // query
	op  string // default operator
	df  string // default field
	sow bool   // split on whitespace
	tag string // tag
}

var _ QueryParser = (*StandardQueryParser)(nil)

// NewStandardQueryParser returns a new StandardQueryParser
func NewStandardQueryParser() *StandardQueryParser {
	return &StandardQueryParser{}
}

// BuildParser builds the query parser
func (qp *StandardQueryParser) BuildParser() string {
	// kv is the list of key-value pair
	// in local-param style query
	kv := []string{"lucene"}
	if qp.df != "" {
		kv = append(kv, fmt.Sprintf("df=%s", qp.df))
	}

	if qp.op != "" {
		kv = append(kv, fmt.Sprintf("q.op=%s", qp.op))
	}

	if qp.sow {
		kv = append(kv, "sow=true")
	}

	if qp.tag != "" {
		kv = append(kv, fmt.Sprintf("tag=%s", qp.tag))
	}

	if qp.q != "" {
		kv = append(kv, fmt.Sprintf("v=%s", qp.q))
	}

	return fmt.Sprintf("{!%s}", strings.Join(kv, " "))
}

// Query sets the query
func (qp *StandardQueryParser) Query(query string) *StandardQueryParser {
	qp.q = query
	return qp
}

// Op sets the default operator
func (qp *StandardQueryParser) Op(op string) *StandardQueryParser {
	qp.op = op
	return qp
}

// Df sets the default field
func (qp *StandardQueryParser) Df(df string) *StandardQueryParser {
	qp.df = df
	return qp
}

// Sow enables split-white-space
func (qp *StandardQueryParser) Sow() *StandardQueryParser {
	qp.sow = true
	return qp
}

// Tag sets the tag param
func (qp *StandardQueryParser) Tag(tag string) *StandardQueryParser {
	qp.tag = tag
	return qp
}

// DisMaxQueryParser is a dismax query parser
type DisMaxQueryParser struct {
	// dismax q parser params
	// reference: https://lucene.apache.org/solr/guide/8_7/the-dismax-q-parser.html
	q   string // query
	alt string // alt query
	qf  string // query fields
	mm  string // minimum should match
	pf  string // phrase field
	ps  string // phrase slop
	qs  string // query slop
	tie string // tie breaker parameter
	bq  string // boost query
	bf  string // boost function
}

var _ QueryParser = (*DisMaxQueryParser)(nil)

// NewDisMaxQueryParser returns a new DisMaxQueryParser
func NewDisMaxQueryParser() *DisMaxQueryParser {
	return &DisMaxQueryParser{}
}

// BuildParser builds the query parser
func (qp *DisMaxQueryParser) BuildParser() string {
	kv := []string{"dismax"}
	if qp.alt != "" {
		kv = append(kv, fmt.Sprintf("q.alt=%s", qp.alt))
	}

	if qp.qf != "" {
		kv = append(kv, fmt.Sprintf("qf=%s", qp.qf))
	}

	if qp.mm != "" {
		kv = append(kv, fmt.Sprintf("mm=%s", qp.mm))
	}

	if qp.pf != "" {
		kv = append(kv, fmt.Sprintf("pf=%s", qp.pf))
	}

	if qp.ps != "" {
		kv = append(kv, fmt.Sprintf("ps=%s", qp.ps))
	}

	if qp.qs != "" {
		kv = append(kv, fmt.Sprintf("qs=%s", qp.qs))
	}

	if qp.tie != "" {
		kv = append(kv, fmt.Sprintf("tie=%s", qp.tie))
	}

	if qp.bq != "" {
		kv = append(kv, fmt.Sprintf("bq=%s", qp.bq))
	}

	if qp.bf != "" {
		kv = append(kv, fmt.Sprintf("bf=%s", qp.bf))
	}

	if qp.q != "" {
		kv = append(kv, fmt.Sprintf("v=%s", qp.q))
	}

	return fmt.Sprintf("{!%s}", strings.Join(kv, " "))
}

// Query sets the query
func (qp *DisMaxQueryParser) Query(query string) *DisMaxQueryParser {
	qp.q = query
	return qp
}

// Alt sets the q.alt param
func (qp *DisMaxQueryParser) Alt(alt string) *DisMaxQueryParser {
	qp.alt = alt
	return qp
}

// WithQf sets the qf param
func (qp *DisMaxQueryParser) WithQf(qf string) *DisMaxQueryParser {
	qp.qf = qf
	return qp
}

// Mm sets the minimum should match param
func (qp *DisMaxQueryParser) Mm(mm string) *DisMaxQueryParser {
	qp.mm = mm
	return qp
}

// Pf sets the phrase field param
func (qp *DisMaxQueryParser) Pf(pf string) *DisMaxQueryParser {
	qp.pf = pf
	return qp
}

// Ps sets the phrase slop param
func (qp *DisMaxQueryParser) Ps(ps string) *DisMaxQueryParser {
	qp.ps = ps
	return qp
}

// Qs sets the query slop param
func (qp *DisMaxQueryParser) Qs(qs string) *DisMaxQueryParser {
	qp.qs = qs
	return qp
}

// Tie sets the tie breaker param param
func (qp *DisMaxQueryParser) Tie(tie string) *DisMaxQueryParser {
	qp.tie = tie
	return qp
}

// Bq sets the boost query param
func (qp *DisMaxQueryParser) Bq(bq string) *DisMaxQueryParser {
	qp.bq = bq
	return qp
}

// Bf sets the boost function param
func (qp *DisMaxQueryParser) Bf(bf string) *DisMaxQueryParser {
	qp.bf = bf
	return qp
}

// ParentQueryParser is a block-join parent query parser
type ParentQueryParser struct {
	which,
	tag,
	filters,
	excludeTags,
	score,
	q string
}

var _ QueryParser = (*ParentQueryParser)(nil)

// NewParentQueryParser returns a new ParentQueryParser
func NewParentQueryParser() *ParentQueryParser {
	return &ParentQueryParser{}
}

// BuildParser builds the query parser
func (qp *ParentQueryParser) BuildParser() string {
	kv := []string{"parent"}

	if qp.which != "" {
		kv = append(kv, fmt.Sprintf("which=%s", qp.which))
	}

	if qp.tag != "" {
		kv = append(kv, fmt.Sprintf("tag=%s", qp.tag))
	}

	if qp.filters != "" {
		kv = append(kv, fmt.Sprintf("filters=%s", qp.filters))
	}

	if qp.excludeTags != "" {
		kv = append(kv, fmt.Sprintf("excludeTags=%s", qp.excludeTags))
	}

	if qp.score != "" {
		kv = append(kv, fmt.Sprintf("score=%s", qp.score))
	}

	if qp.q != "" {
		kv = append(kv, fmt.Sprintf("v=%s", qp.q))
	}

	return fmt.Sprintf("{!%s}", strings.Join(kv, " "))
}

// Which sets the which param
func (qp *ParentQueryParser) Which(which string) *ParentQueryParser {
	qp.which = which
	return qp
}

// Tag sets the tag param
func (qp *ParentQueryParser) Tag(tag string) *ParentQueryParser {
	qp.tag = tag
	return qp
}

// Filters sets the filters param
func (qp *ParentQueryParser) Filters(filters string) *ParentQueryParser {
	qp.filters = filters
	return qp
}

// ExcludeTags sets the excludeTags param
func (qp *ParentQueryParser) ExcludeTags(tags string) *ParentQueryParser {
	qp.excludeTags = tags
	return qp
}

// Score sets the score param
func (qp *ParentQueryParser) Score(score string) *ParentQueryParser {
	qp.score = score
	return qp
}

// Query sets the query
func (qp *ParentQueryParser) Query(query string) *ParentQueryParser {
	qp.q = query
	return qp
}

// ChildrenQueryParser is a block-join children query parser
type ChildrenQueryParser struct {
	of,
	filters,
	excludeTags,
	query string
}

var _ QueryParser = (*ChildrenQueryParser)(nil)

// NewChildrenQueryParser returns a new ChildrenQueryParser
func NewChildrenQueryParser() *ChildrenQueryParser {
	return &ChildrenQueryParser{}
}

// BuildParser builds the query parser
func (qp *ChildrenQueryParser) BuildParser() string {
	kv := []string{"child"}

	if qp.of != "" {
		kv = append(kv, fmt.Sprintf("of=%s", qp.of))
	}

	if qp.filters != "" {
		kv = append(kv, fmt.Sprintf("filters=%s", qp.filters))
	}

	if qp.excludeTags != "" {
		kv = append(kv, fmt.Sprintf("excludeTags=%s", qp.excludeTags))
	}

	if qp.query != "" {
		kv = append(kv, fmt.Sprintf("v=%s", qp.query))
	}

	return fmt.Sprintf("{!%s}", strings.Join(kv, " "))
}

// Query sets the query
func (qp *ChildrenQueryParser) Query(query string) *ChildrenQueryParser {
	qp.query = query
	return qp
}

// Of sets the block-mask 'of' param
func (qp *ChildrenQueryParser) Of(of string) *ChildrenQueryParser {
	qp.of = of
	return qp
}

// Filters sets the filters param
func (qp *ChildrenQueryParser) Filters(filters string) *ChildrenQueryParser {
	qp.filters = filters
	return qp
}

// ExcludeTags sets the excludeTags param
func (qp *ChildrenQueryParser) ExcludeTags(tags string) *ChildrenQueryParser {
	qp.excludeTags = tags
	return qp
}

// FiltersQueryParser is a filters query parser
type FiltersQueryParser struct {
	param,
	excludeTags,
	q string
}

var _ QueryParser = (*FiltersQueryParser)(nil)

// NewFiltersQueryParser returns a new FiltersQueryParser
func NewFiltersQueryParser() *FiltersQueryParser {
	return &FiltersQueryParser{}
}

// BuildParser builds the query parser
func (qp *FiltersQueryParser) BuildParser() string {
	kv := []string{"filters"}

	if qp.param != "" {
		kv = append(kv, fmt.Sprintf("param=%s", qp.param))
	}

	if qp.excludeTags != "" {
		kv = append(kv, fmt.Sprintf("excludeTags=%s", qp.excludeTags))
	}

	if qp.q != "" {
		kv = append(kv, fmt.Sprintf("v=%s", qp.q))
	}

	return fmt.Sprintf("{!%s}", strings.Join(kv, " "))
}

// Param sets the 'param' param
func (qp *FiltersQueryParser) Param(param string) *FiltersQueryParser {
	qp.param = param
	return qp
}

// ExcludeTags sets the excludeTags param
func (qp *FiltersQueryParser) ExcludeTags(tags string) *FiltersQueryParser {
	qp.excludeTags = tags
	return qp
}

// Query sets the query
func (qp *FiltersQueryParser) Query(query string) *FiltersQueryParser {
	qp.q = query
	return qp
}
