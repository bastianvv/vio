import { apiGet } from "./client";

export interface Movie {
  id: number;
  title: string;
  has_poster: boolean;
}

export function getMovies(libraryId: number): Promise<Movie[]> {
  return apiGet<Movie[]>(`/movies?library_id=${libraryId}`);
}
