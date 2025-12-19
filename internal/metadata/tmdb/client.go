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

func (c *Client) SearchTV(
	ctx context.Context,
	name string,
) ([]TVResult, error) {

	q := url.Values{}
	q.Set("api_key", c.apiKey)
	q.Set("query", name)

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		c.baseURL+"/search/tv?"+q.Encode(),
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
		return nil, fmt.Errorf("tmdb search tv failed: %s", resp.Status)
	}

	var raw searchTVResponse
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	out := make([]TVResult, 0, len(raw.Results))
	for _, r := range raw.Results {
		out = append(out, TVResult{
			ID:           fmt.Sprint(r.ID),
			Name:         r.Name,
			OriginalName: r.OriginalName,
			FirstAirYear: parseYear(r.FirstAirDate),
		})
	}

	return out, nil
}

func (c *Client) GetTV(
	ctx context.Context,
	tmdbID string,
) (*TVDetails, error) {

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		fmt.Sprintf("%s/tv/%s?api_key=%s", c.baseURL, tmdbID, c.apiKey),
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
		return nil, fmt.Errorf("tmdb get tv failed: %s", resp.Status)
	}

	var raw tvDetailsResponse
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	var seasons []TVSeason
	for _, s := range raw.Seasons {
		if s.SeasonNumber <= 0 {
			continue // skip specials
		}
		seasons = append(seasons, TVSeason{
			ID:          fmt.Sprint(s.ID),
			Number:      s.SeasonNumber,
			Title:       s.Name,
			Overview:    s.Overview,
			PosterPath:  s.PosterPath,
			AirDateYear: parseYear(s.AirDate),
		})
	}

	return &TVDetails{
		ID:           fmt.Sprint(raw.ID),
		Name:         raw.Name,
		OriginalName: raw.OriginalName,
		Overview:     raw.Overview,
		Status:       raw.Status,
		PosterPath:   raw.PosterPath,
		BackdropPath: raw.BackdropPath,
		Seasons:      seasons,
	}, nil
}

func (c *Client) GetSeason(
	ctx context.Context,
	tvID string,
	seasonNumber int,
) ([]TVEpisode, error) {

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		fmt.Sprintf(
			"%s/tv/%s/season/%d?api_key=%s",
			c.baseURL,
			tvID,
			seasonNumber,
			c.apiKey,
		),
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
		return nil, fmt.Errorf("tmdb get season failed: %s", resp.Status)
	}

	var raw seasonDetailsResponse
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	var eps []TVEpisode
	for _, e := range raw.Episodes {
		eps = append(eps, TVEpisode{
			ID:         fmt.Sprint(e.ID),
			Number:     e.EpisodeNum,
			Title:      e.Name,
			Overview:   e.Overview,
			RuntimeMin: e.Runtime,
			AirYear:    parseYear(e.AirDate),
			StillPath:  e.StillPath,
		})
	}

	return eps, nil
}
