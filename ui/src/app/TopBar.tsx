import { AppBar, Toolbar, Typography, Button, Box } from "@mui/material";

import { useScanActions } from "../hooks/useScanActions";
import { triggerEnrich, hasEnrichHandler } from "../state/enrichActions";

export default function TopBar() {
  const { fullScan, incrementalScan } = useScanActions();

  return (
    <AppBar position="fixed">
      <Toolbar>
        <Typography variant="h6" sx={{ flexGrow: 1 }}>
          Vio
        </Typography>

        <Box sx={{ display: "flex", gap: 1 }}>
          <Button color="inherit" onClick={fullScan}>
            Full Scan
          </Button>
          <Button color="inherit" onClick={incrementalScan}>
            Incr. Scan
          </Button>
          <Button color="inherit" onClick={triggerEnrich}>
            Enrich
          </Button>
          <Button color="inherit">⚙</Button>
        </Box>
      </Toolbar>
    </AppBar>
  );
}
