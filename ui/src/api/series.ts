import { apiGet } from "./client";

export interface Series {
  id: number;
  title: string;
  has_poster: boolean;
  has_backdrop: boolean;
}

export function getSeries(): Promise<Series[]> {
  return apiGet<Series[]>(`/series`);
}

export function getSeriesById(id: number): Promise<Series> {
  return apiGet<Series>(`/series/${id}`);
}
