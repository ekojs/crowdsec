// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/crowdsecurity/crowdsec/pkg/database/ent/alert"
	"github.com/crowdsecurity/crowdsec/pkg/database/ent/event"
	"github.com/crowdsecurity/crowdsec/pkg/database/ent/predicate"
)

// EventQuery is the builder for querying Event entities.
type EventQuery struct {
	config
	limit      *int
	offset     *int
	unique     *bool
	order      []OrderFunc
	fields     []string
	predicates []predicate.Event
	// eager-loading edges.
	withOwner *AlertQuery
	withFKs   bool
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the EventQuery builder.
func (eq *EventQuery) Where(ps ...predicate.Event) *EventQuery {
	eq.predicates = append(eq.predicates, ps...)
	return eq
}

// Limit adds a limit step to the query.
func (eq *EventQuery) Limit(limit int) *EventQuery {
	eq.limit = &limit
	return eq
}

// Offset adds an offset step to the query.
func (eq *EventQuery) Offset(offset int) *EventQuery {
	eq.offset = &offset
	return eq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (eq *EventQuery) Unique(unique bool) *EventQuery {
	eq.unique = &unique
	return eq
}

// Order adds an order step to the query.
func (eq *EventQuery) Order(o ...OrderFunc) *EventQuery {
	eq.order = append(eq.order, o...)
	return eq
}

// QueryOwner chains the current query on the "owner" edge.
func (eq *EventQuery) QueryOwner() *AlertQuery {
	query := &AlertQuery{config: eq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := eq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := eq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(event.Table, event.FieldID, selector),
			sqlgraph.To(alert.Table, alert.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, event.OwnerTable, event.OwnerColumn),
		)
		fromU = sqlgraph.SetNeighbors(eq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first Event entity from the query.
// Returns a *NotFoundError when no Event was found.
func (eq *EventQuery) First(ctx context.Context) (*Event, error) {
	nodes, err := eq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{event.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (eq *EventQuery) FirstX(ctx context.Context) *Event {
	node, err := eq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Event ID from the query.
// Returns a *NotFoundError when no Event ID was found.
func (eq *EventQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = eq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{event.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (eq *EventQuery) FirstIDX(ctx context.Context) int {
	id, err := eq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Event entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when exactly one Event entity is not found.
// Returns a *NotFoundError when no Event entities are found.
func (eq *EventQuery) Only(ctx context.Context) (*Event, error) {
	nodes, err := eq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{event.Label}
	default:
		return nil, &NotSingularError{event.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (eq *EventQuery) OnlyX(ctx context.Context) *Event {
	node, err := eq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Event ID in the query.
// Returns a *NotSingularError when exactly one Event ID is not found.
// Returns a *NotFoundError when no entities are found.
func (eq *EventQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = eq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{event.Label}
	default:
		err = &NotSingularError{event.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (eq *EventQuery) OnlyIDX(ctx context.Context) int {
	id, err := eq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Events.
func (eq *EventQuery) All(ctx context.Context) ([]*Event, error) {
	if err := eq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return eq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (eq *EventQuery) AllX(ctx context.Context) []*Event {
	nodes, err := eq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Event IDs.
func (eq *EventQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := eq.Select(event.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (eq *EventQuery) IDsX(ctx context.Context) []int {
	ids, err := eq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (eq *EventQuery) Count(ctx context.Context) (int, error) {
	if err := eq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return eq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (eq *EventQuery) CountX(ctx context.Context) int {
	count, err := eq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (eq *EventQuery) Exist(ctx context.Context) (bool, error) {
	if err := eq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return eq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (eq *EventQuery) ExistX(ctx context.Context) bool {
	exist, err := eq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the EventQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (eq *EventQuery) Clone() *EventQuery {
	if eq == nil {
		return nil
	}
	return &EventQuery{
		config:     eq.config,
		limit:      eq.limit,
		offset:     eq.offset,
		order:      append([]OrderFunc{}, eq.order...),
		predicates: append([]predicate.Event{}, eq.predicates...),
		withOwner:  eq.withOwner.Clone(),
		// clone intermediate query.
		sql:  eq.sql.Clone(),
		path: eq.path,
	}
}

// WithOwner tells the query-builder to eager-load the nodes that are connected to
// the "owner" edge. The optional arguments are used to configure the query builder of the edge.
func (eq *EventQuery) WithOwner(opts ...func(*AlertQuery)) *EventQuery {
	query := &AlertQuery{config: eq.config}
	for _, opt := range opts {
		opt(query)
	}
	eq.withOwner = query
	return eq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		CreatedAt time.Time `json:"created_at,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.Event.Query().
//		GroupBy(event.FieldCreatedAt).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (eq *EventQuery) GroupBy(field string, fields ...string) *EventGroupBy {
	group := &EventGroupBy{config: eq.config}
	group.fields = append([]string{field}, fields...)
	group.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := eq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return eq.sqlQuery(ctx), nil
	}
	return group
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		CreatedAt time.Time `json:"created_at,omitempty"`
//	}
//
//	client.Event.Query().
//		Select(event.FieldCreatedAt).
//		Scan(ctx, &v)
//
func (eq *EventQuery) Select(fields ...string) *EventSelect {
	eq.fields = append(eq.fields, fields...)
	return &EventSelect{EventQuery: eq}
}

func (eq *EventQuery) prepareQuery(ctx context.Context) error {
	for _, f := range eq.fields {
		if !event.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if eq.path != nil {
		prev, err := eq.path(ctx)
		if err != nil {
			return err
		}
		eq.sql = prev
	}
	return nil
}

func (eq *EventQuery) sqlAll(ctx context.Context) ([]*Event, error) {
	var (
		nodes       = []*Event{}
		withFKs     = eq.withFKs
		_spec       = eq.querySpec()
		loadedTypes = [1]bool{
			eq.withOwner != nil,
		}
	)
	if eq.withOwner != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, event.ForeignKeys...)
	}
	_spec.ScanValues = func(columns []string) ([]interface{}, error) {
		node := &Event{config: eq.config}
		nodes = append(nodes, node)
		return node.scanValues(columns)
	}
	_spec.Assign = func(columns []string, values []interface{}) error {
		if len(nodes) == 0 {
			return fmt.Errorf("ent: Assign called without calling ScanValues")
		}
		node := nodes[len(nodes)-1]
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	if err := sqlgraph.QueryNodes(ctx, eq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}

	if query := eq.withOwner; query != nil {
		ids := make([]int, 0, len(nodes))
		nodeids := make(map[int][]*Event)
		for i := range nodes {
			if nodes[i].alert_events == nil {
				continue
			}
			fk := *nodes[i].alert_events
			if _, ok := nodeids[fk]; !ok {
				ids = append(ids, fk)
			}
			nodeids[fk] = append(nodeids[fk], nodes[i])
		}
		query.Where(alert.IDIn(ids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "alert_events" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.Owner = n
			}
		}
	}

	return nodes, nil
}

func (eq *EventQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := eq.querySpec()
	_spec.Node.Columns = eq.fields
	if len(eq.fields) > 0 {
		_spec.Unique = eq.unique != nil && *eq.unique
	}
	return sqlgraph.CountNodes(ctx, eq.driver, _spec)
}

func (eq *EventQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := eq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %w", err)
	}
	return n > 0, nil
}

func (eq *EventQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   event.Table,
			Columns: event.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: event.FieldID,
			},
		},
		From:   eq.sql,
		Unique: true,
	}
	if unique := eq.unique; unique != nil {
		_spec.Unique = *unique
	}
	if fields := eq.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, event.FieldID)
		for i := range fields {
			if fields[i] != event.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := eq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := eq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := eq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := eq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (eq *EventQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(eq.driver.Dialect())
	t1 := builder.Table(event.Table)
	columns := eq.fields
	if len(columns) == 0 {
		columns = event.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if eq.sql != nil {
		selector = eq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if eq.unique != nil && *eq.unique {
		selector.Distinct()
	}
	for _, p := range eq.predicates {
		p(selector)
	}
	for _, p := range eq.order {
		p(selector)
	}
	if offset := eq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := eq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// EventGroupBy is the group-by builder for Event entities.
type EventGroupBy struct {
	config
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (egb *EventGroupBy) Aggregate(fns ...AggregateFunc) *EventGroupBy {
	egb.fns = append(egb.fns, fns...)
	return egb
}

// Scan applies the group-by query and scans the result into the given value.
func (egb *EventGroupBy) Scan(ctx context.Context, v interface{}) error {
	query, err := egb.path(ctx)
	if err != nil {
		return err
	}
	egb.sql = query
	return egb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (egb *EventGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := egb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by.
// It is only allowed when executing a group-by query with one field.
func (egb *EventGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(egb.fields) > 1 {
		return nil, errors.New("ent: EventGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := egb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (egb *EventGroupBy) StringsX(ctx context.Context) []string {
	v, err := egb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// String returns a single string from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (egb *EventGroupBy) String(ctx context.Context) (_ string, err error) {
	var v []string
	if v, err = egb.Strings(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{event.Label}
	default:
		err = fmt.Errorf("ent: EventGroupBy.Strings returned %d results when one was expected", len(v))
	}
	return
}

// StringX is like String, but panics if an error occurs.
func (egb *EventGroupBy) StringX(ctx context.Context) string {
	v, err := egb.String(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by.
// It is only allowed when executing a group-by query with one field.
func (egb *EventGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(egb.fields) > 1 {
		return nil, errors.New("ent: EventGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := egb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (egb *EventGroupBy) IntsX(ctx context.Context) []int {
	v, err := egb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Int returns a single int from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (egb *EventGroupBy) Int(ctx context.Context) (_ int, err error) {
	var v []int
	if v, err = egb.Ints(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{event.Label}
	default:
		err = fmt.Errorf("ent: EventGroupBy.Ints returned %d results when one was expected", len(v))
	}
	return
}

// IntX is like Int, but panics if an error occurs.
func (egb *EventGroupBy) IntX(ctx context.Context) int {
	v, err := egb.Int(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by.
// It is only allowed when executing a group-by query with one field.
func (egb *EventGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(egb.fields) > 1 {
		return nil, errors.New("ent: EventGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := egb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (egb *EventGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := egb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64 returns a single float64 from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (egb *EventGroupBy) Float64(ctx context.Context) (_ float64, err error) {
	var v []float64
	if v, err = egb.Float64s(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{event.Label}
	default:
		err = fmt.Errorf("ent: EventGroupBy.Float64s returned %d results when one was expected", len(v))
	}
	return
}

// Float64X is like Float64, but panics if an error occurs.
func (egb *EventGroupBy) Float64X(ctx context.Context) float64 {
	v, err := egb.Float64(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by.
// It is only allowed when executing a group-by query with one field.
func (egb *EventGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(egb.fields) > 1 {
		return nil, errors.New("ent: EventGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := egb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (egb *EventGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := egb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bool returns a single bool from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (egb *EventGroupBy) Bool(ctx context.Context) (_ bool, err error) {
	var v []bool
	if v, err = egb.Bools(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{event.Label}
	default:
		err = fmt.Errorf("ent: EventGroupBy.Bools returned %d results when one was expected", len(v))
	}
	return
}

// BoolX is like Bool, but panics if an error occurs.
func (egb *EventGroupBy) BoolX(ctx context.Context) bool {
	v, err := egb.Bool(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (egb *EventGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	for _, f := range egb.fields {
		if !event.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("invalid field %q for group-by", f)}
		}
	}
	selector := egb.sqlQuery()
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := egb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (egb *EventGroupBy) sqlQuery() *sql.Selector {
	selector := egb.sql.Select()
	aggregation := make([]string, 0, len(egb.fns))
	for _, fn := range egb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	// If no columns were selected in a custom aggregation function, the default
	// selection is the fields used for "group-by", and the aggregation functions.
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(egb.fields)+len(egb.fns))
		for _, f := range egb.fields {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	return selector.GroupBy(selector.Columns(egb.fields...)...)
}

// EventSelect is the builder for selecting fields of Event entities.
type EventSelect struct {
	*EventQuery
	// intermediate query (i.e. traversal path).
	sql *sql.Selector
}

// Scan applies the selector query and scans the result into the given value.
func (es *EventSelect) Scan(ctx context.Context, v interface{}) error {
	if err := es.prepareQuery(ctx); err != nil {
		return err
	}
	es.sql = es.EventQuery.sqlQuery(ctx)
	return es.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (es *EventSelect) ScanX(ctx context.Context, v interface{}) {
	if err := es.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from a selector. It is only allowed when selecting one field.
func (es *EventSelect) Strings(ctx context.Context) ([]string, error) {
	if len(es.fields) > 1 {
		return nil, errors.New("ent: EventSelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := es.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (es *EventSelect) StringsX(ctx context.Context) []string {
	v, err := es.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// String returns a single string from a selector. It is only allowed when selecting one field.
func (es *EventSelect) String(ctx context.Context) (_ string, err error) {
	var v []string
	if v, err = es.Strings(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{event.Label}
	default:
		err = fmt.Errorf("ent: EventSelect.Strings returned %d results when one was expected", len(v))
	}
	return
}

// StringX is like String, but panics if an error occurs.
func (es *EventSelect) StringX(ctx context.Context) string {
	v, err := es.String(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from a selector. It is only allowed when selecting one field.
func (es *EventSelect) Ints(ctx context.Context) ([]int, error) {
	if len(es.fields) > 1 {
		return nil, errors.New("ent: EventSelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := es.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (es *EventSelect) IntsX(ctx context.Context) []int {
	v, err := es.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Int returns a single int from a selector. It is only allowed when selecting one field.
func (es *EventSelect) Int(ctx context.Context) (_ int, err error) {
	var v []int
	if v, err = es.Ints(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{event.Label}
	default:
		err = fmt.Errorf("ent: EventSelect.Ints returned %d results when one was expected", len(v))
	}
	return
}

// IntX is like Int, but panics if an error occurs.
func (es *EventSelect) IntX(ctx context.Context) int {
	v, err := es.Int(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from a selector. It is only allowed when selecting one field.
func (es *EventSelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(es.fields) > 1 {
		return nil, errors.New("ent: EventSelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := es.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (es *EventSelect) Float64sX(ctx context.Context) []float64 {
	v, err := es.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64 returns a single float64 from a selector. It is only allowed when selecting one field.
func (es *EventSelect) Float64(ctx context.Context) (_ float64, err error) {
	var v []float64
	if v, err = es.Float64s(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{event.Label}
	default:
		err = fmt.Errorf("ent: EventSelect.Float64s returned %d results when one was expected", len(v))
	}
	return
}

// Float64X is like Float64, but panics if an error occurs.
func (es *EventSelect) Float64X(ctx context.Context) float64 {
	v, err := es.Float64(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from a selector. It is only allowed when selecting one field.
func (es *EventSelect) Bools(ctx context.Context) ([]bool, error) {
	if len(es.fields) > 1 {
		return nil, errors.New("ent: EventSelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := es.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (es *EventSelect) BoolsX(ctx context.Context) []bool {
	v, err := es.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bool returns a single bool from a selector. It is only allowed when selecting one field.
func (es *EventSelect) Bool(ctx context.Context) (_ bool, err error) {
	var v []bool
	if v, err = es.Bools(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{event.Label}
	default:
		err = fmt.Errorf("ent: EventSelect.Bools returned %d results when one was expected", len(v))
	}
	return
}

// BoolX is like Bool, but panics if an error occurs.
func (es *EventSelect) BoolX(ctx context.Context) bool {
	v, err := es.Bool(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (es *EventSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := es.sql.Query()
	if err := es.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
