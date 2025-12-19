import { getLibraries } from "../api/libraries";
import { scanLibrary } from "../api/scans";
import { notifyScanFinished } from "../state/scanEvents";

export function useScanActions() {
  const runScan = async (mode: "scan" | "rescan") => {
    const libraries = await getLibraries();

    for (const lib of libraries) {
      await scanLibrary(lib.id, mode);
    }

    // Notify UI that a scan completed
    notifyScanFinished();
  };

  return {
    fullScan: () => runScan("rescan"),
    incrementalScan: () => runScan("scan"),
  };
}
