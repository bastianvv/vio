import {
  Card,
  CardActionArea,
  CardContent,
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
    <Card
      sx={{
        transition: "transform 0.15s ease, box-shadow 0.15s ease",
        "&:hover": {
          transform: "scale(1.03)",
          boxShadow: 6,
        },
      }}
    >
      <CardActionArea>
        {/* Poster area */}
        <Box
          sx={{
            position: "relative",
            width: "100%",
            aspectRatio: "2 / 3", // ⭐ poster ratio
            bgcolor: "grey.900",
          }}
        >
          {showImage ? (
            <Box
              component="img"
              src={posterUrl}
              alt={title}
              loading="lazy"
              onError={() => setError(true)}
              sx={{
                position: "absolute",
                inset: 0,
                width: "100%",
                height: "100%",
                objectFit: "cover", // ⭐ key line
              }}
            />
          ) : (
            <Box
              sx={{
                position: "absolute",
                inset: 0,
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
                color: "grey.400",
                fontSize: 14,
              }}
            >
              No Poster
            </Box>
          )}
        </Box>

        <CardContent sx={{ p: 1 }}>
          <Typography variant="subtitle2" noWrap title={title}>
            {title}
          </Typography>
        </CardContent>
      </CardActionArea>
    </Card>
  );
}
