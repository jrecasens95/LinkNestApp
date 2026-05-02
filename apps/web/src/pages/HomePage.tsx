import { ShortenForm } from "../components/ShortenForm";

export function HomePage() {
  return (
    <main className="page-shell">
      <section className="hero-section">
        <div className="hero-copy">
          <p className="eyebrow">Self-hosted URL shortener</p>
          <h1>LinkNest</h1>
          <p className="hero-lede">
            Create clean short links for campaigns, events, articles, and private team workflows from your own domain.
          </p>
          <div className="hero-metrics" aria-label="Product highlights">
            <span>Private by design</span>
            <span>PostgreSQL backed</span>
            <span>Multi-user auth</span>
          </div>
        </div>

        <ShortenForm />
      </section>
    </main>
  );
}
