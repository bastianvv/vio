import { useEffect, useState } from "react";
import { CircularProgress, Typography, Box } from "@mui/material";
import { useNavigate } from "react-router-dom";
import { getLibraries } from "../../api/libraries";
import { getMovies } from "../../api/movies";
import type { Movie } from "../../api/movies";
import type { Library } from "../../api/libraries";
import { enrichMovie } from "../../api/enrich";
import { registerEnrichHandler } from "../../state/enrichActions";
import { onScanFinished } from "../../state/scanEvents";
import PosterGrid from "../../components/PosterGrid";

export default function MoviesPage() {
  const navigate = useNavigate();
  const [movies, setMovies] = useState<Movie[] | null>(null);
  const [library, setLibrary] = useState<Library | null>(null);

  async function enrichAllMovies() {
    if (!movies) return;
    await Promise.all(movies.map((m) => enrichMovie(m.id)));
  }

  async function loadMovies() {
    try {
      const libraries = await getLibraries();
      const movieLib = libraries.find((l) => l.type === "movies");

      if (!movieLib) {
        setLibrary(null);
        setMovies([]);
        return;
      }

      setLibrary(movieLib);
      const movies = await getMovies(movieLib.id);
      setMovies(movies);
    } catch (err) {
      console.error(err);
      setMovies([]);
    }
  }

  useEffect(() => {
    if (movies && movies.length > 0) {
      registerEnrichHandler(enrichAllMovies);
    } else {
      registerEnrichHandler(null);
    }

    return () => registerEnrichHandler(null);
  }, [movies]);

  // Initial load (mount only)
  useEffect(() => {
    let cancelled = false;

    (async () => {
      if (!cancelled) {
        await loadMovies();
      }
    })();

    return () => {
      cancelled = true;
    };
  }, []);

  // Refresh after scan
  useEffect(() => {
    return onScanFinished(() => {
      setMovies(null); // show spinner
      loadMovies();
    });
  }, []);

  if (movies === null) {
    return (
      <Box sx={{ display: "flex", justifyContent: "center", mt: 4 }}>
        <CircularProgress />
      </Box>
    );
  }

  if (!library) {
    return (
      <>
        <Typography variant="h4">Movies</Typography>
        <Typography sx={{ mt: 2 }}>No movie libraries configured.</Typography>
      </>
    );
  }

  if (movies.length === 0) {
    return (
      <>
        <Typography variant="h4">Movies</Typography>
        <Typography sx={{ mt: 2 }}>
          No movies found. Run a scan to populate this library.
        </Typography>
      </>
    );
  }

  return (
    <>
      <Typography variant="h4" sx={{ mb: 2 }}>
        Movies
      </Typography>
      <Box sx={{ px: 3, py: 2 }}>
        <PosterGrid
          items={movies.map((m) => ({
            id: m.id,
            title: m.title,
            posterUrl: m.has_poster
              ? `/api/images/movies/${m.id}/poster`
              : null,
          }))}
          onItemClick={(id) =>
            navigate(`/movies/${id}`, {
              state: { libraryId: library.id },
            })
          }
        />
      </Box>
    </>
  );
}
