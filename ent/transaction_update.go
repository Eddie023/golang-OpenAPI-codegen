// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/eddie023/wex-tag/ent/predicate"
	"github.com/eddie023/wex-tag/ent/transaction"
	"github.com/shopspring/decimal"
)

// TransactionUpdate is the builder for updating Transaction entities.
type TransactionUpdate struct {
	config
	hooks    []Hook
	mutation *TransactionMutation
}

// Where appends a list predicates to the TransactionUpdate builder.
func (tu *TransactionUpdate) Where(ps ...predicate.Transaction) *TransactionUpdate {
	tu.mutation.Where(ps...)
	return tu
}

// SetDate sets the "date" field.
func (tu *TransactionUpdate) SetDate(t time.Time) *TransactionUpdate {
	tu.mutation.SetDate(t)
	return tu
}

// SetNillableDate sets the "date" field if the given value is not nil.
func (tu *TransactionUpdate) SetNillableDate(t *time.Time) *TransactionUpdate {
	if t != nil {
		tu.SetDate(*t)
	}
	return tu
}

// SetAmountInUsd sets the "amount_in_usd" field.
func (tu *TransactionUpdate) SetAmountInUsd(d decimal.Decimal) *TransactionUpdate {
	tu.mutation.ResetAmountInUsd()
	tu.mutation.SetAmountInUsd(d)
	return tu
}

// SetNillableAmountInUsd sets the "amount_in_usd" field if the given value is not nil.
func (tu *TransactionUpdate) SetNillableAmountInUsd(d *decimal.Decimal) *TransactionUpdate {
	if d != nil {
		tu.SetAmountInUsd(*d)
	}
	return tu
}

// AddAmountInUsd adds d to the "amount_in_usd" field.
func (tu *TransactionUpdate) AddAmountInUsd(d decimal.Decimal) *TransactionUpdate {
	tu.mutation.AddAmountInUsd(d)
	return tu
}

// SetDescription sets the "description" field.
func (tu *TransactionUpdate) SetDescription(s string) *TransactionUpdate {
	tu.mutation.SetDescription(s)
	return tu
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (tu *TransactionUpdate) SetNillableDescription(s *string) *TransactionUpdate {
	if s != nil {
		tu.SetDescription(*s)
	}
	return tu
}

// Mutation returns the TransactionMutation object of the builder.
func (tu *TransactionUpdate) Mutation() *TransactionMutation {
	return tu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (tu *TransactionUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, tu.sqlSave, tu.mutation, tu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (tu *TransactionUpdate) SaveX(ctx context.Context) int {
	affected, err := tu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (tu *TransactionUpdate) Exec(ctx context.Context) error {
	_, err := tu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tu *TransactionUpdate) ExecX(ctx context.Context) {
	if err := tu.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (tu *TransactionUpdate) check() error {
	if v, ok := tu.mutation.Description(); ok {
		if err := transaction.DescriptionValidator(v); err != nil {
			return &ValidationError{Name: "description", err: fmt.Errorf(`ent: validator failed for field "Transaction.description": %w`, err)}
		}
	}
	return nil
}

func (tu *TransactionUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := tu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(transaction.Table, transaction.Columns, sqlgraph.NewFieldSpec(transaction.FieldID, field.TypeUUID))
	if ps := tu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := tu.mutation.Date(); ok {
		_spec.SetField(transaction.FieldDate, field.TypeTime, value)
	}
	if value, ok := tu.mutation.AmountInUsd(); ok {
		_spec.SetField(transaction.FieldAmountInUsd, field.TypeFloat64, value)
	}
	if value, ok := tu.mutation.AddedAmountInUsd(); ok {
		_spec.AddField(transaction.FieldAmountInUsd, field.TypeFloat64, value)
	}
	if value, ok := tu.mutation.Description(); ok {
		_spec.SetField(transaction.FieldDescription, field.TypeString, value)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, tu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{transaction.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	tu.mutation.done = true
	return n, nil
}

// TransactionUpdateOne is the builder for updating a single Transaction entity.
type TransactionUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *TransactionMutation
}

// SetDate sets the "date" field.
func (tuo *TransactionUpdateOne) SetDate(t time.Time) *TransactionUpdateOne {
	tuo.mutation.SetDate(t)
	return tuo
}

// SetNillableDate sets the "date" field if the given value is not nil.
func (tuo *TransactionUpdateOne) SetNillableDate(t *time.Time) *TransactionUpdateOne {
	if t != nil {
		tuo.SetDate(*t)
	}
	return tuo
}

// SetAmountInUsd sets the "amount_in_usd" field.
func (tuo *TransactionUpdateOne) SetAmountInUsd(d decimal.Decimal) *TransactionUpdateOne {
	tuo.mutation.ResetAmountInUsd()
	tuo.mutation.SetAmountInUsd(d)
	return tuo
}

// SetNillableAmountInUsd sets the "amount_in_usd" field if the given value is not nil.
func (tuo *TransactionUpdateOne) SetNillableAmountInUsd(d *decimal.Decimal) *TransactionUpdateOne {
	if d != nil {
		tuo.SetAmountInUsd(*d)
	}
	return tuo
}

// AddAmountInUsd adds d to the "amount_in_usd" field.
func (tuo *TransactionUpdateOne) AddAmountInUsd(d decimal.Decimal) *TransactionUpdateOne {
	tuo.mutation.AddAmountInUsd(d)
	return tuo
}

// SetDescription sets the "description" field.
func (tuo *TransactionUpdateOne) SetDescription(s string) *TransactionUpdateOne {
	tuo.mutation.SetDescription(s)
	return tuo
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (tuo *TransactionUpdateOne) SetNillableDescription(s *string) *TransactionUpdateOne {
	if s != nil {
		tuo.SetDescription(*s)
	}
	return tuo
}

// Mutation returns the TransactionMutation object of the builder.
func (tuo *TransactionUpdateOne) Mutation() *TransactionMutation {
	return tuo.mutation
}

// Where appends a list predicates to the TransactionUpdate builder.
func (tuo *TransactionUpdateOne) Where(ps ...predicate.Transaction) *TransactionUpdateOne {
	tuo.mutation.Where(ps...)
	return tuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (tuo *TransactionUpdateOne) Select(field string, fields ...string) *TransactionUpdateOne {
	tuo.fields = append([]string{field}, fields...)
	return tuo
}

// Save executes the query and returns the updated Transaction entity.
func (tuo *TransactionUpdateOne) Save(ctx context.Context) (*Transaction, error) {
	return withHooks(ctx, tuo.sqlSave, tuo.mutation, tuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (tuo *TransactionUpdateOne) SaveX(ctx context.Context) *Transaction {
	node, err := tuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (tuo *TransactionUpdateOne) Exec(ctx context.Context) error {
	_, err := tuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tuo *TransactionUpdateOne) ExecX(ctx context.Context) {
	if err := tuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (tuo *TransactionUpdateOne) check() error {
	if v, ok := tuo.mutation.Description(); ok {
		if err := transaction.DescriptionValidator(v); err != nil {
			return &ValidationError{Name: "description", err: fmt.Errorf(`ent: validator failed for field "Transaction.description": %w`, err)}
		}
	}
	return nil
}

func (tuo *TransactionUpdateOne) sqlSave(ctx context.Context) (_node *Transaction, err error) {
	if err := tuo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(transaction.Table, transaction.Columns, sqlgraph.NewFieldSpec(transaction.FieldID, field.TypeUUID))
	id, ok := tuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Transaction.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := tuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, transaction.FieldID)
		for _, f := range fields {
			if !transaction.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != transaction.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := tuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := tuo.mutation.Date(); ok {
		_spec.SetField(transaction.FieldDate, field.TypeTime, value)
	}
	if value, ok := tuo.mutation.AmountInUsd(); ok {
		_spec.SetField(transaction.FieldAmountInUsd, field.TypeFloat64, value)
	}
	if value, ok := tuo.mutation.AddedAmountInUsd(); ok {
		_spec.AddField(transaction.FieldAmountInUsd, field.TypeFloat64, value)
	}
	if value, ok := tuo.mutation.Description(); ok {
		_spec.SetField(transaction.FieldDescription, field.TypeString, value)
	}
	_node = &Transaction{config: tuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, tuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{transaction.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	tuo.mutation.done = true
	return _node, nil
}
