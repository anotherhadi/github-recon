package ghrecon

import (
	"context"

	"github.com/charmbracelet/log"
	"github.com/google/go-github/v72/github"
)

type Recon struct {
	client *github.Client
	logger *log.Logger
	ctx    context.Context
}

func NewRecon(client *github.Client, logger *log.Logger, ctx context.Context) *Recon {
	return &Recon{
		client: client,
		logger: logger,
		ctx:    ctx,
	}
}
