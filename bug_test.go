package bug

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"testing"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"

	"entgo.io/bug/ent"
	"entgo.io/bug/ent/enttest"
	"entgo.io/bug/ent/user"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func TestBugPostrgres(t *testing.T) {
	for version, port := range map[string]int{"14": 15434} {
		t.Run(version, func(t *testing.T) {
			client := enttest.Open(t, dialect.Postgres, fmt.Sprintf("host=localhost port=%d user=postgres dbname=test password=pass sslmode=disable", port))
			defer client.Close()
			test(t, client)
		})
	}
}

func test(t *testing.T, client *ent.Client) {
	ctx := context.Background()

	idBytes := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, idBytes); err != nil {
		t.Fatal(err)
	}

	id := hex.EncodeToString(idBytes)
	client.User.Create().
		SetID(id).
		ExecX(ctx)

	tx, err := client.Tx(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if err := tx.User.Create().
		SetID(id).
		OnConflict(
			sql.ConflictColumns(user.FieldID),
			sql.ResolveWithIgnore(),
		).Exec(ctx); err != nil {
		t.Fatal(err)
	}

	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
}
