package config

import (
	"context"
	"idiom-api-services/domains/identities"
	"idiom-api-services/domains/org_members"
	"idiom-api-services/domains/organizations"
	"idiom-api-services/domains/projects"
	"idiom-api-services/domains/sessions"
	"idiom-api-services/packages/database/postgres"
	"idiom-api-services/packages/jwt"
	"log"
)

var (
	PG               *postgres.Postgres
	IdentityRepo     *identities.Repository
	SessionRepo      *sessions.Repository
	OrganizationRepo *organizations.Repository
	MemberRepo       *org_members.Repository
	ProjectRepo      *projects.Repository
	JWTSettings      *jwt.JWTSettings
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
	OrganizationRepo = organizations.NewRepository(PG)
	MemberRepo = org_members.NewRepository(PG)
	ProjectRepo = projects.NewRepository(PG)

	JWTSettings = jwt.NewJWTSettings(
		"",
		JWTPublicKey,
		0,
		TokenIssuer,
	)

	return nil
}

func Close() {
	if PG != nil {
		PG.Close()
	}
}
