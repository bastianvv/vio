import {
  Card,
  CardActionArea,
  CardContent,
  CardMedia,
  Typography,
  Box,
} from "@mui/material";
import { useState } from "react";

interface Props {
  title: string;
  posterUrl?: string | null;
}

export default function PosterCard({ title, posterUrl }: Props) {
  const [error, setError] = useState(false);

  const showImage = posterUrl && !error;

  return (
    <Card>
      <CardActionArea>
        {showImage ? (
          <CardMedia
            component="img"
            height="240"
            image={posterUrl}
            alt={title}
            loading="lazy"
            onError={() => setError(true)}
          />
        ) : (
          <Box
            sx={{
              height: 240,
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              bgcolor: "grey.800",
              color: "grey.300",
              fontSize: 14,
            }}
          >
            No Poster
          </Box>
        )}

        <CardContent>
          <Typography variant="subtitle2" noWrap title={title}>
            {title}
          </Typography>
        </CardContent>
      </CardActionArea>
    </Card>
  );
}
