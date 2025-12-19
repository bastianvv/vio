import { apiGet } from "./client";

export interface Library {
  id: number;
  name: string;
  type: "movies" | "series";
}

export function getLibraries(): Promise<Library[]> {
  return apiGet<Library[]>("/libraries");
}

export interface CreateLibraryRequest {
  name: string;
  path: string;
  type: "movies" | "series";
}

export async function createLibrary(
  req: CreateLibraryRequest,
): Promise<Library> {
  const res = await fetch("/api/libraries", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(req),
  });

  if (!res.ok) {
    throw new Error(`Failed to create library: ${res.status}`);
  }

  return res.json();
}
