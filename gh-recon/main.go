package ghrecon

import (
	"context"

	"github.com/charmbracelet/log"
	"github.com/google/go-github/v72/github"
)

type Recon struct {
	Client      *github.Client
	Logger      *log.Logger
	Ctx         context.Context
	Silent      bool
	JsonFile    string
	MaxRepoSize int
}

func NewRecon(
	client *github.Client,
	logger *log.Logger,
	ctx context.Context,
	silent bool,
	jsonFile string,
	maxRepoSize int,
) *Recon {
	return &Recon{
		Client:      client,
		Logger:      logger,
		Ctx:         ctx,
		Silent:      silent,
		JsonFile:    jsonFile,
		MaxRepoSize: maxRepoSize,
	}
}
