package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrTransactionAlreadyStarted = errors.New("transaction already started")
	ErrTransactionNotStarted     = errors.New("transaction not started")
	ErrTransactionNotRolledBack  = errors.New("transaction not rolled back")
	ErrTransactionNotCommitted   = errors.New("transaction not committed")
	ErrLockTablesFailed          = errors.New("failed to lock tables")
)

type Worker interface {
	BeginTransaction(ctx context.Context, opts *sql.TxOptions) error
	BeginLockTransaction(ctx context.Context, tables ...string) error
	CommitOrRollback() error
	Rollback() error
	Lock(ctx context.Context, tables ...string) error
	GetTx() *sql.Tx
	DoLock(ctx context.Context, fn func(uow Worker) error, tables ...string) error
	Do(ctx context.Context, fn func(uow Worker) error) error
}

type UnitOfWork struct {
	db *sql.DB
	tx *sql.Tx
}

func NewUnitOfWork(db *sql.DB) Worker {
	return &UnitOfWork{
		db: db,
	}
}

func (u *UnitOfWork) BeginTransaction(ctx context.Context, opts *sql.TxOptions) error {
	if u.tx != nil {
		return ErrTransactionAlreadyStarted
	}

	tx, err := u.db.BeginTx(ctx, opts)
	if err != nil {
		return err
	}
	u.tx = tx

	return nil
}

func (u *UnitOfWork) BeginLockTransaction(ctx context.Context, tables ...string) error {
	err := u.BeginTransaction(ctx, nil)
	if err != nil {
		return err
	}

	err = u.Lock(ctx, tables...)
	if err != nil {
		return err
	}

	return nil
}

func (u *UnitOfWork) Lock(ctx context.Context, tables ...string) error {
	if u.tx == nil {
		return ErrTransactionNotStarted
	}

	q := fmt.Sprintf("LOCK TABLE %s", strings.Join(tables, ", "))

	_, err := u.tx.ExecContext(ctx, q)
	if err != nil {
		return errors.Join(err, ErrLockTablesFailed)
	}

	return nil
}

func (u *UnitOfWork) CommitOrRollback() error {
	if u.tx == nil {
		return ErrTransactionNotStarted
	}

	err := u.tx.Commit()
	if err != nil {
		errRollback := u.Rollback()

		if errRollback != nil {
			return errors.Join(err, ErrTransactionNotRolledBack)
		}

		return errors.Join(err, ErrTransactionNotCommitted)
	}

	u.tx = nil

	return nil
}

func (u *UnitOfWork) Rollback() error {
	if u.tx == nil {
		return ErrTransactionNotStarted
	}

	err := u.tx.Rollback()
	if err != nil {
		return errors.Join(err, ErrTransactionNotRolledBack)
	}

	u.tx = nil

	return nil
}

func (u *UnitOfWork) GetTx() *sql.Tx {
	return u.tx
}

func (u *UnitOfWork) DoLock(ctx context.Context, fn func(uow Worker) error, tables ...string) error {
	err := u.BeginLockTransaction(ctx, tables...)
	if err != nil {
		return err
	}

	err = fn(u)

	if err != nil {
		errRollback := u.Rollback()

		if errRollback != nil {
			return errors.Join(err, ErrTransactionNotRolledBack)
		}

		return err
	}

	return u.CommitOrRollback()
}

func (u *UnitOfWork) Do(ctx context.Context, fn func(uow Worker) error) error {
	err := u.BeginTransaction(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}

	err = fn(u)

	if err != nil {
		errRollback := u.Rollback()

		if errRollback != nil {
			return errors.Join(err, ErrTransactionNotRolledBack)
		}

		return err
	}

	return u.CommitOrRollback()
}
