import { useState } from "react";
import {
  Box,
  Button,
  Typography,
  Paper,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  MenuItem,
} from "@mui/material";
import { useNavigate } from "react-router-dom";

import { createLibrary } from "../../api/libraries";

export default function SetupPage() {
  const navigate = useNavigate();

  const [open, setOpen] = useState(false);
  const [name, setName] = useState("");
  const [path, setPath] = useState("");
  const [type, setType] = useState<"movies" | "series">("movies");
  const [loading, setLoading] = useState(false);

  const handleSubmit = async () => {
    setLoading(true);
    try {
      await createLibrary({ name, path, type });
      // Re-run boot logic
      navigate("/", { replace: true });
    } catch (err) {
      console.error(err);
      alert("Failed to create library");
    } finally {
      setLoading(false);
    }
  };

  return (
    <Box
      sx={{
        height: "100%",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
      }}
    >
      <Paper sx={{ p: 4, maxWidth: 400, textAlign: "center" }}>
        <Typography variant="h5" gutterBottom>
          Welcome to Vio
        </Typography>

        <Typography sx={{ mb: 3 }}>
          No libraries found. Add a library to get started.
        </Typography>

        <Button variant="contained" onClick={() => setOpen(true)}>
          Add Library
        </Button>
      </Paper>

      <Dialog open={open} onClose={() => setOpen(false)} fullWidth>
        <DialogTitle>Add Library</DialogTitle>

        <DialogContent
          sx={{ display: "flex", flexDirection: "column", gap: 2 }}
        >
          <TextField
            label="Name"
            value={name}
            onChange={(e) => setName(e.target.value)}
            fullWidth
          />

          <TextField
            label="Path"
            value={path}
            onChange={(e) => setPath(e.target.value)}
            fullWidth
            placeholder="/media/movies"
          />

          <TextField
            select
            label="Type"
            value={type}
            onChange={(e) => setType(e.target.value as "movies" | "series")}
          >
            <MenuItem value="movies">Movies</MenuItem>
            <MenuItem value="series">Series</MenuItem>
          </TextField>
        </DialogContent>

        <DialogActions>
          <Button onClick={() => setOpen(false)}>Cancel</Button>
          <Button
            variant="contained"
            onClick={handleSubmit}
            disabled={loading || !name || !path}
          >
            Add
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}
