import type { Season } from "../../api/seasons";
import type { Episode } from "../../api/episodes";
import EpisodeRow from "./EpisodeRow";

type Props = {
  season: Season;
  episodes: Episode[];
};

export default function SeasonBlock({ season, episodes }: Props) {
  if (episodes.length === 0) {
    return null;
  }

  return (
    <div style={{ marginBottom: "2rem" }}>
      <h3 style={{ marginBottom: "0.5rem" }}>Season {season.number}</h3>

      <div
        style={{
          display: "flex",
          gap: "1rem",
          alignItems: "flex-start",
        }}
      >
        {/* Season poster */}
        <div
          style={{
            width: "120px",
            height: "180px",
            backgroundColor: "#222",
            flexShrink: 0,
          }}
        >
          {season.has_poster ? (
            <img
              src={`/api/images/seasons/${season.id}/poster`}
              alt={`Season ${season.number}`}
              style={{ width: "100%", height: "100%", objectFit: "cover" }}
            />
          ) : null}
        </div>

        {/* Episodes */}
        <div style={{ flex: 1 }}>
          {episodes.map((episode) => (
            <EpisodeRow key={episode.number} episode={episode} />
          ))}
        </div>
      </div>
    </div>
  );
}
