import { useEffect, useState } from "react";
import { useLocation, useNavigate, useParams } from "react-router-dom";
import { Box, Button, Typography } from "@mui/material";

import TopBar from "../../app/TopBar";
import { getEpisode, getEpisodeFiles, type Episode } from "../../api/episodes";
import { getSeason, type Season } from "../../api/seasons";
import { getSeriesById, type Series } from "../../api/series";

type LocationState = {
  seriesId?: number;
};

export default function EpisodeDetails() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const location = useLocation();

  const [episode, setEpisode] = useState<Episode | null>(null);
  const [season, setSeason] = useState<Season | null>(null);
  const [series, setSeries] = useState<Series | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const state = (location.state as LocationState | null) ?? null;

  useEffect(() => {
    if (!id) return;

    async function load() {
      try {
        setLoading(true);
        setError(null);

        const episodeId = Number(id);
        const episodeData = await getEpisode(episodeId);
        const seasonData = await getSeason(episodeData.season_id);
        const seriesData = await getSeriesById(seasonData.series_id);

        setEpisode(episodeData);
        setSeason(seasonData);
        setSeries(seriesData);
      } catch (err) {
        console.error(err);
        setError("Failed to load episode");
      } finally {
        setLoading(false);
      }
    }

    load();
  }, [id]);

  if (loading) {
    return <div>Loading…</div>;
  }

  if (error || !episode || !season || !series) {
    return <div>{error ?? "Episode not found"}</div>;
  }

  const currentEpisode = episode;
  const currentSeason = season;
  const currentSeries = series;

  const imageUrl = currentEpisode.has_still
    ? `/api/images/episodes/${currentEpisode.id}/still`
    : undefined;

  const backdropUrl = currentSeries.has_backdrop
    ? `/api/images/series/${currentSeries.id}/backdrop`
    : imageUrl;

  async function handlePlayEpisode() {
    try {
      const files = await getEpisodeFiles(currentEpisode.id);
      const file = files.find((f) => !f.is_missing);

      if (!file) {
        alert("No playable file found for this episode");
        return;
      }

      navigate(`/player/${file.id}`);
    } catch (err) {
      console.error(err);
      alert("Failed to load episode media");
    }
  }

  function handleBack() {
    const targetSeriesId = state?.seriesId ?? currentSeries.id;
    navigate(`/series/${targetSeriesId}`);
  }

  return (
    <>
      <TopBar
        left={
          <Button color="inherit" onClick={handleBack}>
            ← Back
          </Button>
        }
        title={`${currentSeries.title} • S${currentSeason.number}E${currentEpisode.number}`}
      />

      <Box
        sx={{
          minHeight: "100vh",
          backgroundImage: backdropUrl ? `url(${backdropUrl})` : undefined,
          backgroundSize: "cover",
          backgroundPosition: "center",
          position: "relative",
        }}
      >
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
            {imageUrl && (
              <Box
                component="img"
                src={imageUrl}
                alt={episode.title}
                sx={{
                  width: 320,
                  aspectRatio: "16 / 9",
                  objectFit: "cover",
                  borderRadius: 1,
                  boxShadow: 3,
                  alignSelf: "flex-start",
                }}
              />
            )}

            <Box sx={{ maxWidth: 700 }}>
              <Typography variant="h4" fontWeight="bold">
                {currentEpisode.title || `Episode ${currentEpisode.number}`}
              </Typography>

              <Typography sx={{ mt: 1, opacity: 0.85 }}>
                {currentSeries.title} • Season {currentSeason.number} • Episode {currentEpisode.number}
              </Typography>

              <Typography sx={{ mt: 2 }}>
                {currentEpisode.overview || "No summary available."}
              </Typography>

              <Button variant="contained" sx={{ mt: 3 }} onClick={handlePlayEpisode}>
                ▶ Play
              </Button>
            </Box>
          </Box>
        </Box>
      </Box>
    </>
  );
}
