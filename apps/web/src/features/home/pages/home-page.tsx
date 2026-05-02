import { Hero } from "../components";
import { FeatureShell } from "../../../shared/components";
import { useDocumentTitle } from "../../../shared/hooks";

export function HomePage() {
  useDocumentTitle("Link Nest App - Home");

  return (
    <FeatureShell>
      <Hero />
      <section className="card-grid">
        <article className="card">
          <h2>Feature-first structure</h2>
          <p>
            Pages, components, hooks, and utilities can grow behind explicit
            module exports.
          </p>
        </article>
        <article className="card">
          <h2>Public API by design</h2>
          <p>
            Import from folder indexes instead of reaching into internals across
            the app.
          </p>
        </article>
      </section>
    </FeatureShell>
  );
}
