import { Box } from "@mui/material";
import { useParams } from "react-router-dom";

export default function VideoPlayer() {
  const { fileId } = useParams<{ fileId: string }>();

  if (!fileId) return null;

  return (
    <Box
      sx={{
        width: "100%",
        height: "100%",
        bgcolor: "black",
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
      }}
    >
      <video
        src={`/api/files/${fileId}/stream`}
        controls
        autoPlay
        style={{
          width: "100%",
          maxHeight: "100vh",
          backgroundColor: "black",
        }}
      />
    </Box>
  );
}
