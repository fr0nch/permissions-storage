//go:build uuid

package model

import "github.com/google/uuid"

// TODO: Add UUID support in the future
type UserID = uuid.UUID

const UserIDSQLColumnName = "user_id"
const UserIDSQLType = "UUID"
