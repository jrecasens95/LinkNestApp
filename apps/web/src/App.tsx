import { Link, Route, Routes } from "react-router-dom";
import { DashboardPage } from "./pages/DashboardPage";
import { HomePage } from "./pages/HomePage";
import { LinkDetailPage } from "./pages/LinkDetailPage";

export function App() {
  return (
    <>
      <nav className="top-nav">
        <Link to="/">LinkNest</Link>
        <div>
          <Link to="/">Create</Link>
          <Link to="/dashboard">Dashboard</Link>
        </div>
      </nav>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/dashboard" element={<DashboardPage />} />
        <Route path="/dashboard/links/:id" element={<LinkDetailPage />} />
      </Routes>
    </>
  );
}
