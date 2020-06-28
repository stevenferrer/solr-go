package schema

// Schema is a schema
type Schema struct {
	Name       string      `json:"name"`
	Version    float64     `json:"version"`
	UniqueKey  string      `json:"uniqueKey"`
	FieldTypes []FieldType `json:"fieldTypes,omitempty"`
	Fields     []Field     `json:"fields,omitempty"`
	CopyFields []CopyField `json:"copyFields"`
}

// FieldType is a field type
type FieldType struct {
	Name                      string    `json:"name"`
	Class                     string    `json:"class"`
	PositionIncrementGap      string    `json:"positionIncrementGap,omitempty"`
	AutoGeneratePhraseQueries string    `json:"autoGeneratePhraseQueries,omitempty"`
	SynonymQueryStyle         string    `json:"synonymQueryStyle,omitempty"`
	EnableGraphQueries        string    `json:"enableGraphQueries,omitempty"`
	DocValuesFormat           string    `json:"docValuesFormat,omitempty"`
	PostingsFormat            string    `json:"postingsFormat,omitempty"`
	Indexed                   bool      `json:"indexed,omitempty"`
	Stored                    bool      `json:"stored,omitempty"`
	DocValues                 bool      `json:"docValues,omitempty"`
	SortMissingFirst          bool      `json:"sortMissingFirst,omitempty"`
	SortMissingLast           bool      `json:"sortMissingLast,omitempty"`
	MultiValued               bool      `json:"multiValued,omitempty"`
	Uninvertible              bool      `json:"uninvertible,omitempty"`
	OmitNorms                 bool      `json:"omitNorms,omitempty"`
	OmitTermFreqAndPositions  bool      `json:"omitTermFreqAndPositions,omitempty"`
	OmitPositions             bool      `json:"omitPositions,omitempty"`
	TermVectors               bool      `json:"termVectors,omitempty"`
	TermPositions             bool      `json:"termPositions,omitempty"`
	TermOffsets               bool      `json:"termOffsets,omitempty"`
	TermPayloads              bool      `json:"termPayloads,omitempty"`
	Required                  bool      `json:"required,omitempty"`
	UseDocValuesAsStored      bool      `json:"useDocValuesAsStored,omitempty"`
	Large                     bool      `json:"large,omitempty"`
	MaxCharsForDocValues      string    `json:"maxCharsForDocValues,omitempty"`
	Geo                       string    `json:"geo,omitempty"`
	MaxDistErr                string    `json:"maxDistErr,omitempty"`
	DistErrPct                string    `json:"distErrPct,omitempty"`
	DistanceUnits             string    `json:"distanceUnits,omitempty"`
	SubFieldSuffix            string    `json:"subFieldSuffix,omitempty"`
	Dimension                 string    `json:"dimension,omitempty"`
	Analyzer                  *Analyzer `json:"analyzer,omitempty"`
	IndexAnalyzer             *Analyzer `json:"indexAnalyzer,omitempty"`
	QueryAnalyzer             *Analyzer `json:"queryAnalyzer,omitempty"`
}

// Tokenizer is a tokenizer
type Tokenizer struct {
	Class                 string `json:"class"`
	Delimeter             string `json:"delimiter,omitempty"`
	OutputUnknownUnigrams string `json:"outputUnknownUnigrams,omitempty"`
	DecompoundMode        string `json:"decompoundMode,omitempty"`
	Mode                  string `json:"mode,omitempty"`
}

// Analyzer is an analyzer
type Analyzer struct {
	Tokenizer   *Tokenizer `json:"tokenizer"`
	Filters     []Filter   `json:"filters,omitempty"`
	CharFilters []Filter   `json:"charFilters,omitempty"`
}

// Filter is a filter
type Filter struct {
	Class               string `json:"class"`
	Encoder             string `json:"encoder,omitempty"`
	Inject              string `json:"inject,omitempty"`
	Words               string `json:"words,omitempty"`
	IgnoreCase          string `json:"ignoreCase,omitempty"`
	Articles            string `json:"articles,omitempty"`
	Language            string `json:"language,omitempty"`
	Format              string `json:"format,omitempty"`
	Protected           string `json:"protected,omitempty"`
	Expand              string `json:"expand,omitempty"`
	Synonyms            string `json:"synonyms,omitempty"`
	CatenateNumbers     string `json:"catenateNumbers,omitempty"`
	GenerateNumberParts string `json:"generateNumberParts,omitempty"`
	SplitOnCaseChange   string `json:"splitOnCaseChange,omitempty"`
	GenerateWordParts   string `json:"generateWordParts,omitempty"`
	CatenateAll         string `json:"catenateAll,omitempty"`
	CatenateWords       string `json:"catenateWords,omitempty"`
	MaxPosQuestion      string `json:"maxPosQuestion,omitempty"`
	MaxFractionAsterisk string `json:"maxFractionAsterisk,omitempty"`
	MaxPosAsterisk      string `json:"maxPosAsterisk,omitempty"`
	WithOriginal        string `json:"withOriginal,omitempty"`
	StemDerivational    string `json:"stemDerivational,omitempty"`
	MinimumLength       string `json:"minimumLength,omitempty"`
	Dictionary          string `json:"dictionary,omitempty"`
	Tags                string `json:"tags,omitempty"`
	Replacement         string `json:"replacement,omitempty"`
	Pattern             string `json:"pattern,omitempty"`
	PreserveOriginal    string `json:"preserveOriginal,omitempty"`
	Replace             string `json:"replace,omitempty"`
}

// Field is a field
type Field struct {
	Name                 string `json:"name"`
	Type                 string `json:"type"`
	DocValues            bool   `json:"docValues,omitempty"`
	Indexed              bool   `json:"indexed,omitempty"`
	Stored               bool   `json:"stored,omitempty"`
	MultiValued          bool   `json:"multiValued,omitempty"`
	Required             bool   `json:"requied,omitempty"`
	UseDocValuesAsStored bool   `json:"useDocValuesAsStored,omitempty"`
}

// CopyField is a copy field
type CopyField struct {
	Source   string `json:"source,omitempty"`
	Dest     string `json:"dest,omitempty"`
	MaxChars int    `json:"maxchars,omitempty"`
}
