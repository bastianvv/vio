import type { Episode } from "../../api/episodes";

type Props = {
  episode: Episode;
  onPlay?: (episodeId: number) => void;
};

export default function EpisodeRow({ episode, onPlay }: Props) {
  const handlePlay = () => {
    onPlay?.(episode.id);
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

      <button
        onClick={(e) => {
          e.stopPropagation();
          handlePlay();
        }}
      >
        ▶
      </button>
    </div>
  );
}
