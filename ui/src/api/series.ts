import { apiGet } from "./client";

export interface Series {
  id: number;
  title: string;
  has_poster: boolean;
}

export function getSeries(): Promise<Series[]> {
  return apiGet<Series[]>(`/series`);
}
