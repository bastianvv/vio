import type { ReactNode } from "react";
import { AppBar, Toolbar, Typography, Button, Box } from "@mui/material";

import { useScanActions } from "../hooks/useScanActions";
import { triggerEnrich } from "../state/enrichActions";

type Props = {
  left?: ReactNode;
  title?: ReactNode;
};

export default function TopBar({ left, title }: Props) {
  const { fullScan, incrementalScan } = useScanActions();

  return (
    <AppBar
      position="fixed"
      sx={{ zIndex: (theme) => theme.zIndex.drawer + 1 }}
    >
      <Toolbar>
        {left && (
          <Box sx={{ mr: 2, display: "flex", alignItems: "center" }}>
            {left}
          </Box>
        )}

        <Typography variant="h6" sx={{ flexGrow: 1 }}>
          {title ?? "Vio"}
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
