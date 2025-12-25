import { useEffect, useState } from "react";
import { useParams, useNavigate, useLocation } from "react-router-dom";
import { Box, Typography, Button } from "@mui/material";

import TopBar from "../../app/TopBar";
import { getMovieById } from "../../api/movies";
import type { Movie } from "../../api/movies";

export default function MovieDetails() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const location = useLocation();

  const { libraryId } = (location.state as { libraryId: number }) ?? {};

  const [movie, setMovie] = useState<Movie | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!id || !libraryId) return;

    (async () => {
      const data = await getMovieById(Number(id), libraryId);
      setMovie(data);
      setLoading(false);
    })();
  }, [id, libraryId]);

  if (!libraryId) {
    return <div>Missing library context</div>;
  }

  if (loading || !movie) {
    return <div>Loading…</div>;
  }

  return (
    <>
      <TopBar
        left={
          <Button color="inherit" onClick={() => navigate("/movies")}>
            ← Back
          </Button>
        }
        title={movie.title}
      />

      {/* Backdrop container */}
      <Box
        sx={{
          minHeight: "100vh",
          backgroundImage: movie.has_backdrop
            ? `url(/api/images/movies/${movie.id}/backdrop)`
            : undefined,
          backgroundSize: "cover",
          backgroundPosition: "center",
          position: "relative",
        }}
      >
        {/* Dark overlay */}
        <Box
          sx={{
            position: "absolute",
            inset: 0,
            background:
              "linear-gradient(to right, rgba(0,0,0,0.85) 0%, rgba(0,0,0,0.5) 40%, rgba(0,0,0,0.2) 100%)",
          }}
        />
        <Box
          sx={{
            minHeight: "100vh",
            background:
              "linear-gradient(to top, rgba(0,0,0,0.9), rgba(0,0,0,0.4))",
            padding: 4,
            display: "flex",
            alignItems: "flex-end",
          }}
        >
          <Box sx={{ maxWidth: 900 }}>
            <Box
              sx={{
                position: "relative",
                zIndex: 1,
                maxWidth: 1100,
                mx: "auto",
                px: 4,
                py: 6,
                display: "flex",
                gap: 3,
                color: "#fff",
              }}
            >
              {movie.has_poster && (
                <Box
                  component="img"
                  src={`/api/images/movies/${movie.id}/poster`}
                  alt={movie.title}
                  sx={{
                    width: 200,
                    borderRadius: 1,
                    boxShadow: 3,
                  }}
                />
              )}

              <Box>
                <Typography variant="h4" fontWeight="bold">
                  {movie.title}
                </Typography>

                {movie.original_title &&
                  movie.original_title !== movie.title && (
                    <Typography
                      variant="subtitle1"
                      sx={{ fontStyle: "italic", opacity: 0.85 }}
                    >
                      {movie.original_title}
                    </Typography>
                  )}

                <Typography sx={{ mt: 1, opacity: 0.8 }}>
                  {movie.year}
                </Typography>

                <Typography sx={{ mt: 2 }}>{movie.overview}</Typography>

                <Button variant="contained" sx={{ mt: 3 }}>
                  ▶ Play
                </Button>
              </Box>
            </Box>
          </Box>
        </Box>
      </Box>
    </>
  );
}
