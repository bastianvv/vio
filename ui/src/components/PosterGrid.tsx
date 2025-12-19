import Grid from "@mui/material/GridLegacy";
import PosterCard from "./PosterCard";

interface Item {
  id: number;
  title: string;
  posterUrl?: string | null;
}

interface Props {
  items: Item[];
}

export default function PosterGrid({ items }: Props) {
  return (
    <Grid container spacing={2}>
      {items.map((item) => (
        <Grid item key={item.id} xs={6} sm={4} md={3} lg={2}>
          <PosterCard title={item.title} posterUrl={item.posterUrl} />
        </Grid>
      ))}
    </Grid>
  );
}
