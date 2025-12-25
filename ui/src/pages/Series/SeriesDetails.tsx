import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { Box, Button } from "@mui/material";

import type { Series } from "../../api/series";
import type { Season } from "../../api/seasons";
import type { Episode } from "../../api/episodes";
import { getSeriesById } from "../../api/series";
import { getSeasonsBySeries } from "../../api/seasons";
import { getEpisodesBySeason } from "../../api/episodes";

import SeasonBlock from "../../components/series/SeasonBlock";
import TopBar from "../../app/TopBar";

type SeasonWithEpisodes = {
  season: Season;
  episodes: Episode[];
};

export default function SeriesDetails() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const [series, setSeries] = useState<Series | null>(null);
  const [seasons, setSeasons] = useState<SeasonWithEpisodes[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const HERO_HEIGHT = 600; // px — tweakable

  useEffect(() => {
    if (!id) return;

    async function load() {
      try {
        setLoading(true);

        const seriesId = Number(id);
        const seriesData = await getSeriesById(seriesId);
        const seasonsData = await getSeasonsBySeries(seriesId);

        const seasonsWithEpisodes: SeasonWithEpisodes[] = await Promise.all(
          seasonsData.map(async (season) => ({
            season,
            episodes: await getEpisodesBySeason(season.id),
          })),
        );

        setSeries(seriesData);
        setSeasons(seasonsWithEpisodes);
      } catch (err) {
        console.error(err);
        setError("Failed to load series");
      } finally {
        setLoading(false);
      }
    }

    load();
  }, [id]);

  if (loading) {
    return <div>Loading…</div>;
  }

  if (error || !series) {
    return <div>{error ?? "Series not found"}</div>;
  }

  return (
    <>
      <TopBar
        left={
          <Button color="inherit" onClick={() => navigate("/series")}>
            ← Back
          </Button>
        }
        title={series.title}
      />

      {/* HERO */}
      <Box
        sx={{
          height: HERO_HEIGHT,
          backgroundImage: series.has_backdrop
            ? `url(/api/images/series/${series.id}/backdrop)`
            : undefined,
          backgroundSize: "cover",
          backgroundPosition: "center",
          position: "relative",
        }}
      >
        {/* Overlay */}
        <Box
          sx={{
            position: "absolute",
            inset: 0,
            background:
              "linear-gradient(to bottom, rgba(0,0,0,0.75), rgba(0,0,0,0.9))",
          }}
        />

        {/* Hero content (optional: title, metadata later) */}
        <Box
          sx={{
            position: "relative",
            zIndex: 1,
            maxWidth: 1100,
            mx: "auto",
            px: 3,
            pt: 6,
            color: "#fff",
          }}
        >
          {/* You can add series metadata here later */}
        </Box>
      </Box>

      {/* SCROLLING CONTENT */}
      <Box
        sx={{
          maxWidth: 1100,
          mx: "auto",
          px: 3,
          py: 4,
        }}
      >
        {seasons.map(({ season, episodes }) => (
          <SeasonBlock key={season.id} season={season} episodes={episodes} />
        ))}
      </Box>
    </>
  );
}
