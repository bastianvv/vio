import { apiPost } from "./client";

export function enrichMovie(movieId: number) {
  return apiPost(`/movies/${movieId}/enrich`);
}

export function enrichSeries(seriesId: number) {
  return apiPost(`/series/${seriesId}/enrich`);
}
