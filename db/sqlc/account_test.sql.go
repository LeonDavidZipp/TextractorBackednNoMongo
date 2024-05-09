package db

import (
	"context"
	"testing"
	"time"
	"database/sql"
	"github.com/stretchr/testify/require"
	"github.com/LeonDavidZipp/Textractor"
	"net/mail"
)


func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner : util.RandomName(),
		Email : util.RandomEmail(),
		GoogleID : nil,
		FacebookID : nil,
	}

}
