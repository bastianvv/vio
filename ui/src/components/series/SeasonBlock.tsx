import type { Season } from "../../api/seasons";
import type { Episode } from "../../api/episodes";
import EpisodeRow from "./EpisodeRow";
import { useNavigate } from "react-router-dom";
import { getEpisodeFiles } from "../../api/episodes"; // 👈 you already have this

type Props = {
  season: Season;
  episodes: Episode[];
};

export default function SeasonBlock({ season, episodes }: Props) {
  const navigate = useNavigate();

  async function handlePlayEpisode(episodeId: number) {
    try {
      const files = await getEpisodeFiles(episodeId);

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

  function handleEpisodeDetails(episodeId: number) {
    navigate(`/episodes/${episodeId}`, {
      state: { seriesId: season.series_id },
    });
  }

  if (episodes.length === 0) {
    return null;
  }

  return (
    <div style={{ marginBottom: "2rem" }}>
      <h3 style={{ marginBottom: "0.5rem" }}>Season {season.number}</h3>

      <div style={{ display: "flex", gap: "1rem" }}>
        {/* Poster */}
        <div style={{ width: 120, height: 180, background: "#222" }}>
          {season.has_poster && (
            <img
              src={`/api/images/seasons/${season.id}/poster`}
              alt={`Season ${season.number}`}
              style={{ width: "100%", height: "100%", objectFit: "cover" }}
            />
          )}
        </div>

        {/* Episodes */}
        <div style={{ flex: 1 }}>
          {episodes.map((episode) => (
            <EpisodeRow
              key={episode.id}
              episode={episode}
              onPlay={handlePlayEpisode}
              onDetails={handleEpisodeDetails}
            />
          ))}
        </div>
      </div>
    </div>
  );
}
