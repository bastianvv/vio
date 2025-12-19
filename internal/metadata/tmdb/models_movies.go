package tmdb

type searchMovieResponse struct {
	Results []struct {
		ID            int    `json:"id"`
		Title         string `json:"title"`
		OriginalTitle string `json:"original_title"`
		ReleaseDate   string `json:"release_date"`
	} `json:"results"`
}

type movieDetailsResponse struct {
	ID            int     `json:"id"`
	Title         string  `json:"title"`
	OriginalTitle string  `json:"original_title"`
	Overview      string  `json:"overview"`
	Runtime       int     `json:"runtime"`
	PosterPath    *string `json:"poster_path"`
	BackdropPath  *string `json:"backdrop_path"`
	ReleaseDate   string  `json:"release_date"`
}

type MovieResult struct {
	ID            string
	Title         string
	OriginalTitle string
	ReleaseYear   int
}

type MovieDetails struct {
	ID            string
	Title         string
	OriginalTitle string
	Overview      string
	RuntimeMin    int
	PosterPath    *string
	BackdropPath  *string
	ReleaseYear   int
}
