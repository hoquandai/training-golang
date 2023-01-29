package main

import (
  "database/sql"
  "context"
  "time"
  "fmt"
  "log"
  _ "github.com/jackc/pgx/v5/stdlib"
)

type Bird struct {
  Species string
  Description string
}

func main() {
  db, err := sql.Open("pgx", "postgresql://postgres:postgres@localhost:5555/bird_encyclopedia")
  if err != nil {
    log.Fatalf("could not connect to database: %v", err)
  }

  if err := db.Ping(); err != nil {
    log.Fatalf("unable to reach database: %v", err)
  }

  fmt.Println("database is reachable")

  row := db.QueryRow("SELECT bird, description FROM birds LIMIT 1")
  bird := Bird{}

  if err := row.Scan(&bird.Species, &bird.Description); err != nil {
    log.Fatalf("could not scan row: %v", err)
  }

  fmt.Printf("found bird: %+v\n", bird)

  rows, err := db.Query("SELECT bird, description FROM birds")
  if err != nil {
    log.Fatalf("could not excute query: %v", err)
  }

  birds := []Bird{}

  for rows.Next() {
    bird := Bird{}

    if err := rows.Scan(&bird.Species, &bird.Description); err != nil {
      log.Fatalf("could not scan row: %v", err)
    }

    birds = append(birds, bird)
  }

  fmt.Printf("found %d birds: %+v\n", len(birds), birds)

  birdName := "eagle"
  birdRow := db.QueryRow("SELECT bird, description FROM birds WHERE bird = $1", birdName)

  if err := birdRow.Scan(&bird.Species, &bird.Description); err != nil {
    log.Fatalf("could not scan row: %v", err)
  }

  fmt.Printf("found bird with name %v, bird: %+v\n", birdName, bird)

  ctx := context.Background()
  cctx, _ := context.WithTimeout(ctx, 300 * time.Millisecond)
  _, err = db.QueryContext(cctx, "SELECT * FROM pg_sleep(1)")

  if err != nil {
    log.Fatalf("could not excute query: %v", err)
  }
}
