package index

import (
	"encoding/json"
)

// Docs contains the list of documents to be indexed
type Docs struct {
	docs []interface{}
}

// NewDocs is a factory for *Docs
func NewDocs() *Docs {
	return &Docs{docs: []interface{}{}}
}

// AddDoc adds a document to the list of documents
func (d *Docs) AddDoc(doc interface{}) {
	d.docs = append(d.docs, doc)
}

// Count counts the number of documents
func (d *Docs) Count() int {
	return len(d.docs)
}

// Marshal marshals the documents
func (d *Docs) Marshal() ([]byte, error) {
	return json.Marshal(d.docs)
}
