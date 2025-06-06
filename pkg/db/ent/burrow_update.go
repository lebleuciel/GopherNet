// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"gophernet/pkg/db/ent/burrow"
	"gophernet/pkg/db/ent/predicate"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// BurrowUpdate is the builder for updating Burrow entities.
type BurrowUpdate struct {
	config
	hooks    []Hook
	mutation *BurrowMutation
}

// Where appends a list predicates to the BurrowUpdate builder.
func (bu *BurrowUpdate) Where(ps ...predicate.Burrow) *BurrowUpdate {
	bu.mutation.Where(ps...)
	return bu
}

// SetName sets the "name" field.
func (bu *BurrowUpdate) SetName(s string) *BurrowUpdate {
	bu.mutation.SetName(s)
	return bu
}

// SetNillableName sets the "name" field if the given value is not nil.
func (bu *BurrowUpdate) SetNillableName(s *string) *BurrowUpdate {
	if s != nil {
		bu.SetName(*s)
	}
	return bu
}

// SetDepth sets the "depth" field.
func (bu *BurrowUpdate) SetDepth(f float64) *BurrowUpdate {
	bu.mutation.ResetDepth()
	bu.mutation.SetDepth(f)
	return bu
}

// SetNillableDepth sets the "depth" field if the given value is not nil.
func (bu *BurrowUpdate) SetNillableDepth(f *float64) *BurrowUpdate {
	if f != nil {
		bu.SetDepth(*f)
	}
	return bu
}

// AddDepth adds f to the "depth" field.
func (bu *BurrowUpdate) AddDepth(f float64) *BurrowUpdate {
	bu.mutation.AddDepth(f)
	return bu
}

// SetWidth sets the "width" field.
func (bu *BurrowUpdate) SetWidth(f float64) *BurrowUpdate {
	bu.mutation.ResetWidth()
	bu.mutation.SetWidth(f)
	return bu
}

// SetNillableWidth sets the "width" field if the given value is not nil.
func (bu *BurrowUpdate) SetNillableWidth(f *float64) *BurrowUpdate {
	if f != nil {
		bu.SetWidth(*f)
	}
	return bu
}

// AddWidth adds f to the "width" field.
func (bu *BurrowUpdate) AddWidth(f float64) *BurrowUpdate {
	bu.mutation.AddWidth(f)
	return bu
}

// SetIsOccupied sets the "is_occupied" field.
func (bu *BurrowUpdate) SetIsOccupied(b bool) *BurrowUpdate {
	bu.mutation.SetIsOccupied(b)
	return bu
}

// SetNillableIsOccupied sets the "is_occupied" field if the given value is not nil.
func (bu *BurrowUpdate) SetNillableIsOccupied(b *bool) *BurrowUpdate {
	if b != nil {
		bu.SetIsOccupied(*b)
	}
	return bu
}

// SetAge sets the "age" field.
func (bu *BurrowUpdate) SetAge(i int) *BurrowUpdate {
	bu.mutation.ResetAge()
	bu.mutation.SetAge(i)
	return bu
}

// SetNillableAge sets the "age" field if the given value is not nil.
func (bu *BurrowUpdate) SetNillableAge(i *int) *BurrowUpdate {
	if i != nil {
		bu.SetAge(*i)
	}
	return bu
}

// AddAge adds i to the "age" field.
func (bu *BurrowUpdate) AddAge(i int) *BurrowUpdate {
	bu.mutation.AddAge(i)
	return bu
}

// SetUpdatedAt sets the "updated_at" field.
func (bu *BurrowUpdate) SetUpdatedAt(t time.Time) *BurrowUpdate {
	bu.mutation.SetUpdatedAt(t)
	return bu
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (bu *BurrowUpdate) SetNillableUpdatedAt(t *time.Time) *BurrowUpdate {
	if t != nil {
		bu.SetUpdatedAt(*t)
	}
	return bu
}

// Mutation returns the BurrowMutation object of the builder.
func (bu *BurrowUpdate) Mutation() *BurrowMutation {
	return bu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (bu *BurrowUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, bu.sqlSave, bu.mutation, bu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (bu *BurrowUpdate) SaveX(ctx context.Context) int {
	affected, err := bu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (bu *BurrowUpdate) Exec(ctx context.Context) error {
	_, err := bu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (bu *BurrowUpdate) ExecX(ctx context.Context) {
	if err := bu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (bu *BurrowUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(burrow.Table, burrow.Columns, sqlgraph.NewFieldSpec(burrow.FieldID, field.TypeInt))
	if ps := bu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := bu.mutation.Name(); ok {
		_spec.SetField(burrow.FieldName, field.TypeString, value)
	}
	if value, ok := bu.mutation.Depth(); ok {
		_spec.SetField(burrow.FieldDepth, field.TypeFloat64, value)
	}
	if value, ok := bu.mutation.AddedDepth(); ok {
		_spec.AddField(burrow.FieldDepth, field.TypeFloat64, value)
	}
	if value, ok := bu.mutation.Width(); ok {
		_spec.SetField(burrow.FieldWidth, field.TypeFloat64, value)
	}
	if value, ok := bu.mutation.AddedWidth(); ok {
		_spec.AddField(burrow.FieldWidth, field.TypeFloat64, value)
	}
	if value, ok := bu.mutation.IsOccupied(); ok {
		_spec.SetField(burrow.FieldIsOccupied, field.TypeBool, value)
	}
	if value, ok := bu.mutation.Age(); ok {
		_spec.SetField(burrow.FieldAge, field.TypeInt, value)
	}
	if value, ok := bu.mutation.AddedAge(); ok {
		_spec.AddField(burrow.FieldAge, field.TypeInt, value)
	}
	if value, ok := bu.mutation.UpdatedAt(); ok {
		_spec.SetField(burrow.FieldUpdatedAt, field.TypeTime, value)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, bu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{burrow.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	bu.mutation.done = true
	return n, nil
}

// BurrowUpdateOne is the builder for updating a single Burrow entity.
type BurrowUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *BurrowMutation
}

// SetName sets the "name" field.
func (buo *BurrowUpdateOne) SetName(s string) *BurrowUpdateOne {
	buo.mutation.SetName(s)
	return buo
}

// SetNillableName sets the "name" field if the given value is not nil.
func (buo *BurrowUpdateOne) SetNillableName(s *string) *BurrowUpdateOne {
	if s != nil {
		buo.SetName(*s)
	}
	return buo
}

// SetDepth sets the "depth" field.
func (buo *BurrowUpdateOne) SetDepth(f float64) *BurrowUpdateOne {
	buo.mutation.ResetDepth()
	buo.mutation.SetDepth(f)
	return buo
}

// SetNillableDepth sets the "depth" field if the given value is not nil.
func (buo *BurrowUpdateOne) SetNillableDepth(f *float64) *BurrowUpdateOne {
	if f != nil {
		buo.SetDepth(*f)
	}
	return buo
}

// AddDepth adds f to the "depth" field.
func (buo *BurrowUpdateOne) AddDepth(f float64) *BurrowUpdateOne {
	buo.mutation.AddDepth(f)
	return buo
}

// SetWidth sets the "width" field.
func (buo *BurrowUpdateOne) SetWidth(f float64) *BurrowUpdateOne {
	buo.mutation.ResetWidth()
	buo.mutation.SetWidth(f)
	return buo
}

// SetNillableWidth sets the "width" field if the given value is not nil.
func (buo *BurrowUpdateOne) SetNillableWidth(f *float64) *BurrowUpdateOne {
	if f != nil {
		buo.SetWidth(*f)
	}
	return buo
}

// AddWidth adds f to the "width" field.
func (buo *BurrowUpdateOne) AddWidth(f float64) *BurrowUpdateOne {
	buo.mutation.AddWidth(f)
	return buo
}

// SetIsOccupied sets the "is_occupied" field.
func (buo *BurrowUpdateOne) SetIsOccupied(b bool) *BurrowUpdateOne {
	buo.mutation.SetIsOccupied(b)
	return buo
}

// SetNillableIsOccupied sets the "is_occupied" field if the given value is not nil.
func (buo *BurrowUpdateOne) SetNillableIsOccupied(b *bool) *BurrowUpdateOne {
	if b != nil {
		buo.SetIsOccupied(*b)
	}
	return buo
}

// SetAge sets the "age" field.
func (buo *BurrowUpdateOne) SetAge(i int) *BurrowUpdateOne {
	buo.mutation.ResetAge()
	buo.mutation.SetAge(i)
	return buo
}

// SetNillableAge sets the "age" field if the given value is not nil.
func (buo *BurrowUpdateOne) SetNillableAge(i *int) *BurrowUpdateOne {
	if i != nil {
		buo.SetAge(*i)
	}
	return buo
}

// AddAge adds i to the "age" field.
func (buo *BurrowUpdateOne) AddAge(i int) *BurrowUpdateOne {
	buo.mutation.AddAge(i)
	return buo
}

// SetUpdatedAt sets the "updated_at" field.
func (buo *BurrowUpdateOne) SetUpdatedAt(t time.Time) *BurrowUpdateOne {
	buo.mutation.SetUpdatedAt(t)
	return buo
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (buo *BurrowUpdateOne) SetNillableUpdatedAt(t *time.Time) *BurrowUpdateOne {
	if t != nil {
		buo.SetUpdatedAt(*t)
	}
	return buo
}

// Mutation returns the BurrowMutation object of the builder.
func (buo *BurrowUpdateOne) Mutation() *BurrowMutation {
	return buo.mutation
}

// Where appends a list predicates to the BurrowUpdate builder.
func (buo *BurrowUpdateOne) Where(ps ...predicate.Burrow) *BurrowUpdateOne {
	buo.mutation.Where(ps...)
	return buo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (buo *BurrowUpdateOne) Select(field string, fields ...string) *BurrowUpdateOne {
	buo.fields = append([]string{field}, fields...)
	return buo
}

// Save executes the query and returns the updated Burrow entity.
func (buo *BurrowUpdateOne) Save(ctx context.Context) (*Burrow, error) {
	return withHooks(ctx, buo.sqlSave, buo.mutation, buo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (buo *BurrowUpdateOne) SaveX(ctx context.Context) *Burrow {
	node, err := buo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (buo *BurrowUpdateOne) Exec(ctx context.Context) error {
	_, err := buo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (buo *BurrowUpdateOne) ExecX(ctx context.Context) {
	if err := buo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (buo *BurrowUpdateOne) sqlSave(ctx context.Context) (_node *Burrow, err error) {
	_spec := sqlgraph.NewUpdateSpec(burrow.Table, burrow.Columns, sqlgraph.NewFieldSpec(burrow.FieldID, field.TypeInt))
	id, ok := buo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Burrow.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := buo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, burrow.FieldID)
		for _, f := range fields {
			if !burrow.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != burrow.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := buo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := buo.mutation.Name(); ok {
		_spec.SetField(burrow.FieldName, field.TypeString, value)
	}
	if value, ok := buo.mutation.Depth(); ok {
		_spec.SetField(burrow.FieldDepth, field.TypeFloat64, value)
	}
	if value, ok := buo.mutation.AddedDepth(); ok {
		_spec.AddField(burrow.FieldDepth, field.TypeFloat64, value)
	}
	if value, ok := buo.mutation.Width(); ok {
		_spec.SetField(burrow.FieldWidth, field.TypeFloat64, value)
	}
	if value, ok := buo.mutation.AddedWidth(); ok {
		_spec.AddField(burrow.FieldWidth, field.TypeFloat64, value)
	}
	if value, ok := buo.mutation.IsOccupied(); ok {
		_spec.SetField(burrow.FieldIsOccupied, field.TypeBool, value)
	}
	if value, ok := buo.mutation.Age(); ok {
		_spec.SetField(burrow.FieldAge, field.TypeInt, value)
	}
	if value, ok := buo.mutation.AddedAge(); ok {
		_spec.AddField(burrow.FieldAge, field.TypeInt, value)
	}
	if value, ok := buo.mutation.UpdatedAt(); ok {
		_spec.SetField(burrow.FieldUpdatedAt, field.TypeTime, value)
	}
	_node = &Burrow{config: buo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, buo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{burrow.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	buo.mutation.done = true
	return _node, nil
}
