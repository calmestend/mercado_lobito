package main

import (
	"github.com/calmestend/mercado_lobito/internal/auth"
	"github.com/calmestend/mercado_lobito/internal/db"
	"github.com/calmestend/mercado_lobito/internal/router"
	"github.com/calmestend/mercado_lobito/pkg/env"
)

func main() {
	env.Init()
	dbConn := db.Init()
	auth.SetDBConnection(dbConn)
	router.Init()
}
