export async function scanLibrary(
  libraryId: number,
  mode: "scan" | "rescan",
): Promise<void> {
  const res = await fetch(`/api/libraries/${libraryId}/${mode}`, {
    method: "POST",
  });

  if (!res.ok) {
    throw new Error(`Failed to ${mode} library ${libraryId}`);
  }
}
