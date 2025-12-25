import { Routes, Route } from "react-router-dom";

import BootRouter from "./BootRouter";
import SetupPage from "../pages/Setup/SetupPage";
import MoviesPage from "../pages/Movies/MoviesPage";
import MovieDetails from "../pages/Movies/MoviesDetails";
import SeriesPage from "../pages/Series/SeriesPage";
import SeriesDetails from "../pages/Series/SeriesDetails";
import VideoPlayer from "../components/player/VideoPlayer";

export default function AppRoutes() {
  return (
    <Routes>
      {/* Boot decision */}
      <Route path="/" element={<BootRouter />} />

      {/* Pages */}
      <Route path="/setup" element={<SetupPage />} />
      <Route path="/movies" element={<MoviesPage />} />
      <Route path="/movies/:id" element={<MovieDetails />} />
      <Route path="/series" element={<SeriesPage />} />
      <Route path="/series/:id" element={<SeriesDetails />} />
      <Route path="/player/:fileId" element={<VideoPlayer />} />

      {/* Fallback */}
      <Route path="*" element={<BootRouter />} />
    </Routes>
  );
}
