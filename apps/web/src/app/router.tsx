import { Route, Routes } from "react-router-dom";
import { HomePage } from "../features/home";

function AboutPage() {
  return (
    <section className="card-grid">
      <article className="card">
        <h2>About this starter</h2>
        <p>This frontend starter is organized by features, not by page-level sprawl.</p>
      </article>
    </section>
  );
}

export function AppRouter() {
  return (
    <Routes>
      <Route path="/" element={<HomePage />} />
      <Route path="/about" element={<AboutPage />} />
    </Routes>
  );
}
