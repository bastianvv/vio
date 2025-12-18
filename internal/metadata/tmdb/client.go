package tmdb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Client struct {
	apiKey  string
	baseURL string
	http    *http.Client
}

func New(apiKey string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: "https://api.themoviedb.org/3",
		http:    &http.Client{},
	}
}

func (c *Client) SearchMovie(
	ctx context.Context,
	title string,
	year int,
) ([]MovieResult, error) {

	q := url.Values{}
	q.Set("api_key", c.apiKey)
	q.Set("query", title)
	if year > 0 {
		q.Set("year", fmt.Sprint(year))
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		c.baseURL+"/search/movie?"+q.Encode(),
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tmdb search failed: %s", resp.Status)
	}

	var raw searchMovieResponse
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	out := make([]MovieResult, 0, len(raw.Results))
	for _, r := range raw.Results {
		out = append(out, MovieResult{
			ID:            fmt.Sprint(r.ID),
			Title:         r.Title,
			OriginalTitle: r.OriginalTitle,
			ReleaseYear:   parseYear(r.ReleaseDate),
		})
	}

	return out, nil
}

func (c *Client) GetMovie(
	ctx context.Context,
	tmdbID string,
) (*MovieDetails, error) {

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		fmt.Sprintf("%s/movie/%s?api_key=%s", c.baseURL, tmdbID, c.apiKey),
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tmdb get failed: %s", resp.Status)
	}

	var raw movieDetailsResponse
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	return &MovieDetails{
		ID:            fmt.Sprint(raw.ID),
		Title:         raw.Title,
		OriginalTitle: raw.OriginalTitle,
		Overview:      raw.Overview,
		RuntimeMin:    raw.Runtime,
		PosterPath:    raw.PosterPath,
		BackdropPath:  raw.BackdropPath,
		ReleaseYear:   parseYear(raw.ReleaseDate),
	}, nil
}
