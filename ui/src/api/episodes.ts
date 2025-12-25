// src/api/episodes.ts
import { apiGet } from "./client";

export interface Episode {
  id: number;
  season_id: number;
  number: number;
  title: string;
  overview?: string;
  runtime_min?: number;
  has_still?: boolean;
}

export async function getEpisodesBySeason(
  seasonId: number,
): Promise<Episode[]> {
  return apiGet<Episode[]>(`/seasons/${seasonId}/episodes`);
}

export async function getEpisode(id: number): Promise<Episode> {
  return apiGet<Episode>(`/episodes/${id}`);
}
