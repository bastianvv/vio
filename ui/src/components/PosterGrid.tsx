import Grid from "@mui/material/GridLegacy";
import PosterCard from "./PosterCard";

interface Item {
  id: number;
  title: string;
  posterUrl?: string | null;
}

interface Props {
  items: Item[];
  onItemClick?: (id: number) => void;
}

export default function PosterGrid({ items, onItemClick }: Props) {
  return (
    <Grid container spacing={2}>
      {items.map((item) => (
        <Grid
          item
          key={item.id}
          xs={6}
          sm={4}
          md={3}
          lg={2}
          onClick={() => onItemClick?.(item.id)}
          sx={{ cursor: onItemClick ? "pointer" : "default" }}
        >
          <PosterCard title={item.title} posterUrl={item.posterUrl} />
        </Grid>
      ))}
    </Grid>
  );
}
