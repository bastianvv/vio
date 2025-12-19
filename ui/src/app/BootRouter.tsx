import { useEffect, useState } from "react";
import { Navigate } from "react-router-dom";
import { CircularProgress, Box } from "@mui/material";

import { getLibraries } from "../api/libraries";
import type { Library } from "../api/libraries";

type BootState = { status: "loading" } | { status: "ready"; redirect: string };

export default function BootRouter() {
  const [state, setState] = useState<BootState>({ status: "loading" });

  useEffect(() => {
    getLibraries()
      .then((libraries: Library[]) => {
        if (libraries.length === 0) {
          setState({ status: "ready", redirect: "/setup" });
          return;
        }

        const hasMovies = libraries.some((l) => l.type === "movies");
        if (hasMovies) {
          setState({ status: "ready", redirect: "/movies" });
          return;
        }

        setState({ status: "ready", redirect: "/series" });
      })
      .catch(() => {
        // Fail safe: setup screen
        setState({ status: "ready", redirect: "/setup" });
      });
  }, []);

  if (state.status === "loading") {
    return (
      <Box
        sx={{
          height: "100%",
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
        }}
      >
        <CircularProgress />
      </Box>
    );
  }

  return <Navigate to={state.redirect} replace />;
}
