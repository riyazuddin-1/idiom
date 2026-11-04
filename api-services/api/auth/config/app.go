package config

import (
	"context"
	"idiom-api-services/domains/identities"
	"idiom-api-services/domains/projects"
	"idiom-api-services/domains/sessions"
	"idiom-api-services/packages/database/postgres"
	"idiom-api-services/packages/email"
	"idiom-api-services/packages/jwt"
	"log"
)

var (
	PG           *postgres.Postgres
	IdentityRepo *identities.Repository
	SessionRepo  *sessions.Repository
	ProjectRepo  *projects.Repository
	JWTSettings  *jwt.JWTSettings
	EmailSender  email.Sender
)

func Initialize(ctx context.Context) error {
	var err error
	pg, err := postgres.Init(ctx, PostgresDSN)
	if err != nil {
		log.Fatalf("failed to initialize postgres: %v", err)
		return err
	}

	PG = pg
	IdentityRepo = identities.NewRepository(PG)
	SessionRepo = sessions.NewRepository(PG)
	ProjectRepo = projects.NewRepository(PG)

	JWTSettings = jwt.NewJWTSettings(
		JWTPrivateKey,
		JWTPublicKey,
		TokenExpirationInSeconds,
		TokenIssuer,
	)

	EmailSender = email.NewHTTPMailSender()

	return nil
}

func Close() {
	if PG != nil {
		PG.Close()
	}
}
