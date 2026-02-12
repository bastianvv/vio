import type { Episode } from "../../api/episodes";

type Props = {
  episode: Episode;
  onPlay?: (episodeId: number) => void;
  onDetails?: (episodeId: number) => void;
};

export default function EpisodeRow({ episode, onPlay, onDetails }: Props) {
  const handlePlay = () => {
    onPlay?.(episode.id);
  };

  const handleDetails = () => {
    onDetails?.(episode.id);
  };

  return (
    <div
      style={{
        display: "flex",
        justifyContent: "space-between",
        alignItems: "center",
        padding: "0.5rem 0",
        borderBottom: "1px solid #333",
        cursor: "pointer",
      }}
      onClick={handlePlay}
    >
      <div>
        <strong>Ep {episode.number}</strong> – {episode.title}
      </div>

      <div style={{ display: "flex", gap: "0.5rem" }}>
        <button
          onClick={(e) => {
            e.stopPropagation();
            handlePlay();
          }}
        >
          ▶
        </button>

        <button
          onClick={(e) => {
            e.stopPropagation();
            handleDetails();
          }}
        >
          Details
        </button>
      </div>
    </div>
  );
}
