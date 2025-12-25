// src/api/seasons.ts
import { apiGet } from "./client";

export interface Season {
  id: number;
  series_id: number;
  number: number;
  has_poster: boolean;
}

export async function getSeasonsBySeries(seriesId: number): Promise<Season[]> {
  return apiGet<Season[]>(`/series/${seriesId}/seasons`);
}

export async function getSeason(id: number): Promise<Season> {
  return apiGet<Season>(`/seasons/${id}`);
}
