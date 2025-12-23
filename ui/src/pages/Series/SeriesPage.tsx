import { useEffect, useState } from "react";
import { CircularProgress, Typography, Box } from "@mui/material";

import { getLibraries } from "../../api/libraries";
import { getSeries } from "../../api/series";
import type { Series } from "../../api/series";
import type { Library } from "../../api/libraries";

import { onScanFinished } from "../../state/scanEvents";
import PosterGrid from "../../components/PosterGrid";

export default function SeriesPage() {
  const [series, setSeries] = useState<Series[] | null>(null);
  const [library, setLibrary] = useState<Library | null>(null);

  async function loadSeries() {
    try {
      const libraries = await getLibraries();
      const seriesLib = libraries.find((l) => l.type === "series");

      if (!seriesLib) {
        setLibrary(null);
        setSeries([]);
        return;
      }

      setLibrary(seriesLib);
      const series = await getSeries();
      setSeries(series);
    } catch (err) {
      console.error(err);
      setSeries([]);
    }
  }

  // Initial load (mount only)
  useEffect(() => {
    let cancelled = false;

    (async () => {
      if (!cancelled) {
        await loadSeries();
      }
    })();

    return () => {
      cancelled = true;
    };
  }, []);

  // Refresh after scan
  useEffect(() => {
    return onScanFinished(() => {
      setSeries(null); // show spinner
      loadSeries();
    });
  }, []);

  if (series === null) {
    return (
      <Box sx={{ display: "flex", justifyContent: "center", mt: 4 }}>
        <CircularProgress />
      </Box>
    );
  }

  if (!library) {
    return (
      <>
        <Typography variant="h4">Series</Typography>
        <Typography sx={{ mt: 2 }}>No series libraries configured.</Typography>
      </>
    );
  }

  if (series.length === 0) {
    return (
      <>
        <Typography variant="h4">Series</Typography>
        <Typography sx={{ mt: 2 }}>
          No series found. Run a scan to populate this library.
        </Typography>
      </>
    );
  }

  return (
    <>
      <Typography variant="h4" sx={{ mb: 2 }}>
        Series
      </Typography>

      <PosterGrid
        items={series.map((s) => ({
          id: s.id,
          title: s.title,
          posterUrl: s.has_poster ? `/api/images/series/${s.id}/poster` : null,
        }))}
      />
    </>
  );
}
