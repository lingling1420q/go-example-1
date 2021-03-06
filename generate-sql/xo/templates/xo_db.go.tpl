// XODB is the common interface for database operations that can be used with
// types from schema '{{ schema .Schema }}'.
//
// This should work with database/sql.DB and database/sql.Tx.
type XODB interface {
	Exec(string, ...interface{}) (sql.Result, error)
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryRow(string, ...interface{}) *sql.Row
}

// XOLog provides the log func used by generated queries.
var XOLog = func(string, ...interface{}) { }

// ScannerValuer is the common interface for types that implement both the
// database/sql.Scanner and sql/driver.Valuer interfaces.
type ScannerValuer interface {
	sql.Scanner
	driver.Valuer
}

// StringSlice is a slice of strings.
type StringSlice []string

// Cursor specifies an index to sort by, the direction of the sort, an offset, and a limit.
type Cursor struct {
	Offset *int32
	Limit  *int32
	Index  *string
	Desc   *bool
	Dead   *bool
	After  *graphql.ID
	First  *int32
	Before *graphql.ID
	Last   *int32
}

var (
	defaultOffset int32  = 0
	defaultLimit  int32  = 50
	defaultIndex  string = "id"
	defaultDesc   bool   = false
	defaultDead   bool   = false
)

// DefaultCursor will get the 50 first non-deleted IDs from a table.
var DefaultCursor = Cursor{
	Offset: &defaultOffset,
	Limit:  &defaultLimit,
	Index:  &defaultIndex,
	Desc:   &defaultDesc,
	Dead:   &defaultDead,
}

// sqlConjunctionMap supported conjunction, related to graphql enum: FilterConjunction
var sqlConjunctionMap = map[string]struct{}{
	"AND":{},
	"OR":{},
}

// filterPair item of filter
type filterPair struct{
	fieldName string
	option string
	value interface{}
}

// filterArguments filter arguments
type filterArguments struct{
	filterPairs []*filterPair
	conjunction string
	conjCnt int
}

// quoteEscapeRegex is the regex to match escaped characters in a string.
var quoteEscapeRegex = regexp.MustCompile(`([^\\]([\\]{2})*)\\"`)

// Scan satisfies the sql.Scanner interface for StringSlice.
func (ss *StringSlice) Scan(src interface{}) error {
	buf, ok := src.([]byte)
	if !ok {
		return errors.New("invalid StringSlice")
	}

	// change quote escapes for csv parser
	str := quoteEscapeRegex.ReplaceAllString(string(buf), `$1""`)
	str = strings.Replace(str, `\\`, `\`, -1)

	// remove braces
	str = str[1:len(str)-1]

	// bail if only one
	if len(str) == 0 {
		*ss = StringSlice([]string{})
		return nil
	}

	// parse with csv reader
	cr := csv.NewReader(strings.NewReader(str))
	slice, err := cr.Read()
	if err != nil {
		fmt.Printf("exiting!: %v\n", err)
		return err
	}

	*ss = StringSlice(slice)

	return nil
}

// Value satisfies the driver.Valuer interface for StringSlice.
func (ss StringSlice) Value() (driver.Value, error) {
	v := make([]string, len(ss))
	for i, s := range ss {
		v[i] = `"` + strings.Replace(strings.Replace(s, `\`, `\\\`, -1), `"`, `\"`, -1) + `"`
	}
	return "{" + strings.Join(v, ",") + "}", nil
}

// Slice is a slice of ScannerValuers.
type Slice []ScannerValuer


// GraphQL related types
const GraphQLCommonTypes = `
	type PageInfo {
		hasNextPage: Boolean!
		hasPreviousPage: Boolean!
		startCursor: ID
		endCursor: ID
	}
	scalar Time
	enum FilterConjunction{
		AND
		OR
	}
`

// PageInfoResolver defines the GraphQL PageInfo type
type PageInfoResolver struct {
	startCursor     graphql.ID
	endCursor       graphql.ID
	hasNextPage     bool
	hasPreviousPage bool
}

// StartCursor returns the start cursor (global id)
func (r PageInfoResolver) StartCursor() *graphql.ID {
	return &r.startCursor
}

// EndCursor returns the end cursor (global id)
func (r PageInfoResolver) EndCursor() *graphql.ID {
	return &r.endCursor
}

// HasNextPage returns if next page is available
func (r PageInfoResolver) HasNextPage() bool {
	return r.hasNextPage
}

// HasPreviousPage returns if previous page is available
func (r PageInfoResolver) HasPreviousPage() bool {
	return r.hasNextPage
}

func encodeCursor(typeName string, id int) graphql.ID {
	return graphql.ID(base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%d", typeName, id))))
}

type dbContext struct{}

// DBCtx is the key for setting DB Context.WithValue
var DBCtx = dbContext{}

// Bool returns a nullable bool.
func Bool(b bool) sql.NullBool {
	return sql.NullBool{Bool: b, Valid: true}
}

// BoolPointer converts bool pointer to sql.NullBool
func BoolPointer(b *bool) sql.NullBool {
	if b == nil {
		return sql.NullBool{}
	}
	return sql.NullBool{Bool: *b, Valid: true}
}

// PointerBool converts bool to pointer to bool
func PointerBool(b sql.NullBool) *bool {
	if !b.Valid {
		return nil
	}
	return &b.Bool
}

// Int64 returns a nullable int64
func Int64(i int64) sql.NullInt64 {
	return sql.NullInt64{Int64: i, Valid: true}
}

// Int64Pointer converts a int64 pointer to sql.NullInt64
func Int64Pointer(i *int64) sql.NullInt64 {
	if i == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: *i, Valid: true}
}

// PointerInt64 converts sql.NullInt64 to pointer to int64
func PointerInt64(i sql.NullInt64) *int64 {
	if !i.Valid {
		return nil
	}
	return &i.Int64
}

// Float64 returns a nullable float64
func Float64(i float64) sql.NullFloat64 {
	return sql.NullFloat64{Float64: i, Valid: true}
}

// Float64Pointer converts a float64 pointer to sql.NullFloat64
func Float64Pointer(i *float64) sql.NullFloat64 {
	if i == nil {
		return sql.NullFloat64{}
	}
	return sql.NullFloat64{Float64: *i, Valid: true}
}

// PointerFloat64 converts sql.NullFloat64 to pointer to float64
func PointerFloat64(i sql.NullFloat64) *float64 {
	if !i.Valid {
		return nil
	}
	return &i.Float64
}

// String returns a nullable string
func String(s string) sql.NullString {
	return sql.NullString{String: s, Valid: true}
}

// StringPointer converts string pointer to sql.NullString
func StringPointer(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *s, Valid: true}
}

// PointerString converts sql.NullString to pointer to string
func PointerString(s sql.NullString) *string {
	if !s.Valid {
		return nil
	}
	return &s.String
}

// Time returns a nullable Time
func Time(t time.Time) pq.NullTime {
	return pq.NullTime{Time: t, Valid: true}
}

// TimePointer converts time.Time pointer to pq.NullTime
func TimePointer(t *time.Time) pq.NullTime {
	if t == nil {
		return pq.NullTime{}
	}
	return pq.NullTime{Time: *t, Valid: true}
}

// TimeGqlPointer converts graphql.Time pointer to pq.NullTime
func TimeGqlPointer(t *graphql.Time) pq.NullTime {
	if t == nil {
		return pq.NullTime{}
	}
	return pq.NullTime{Time: t.Time, Valid: true}
}

// PointerTime converts pq.NullTIme to pointer to time.Time
func PointerTime(t pq.NullTime) *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

// PointerGqlTime converts pq.NullType to pointer to graphql.Time
func PointerGqlTime(t pq.NullTime) *graphql.Time {
	if !t.Valid {
		return nil
	}
	return &graphql.Time{t.Time}
}

// PointerStringInt64 converts Int64 pointer to string pointer
func PointerStringInt64(i *int64) *string {
	if i == nil {
		return nil
	}
	s := strconv.Itoa(int(*i))
	return &s
}

// PointerStringSqlInt64 converts sql.NullInt64 pointer to graphql.ID pointer
func PointerStringSqlInt64(i sql.NullInt64) *string {
	if !i.Valid {
		return nil
	}
	s := strconv.Itoa(int(i.Int64))
	return &s
}

// PointerStringFloat64 converts Float64 pointer to string pointer
func PointerStringFloat64(i *float64) *string {
	if i == nil {
		return nil
	}
	s :=fmt.Sprintf("%.6f", *i)
	return &s
}

// PointerFloat64SqlFloat64 converts sql.NullFloat64 pointer to graphql.ID pointer
func PointerFloat64SqlFloat64(i sql.NullFloat64) *float64 {
	if !i.Valid {
		return nil
	}
	s := i.Float64
	return &s
}

type RootResolver struct{}

func GetRootSchema(extraQueries, extraMutations, extraTypes string) string {
	return `
	schema {
		query: Query
		mutation: Mutation
	}

	type Query {
` +  extraQueries +
{{- range $type, $_ := .TypeMap }}
	Get{{ $type }}Queries() +
{{- end -}}
`}

type Mutation {
` + extraMutations +
{{- range $type, $_ := .TypeMap }}
	Get{{ $type }}Mutations() +
{{- end -}}
`}

` + extraTypes +
{{- range $type, $_ := .TypeMap }}
	GraphQL{{ $type }}Types +
{{- end }}
GraphQLCommonTypes
}
