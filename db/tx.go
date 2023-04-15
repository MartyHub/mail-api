package db

import "github.com/jackc/pgx/v5"

func TxReadOnly() pgx.TxOptions {
	return pgx.TxOptions{
		IsoLevel:       pgx.RepeatableRead,
		AccessMode:     pgx.ReadOnly,
		DeferrableMode: pgx.NotDeferrable,
	}
}

func TxWrite() pgx.TxOptions {
	return pgx.TxOptions{
		IsoLevel:       pgx.RepeatableRead,
		AccessMode:     pgx.ReadWrite,
		DeferrableMode: pgx.NotDeferrable,
	}
}
