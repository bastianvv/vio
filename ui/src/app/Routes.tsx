import { Routes, Route } from "react-router-dom";

import BootRouter from "./BootRouter";
import SetupPage from "../pages/Setup/SetupPage";
import MoviesPage from "../pages/Movies/MoviesPage";
import SeriesPage from "../pages/Series/SeriesPage";

export default function AppRoutes() {
  return (
    <Routes>
      {/* Boot decision */}
      <Route path="/" element={<BootRouter />} />

      {/* Pages */}
      <Route path="/setup" element={<SetupPage />} />
      <Route path="/movies" element={<MoviesPage />} />
      <Route path="/series" element={<SeriesPage />} />

      {/* Fallback */}
      <Route path="*" element={<BootRouter />} />
    </Routes>
  );
}
