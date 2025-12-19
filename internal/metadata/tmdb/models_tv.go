package tmdb

type searchTVResponse struct {
	Results []struct {
		ID           int    `json:"id"`
		Name         string `json:"name"`
		OriginalName string `json:"original_name"`
		FirstAirDate string `json:"first_air_date"`
	} `json:"results"`
}

type tvDetailsResponse struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	OriginalName string  `json:"original_name"`
	Overview     string  `json:"overview"`
	Status       string  `json:"status"`
	PosterPath   *string `json:"poster_path"`
	BackdropPath *string `json:"backdrop_path"`
	Seasons      []struct {
		ID           int     `json:"id"`
		SeasonNumber int     `json:"season_number"`
		Name         string  `json:"name"`
		Overview     string  `json:"overview"`
		PosterPath   *string `json:"poster_path"`
		AirDate      string  `json:"air_date"`
	} `json:"seasons"`
}

type seasonDetailsResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Overview string `json:"overview"`
	Episodes []struct {
		ID         int     `json:"id"`
		EpisodeNum int     `json:"episode_number"`
		Name       string  `json:"name"`
		Overview   string  `json:"overview"`
		Runtime    int     `json:"runtime"`
		AirDate    string  `json:"air_date"`
		StillPath  *string `json:"still_path"`
	} `json:"episodes"`
}

type TVResult struct {
	ID           string
	Name         string
	OriginalName string
	FirstAirYear int
}

type TVDetails struct {
	ID           string
	Name         string
	OriginalName string
	Overview     string
	Status       string
	PosterPath   *string
	BackdropPath *string
	Seasons      []TVSeason
}

type TVSeason struct {
	ID          string
	Number      int
	Title       string
	Overview    string
	PosterPath  *string
	AirDateYear int
}

type TVEpisode struct {
	ID         string
	Number     int
	Title      string
	Overview   string
	RuntimeMin int
	AirYear    int
	StillPath  *string
}
