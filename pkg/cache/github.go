package cache

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/deparr/api/pkg/model"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

var ghUpdateTime = time.Hour * 2

type (
	RepoFrag struct {
		Owner struct {
			Login string
		}
		Name           string
		Url            string
		Description    string
		StargazerCount int
		Languages      struct {
			Edges []struct {
				Node struct {
					Name  string
					Color string
				}
				Size int
			}
		} `graphql:"languages(first:5,orderBy:{direction:DESC,field:SIZE})"`
		PushedAt string
	}

	pinsQuery struct {
		User struct {
			PinnedItems struct {
				TotalCount int
				Nodes      []struct {
					RepoFrag `graphql:"... on Repository"`
				}
			} `graphql:"pinnedItems(first: 6, types: [REPOSITORY])"`
		} `graphql:"user(login: \"deparr\")"`
	}

	recentsQuery struct {
		User struct {
			Repositories struct {
				Nodes []RepoFrag
			} `graphql:"repositories(first:6,orderBy:{field:PUSHED_AT,direction:DESC},visibility:PUBLIC)"`
		} `graphql:"user(login: \"deparr\")"`
	}
)

func makeGhClient(ctx context.Context) *githubv4.Client {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_ACCESS_TOKEN")},
	)
	httpClient := oauth2.NewClient(ctx, src)
	client := githubv4.NewClient(httpClient)

	return client
}

func fetchPinned(ctx context.Context, client *githubv4.Client) (*pinsQuery, error) {
	q := pinsQuery{}
	err := client.Query(ctx, &q, nil)
	if err != nil {
		return nil, err
	}
	return &q, nil
}

func fetchRecent(ctx context.Context, client *githubv4.Client) (*recentsQuery, error) {
	q := recentsQuery{}
	err := client.Query(ctx, &q, nil)
	if err != nil {
		return nil, err
	}
	return &q, nil
}

func cleanRepoQueryRes(repos []RepoFrag) []model.Repository {
	clean := make([]model.Repository, len(repos))

	for i, repo := range repos {

		var cleaned model.Repository
		cleaned.Owner = repo.Owner.Login
		cleaned.Name = repo.Name
		cleaned.Url = repo.Url
		cleaned.Desc = repo.Description
		cleaned.Stars = repo.StargazerCount
		cleaned.PushedAt = repo.PushedAt

		langs := make(model.RepoLangs, len(repo.Languages.Edges))
		totalSize := 0
		for j, rlang := range repo.Languages.Edges {
			langs[j].Percent = rlang.Size
			langs[j].Name = rlang.Node.Name
			langs[j].Color = rlang.Node.Color
			totalSize += rlang.Size
		}

		for j := range langs {
			langs[j].Percent = langs[j].Percent * 100 / totalSize
		}

		cleaned.Language = langs

		clean[i] = cleaned
	}
	return clean
}

func updateRepoCache(t time.Time) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*80)
	client := makeGhClient(ctx)

	pinnedRes, err := fetchPinned(ctx, client)
	var pinnedClean []model.Repository
	if err != nil {
		slog.Error("fetching cache(gh.pinned", "err", err)
	} else {
		nodes := pinnedRes.User.PinnedItems.Nodes
		frags := make([]RepoFrag, len(nodes))
		for i := range nodes {
			frags[i] = nodes[i].RepoFrag
		}
		pinnedClean = cleanRepoQueryRes(frags)
	}

	recentRes, err := fetchRecent(ctx, client)
	var recentClean []model.Repository
	if err != nil {
		slog.Error("fetching cache(gh.recent", "err", err)
	} else {
		recentClean = cleanRepoQueryRes(recentRes.User.Repositories.Nodes)
	}

	globalCache.setGithub("pinned", pinnedClean, t)
	globalCache.setGithub("recent", recentClean, t)

	cancel()
}

func initRepoCache() {
	cacheTickers.github = time.NewTicker(ghUpdateTime)
	toggleChans.github = make(chan bool, 1)
	slog.Info("initing cache(gh)")
	updateRepoCache(time.Now())

	go func() {
		updateTicker := cacheTickers.github
		toggleChan := toggleChans.github
		for {
			select {
			case reset := <-toggleChan:
				slog.Info("toggling cache(gh)", "on", reset)
				if reset {
					updateRepoCache(time.Now())
					updateTicker.Reset(ghUpdateTime)
				} else {
					updateTicker.Stop()
					// <-updateTicker.C
				}
			case t := <-updateTicker.C:
				slog.Info("refetching cache(gh)")
				updateRepoCache(t)
			}
		}
	}()

}
