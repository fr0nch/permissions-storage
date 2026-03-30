//go:build !uuid

package model

type UserID = uint64

const UserIDSQLColumnName = "steamid64"
const UserIDSQLType = "BIGINT"
