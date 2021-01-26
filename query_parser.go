package solr

import (
	"errors"
	"fmt"
	"strings"
)

// QueryParser is an abstraction of a query parser
// e.g. standard (lucene), dismax, edismax, boost, block join etc.
type QueryParser interface {
	// BuildParser builds the query from the specified parameters
	BuildParser() (string, error)
}

// StandardQueryParser is a standard query parser a.k.a. lucene
type StandardQueryParser struct {
	// standard q parser params
	// reference: https://lucene.apache.org/solr/guide/8_7/the-standard-q-parser.html
	q   string // query
	op  string // default operator
	df  string // default field
	sow bool   // split on whitespace
}

var _ QueryParser = (*StandardQueryParser)(nil)

// NewStandardQueryParser returns a new StdQueryParser
func NewStandardQueryParser(q string) *StandardQueryParser {
	return &StandardQueryParser{q: q}
}

// BuildParser builds the query parser
func (qp *StandardQueryParser) BuildParser() (string, error) {
	if qp.q == "" {
		return "", errors.New("'q' parameter is required")
	}

	// kv is the list of key-value pair
	// in local-param style query
	kv := []string{"lucene"}
	if qp.df != "" {
		kv = append(kv, "df='"+qp.df+"'")
	}

	if qp.op != "" {
		kv = append(kv, "q.op='"+qp.op+"'")
	}

	if qp.sow {
		kv = append(kv, "sow=true")
	}

	return fmt.Sprintf("{!%s}%s", strings.Join(kv, " "), qp.q), nil
}

// WithQ sets the query
func (qp *StandardQueryParser) WithQ(q string) *StandardQueryParser {
	qp.q = q
	return qp
}

// WithOp sets the default operator
func (qp *StandardQueryParser) WithOp(op string) *StandardQueryParser {
	qp.op = op
	return qp
}

// WithDf sets the default field
func (qp *StandardQueryParser) WithDf(df string) *StandardQueryParser {
	qp.df = df
	return qp
}

// Sow enables split-white-space
func (qp *StandardQueryParser) Sow() *StandardQueryParser {
	qp.sow = true
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

// BuildParser implements the QueryParserInterface
func (qp *DisMaxQueryParser) BuildParser() (string, error) {
	if qp.q == "" {
		return "", errors.New("'q' parameter is required")
	}

	kv := []string{"dismax"}
	if qp.alt != "" {
		kv = append(kv, "q.alt='"+qp.alt+"'")
	}

	if qp.qf != "" {
		kv = append(kv, "qf='"+qp.qf+"'")
	}

	if qp.mm != "" {
		kv = append(kv, "mm='"+qp.mm+"'")
	}

	if qp.pf != "" {
		kv = append(kv, "qf='"+qp.pf+"'")
	}

	if qp.ps != "" {
		kv = append(kv, "ps='"+qp.ps+"'")
	}

	if qp.qs != "" {
		kv = append(kv, "qs='"+qp.qs+"'")
	}

	if qp.tie != "" {
		kv = append(kv, "tie='"+qp.tie+"'")
	}

	if qp.bq != "" {
		kv = append(kv, "bq='"+qp.bq+"'")
	}

	if qp.bf != "" {
		kv = append(kv, "bf='"+qp.bf+"'")
	}

	return fmt.Sprintf("{!%s}%s", strings.Join(kv, " "), qp.q), nil
}

// NewDisMaxQueryParser returns a new dismax query parser
func NewDisMaxQueryParser(q string) *DisMaxQueryParser {
	return &DisMaxQueryParser{q: q}
}

// WithQ sets the query param
func (qp *DisMaxQueryParser) WithQ(q string) *DisMaxQueryParser {
	qp.q = q
	return qp
}

// WithAlt sets the q.alt param
func (qp *DisMaxQueryParser) WithAlt(alt string) *DisMaxQueryParser {
	qp.alt = alt
	return qp
}

// WithQf sets the qf param
func (qp *DisMaxQueryParser) WithQf(qf string) *DisMaxQueryParser {
	qp.qf = qf
	return qp
}

// WithMm sets the minimum should match param
func (qp *DisMaxQueryParser) WithMm(mm string) *DisMaxQueryParser {
	qp.mm = mm
	return qp
}

// WithPf sets the phrase field param
func (qp *DisMaxQueryParser) WithPf(pf string) *DisMaxQueryParser {
	qp.pf = pf
	return qp
}

// WithPs sets the phrase slop param
func (qp *DisMaxQueryParser) WithPs(ps string) *DisMaxQueryParser {
	qp.ps = ps
	return qp
}

// WithQs sets the query slop param
func (qp *DisMaxQueryParser) WithQs(qs string) *DisMaxQueryParser {
	qp.qs = qs
	return qp
}

// WithTie sets the tie breaker param param
func (qp *DisMaxQueryParser) WithTie(tie string) *DisMaxQueryParser {
	qp.tie = tie
	return qp
}

// WithBq sets the boost query param
func (qp *DisMaxQueryParser) WithBq(bq string) *DisMaxQueryParser {
	qp.bq = bq
	return qp
}

// WithBf sets the boost function param
func (qp *DisMaxQueryParser) WithBf(bf string) *DisMaxQueryParser {
	qp.bf = bf
	return qp
}

// BlockJoinParentQueryParser is a block-join parent query parser
type BlockJoinParentQueryParser struct {
	which   string
	tag     string
	filters string
	score   string
	query   string
}
