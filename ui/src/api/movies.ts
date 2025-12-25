import { apiGet } from "./client";

export interface Movie {
  library_id: number;
  id: number;
  title: string;
  original_title: string;
  year: number;
  overview: string;
  has_poster: boolean;
  has_backdrop: boolean;
}

export function getMovies(libraryId: number): Promise<Movie[]> {
  return apiGet<Movie[]>(`/movies?library_id=${libraryId}`);
}

export function getMovieById(
  movieId: number,
  libraryId: number,
): Promise<Movie> {
  return apiGet<Movie>(`/movies/${movieId}?library_id=${libraryId}`);
}
