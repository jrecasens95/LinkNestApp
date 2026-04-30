import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { deleteLink, listLinks, updateLink, type ShortLink } from "../api/links";

export function DashboardPage() {
  const [links, setLinks] = useState<ShortLink[]>([]);
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(true);
  const [busyID, setBusyID] = useState<number | null>(null);

  useEffect(() => {
    loadLinks();
  }, []);

  async function loadLinks() {
    setIsLoading(true);
    setError("");

    try {
      setLinks(await listLinks());
    } catch (err) {
      setError(err instanceof Error ? err.message : "Could not load links.");
    } finally {
      setIsLoading(false);
    }
  }

  async function toggleLink(link: ShortLink) {
    setBusyID(link.id);
    setError("");

    try {
      const updated = await updateLink(link.id, { is_active: !link.is_active });
      setLinks((current) => current.map((item) => (item.id === updated.id ? updated : item)));
    } catch (err) {
      setError(err instanceof Error ? err.message : "Could not update this link.");
    } finally {
      setBusyID(null);
    }
  }

  async function removeLink(link: ShortLink) {
    const confirmed = window.confirm(`Delete ${link.short_url}?`);
    if (!confirmed) {
      return;
    }

    setBusyID(link.id);
    setError("");

    try {
      await deleteLink(link.id);
      setLinks((current) => current.filter((item) => item.id !== link.id));
    } catch (err) {
      setError(err instanceof Error ? err.message : "Could not delete this link.");
    } finally {
      setBusyID(null);
    }
  }

  return (
    <main className="dashboard-shell">
      <header className="dashboard-header">
        <div>
          <p className="eyebrow">Private dashboard</p>
          <h1>Links</h1>
          <p className="dashboard-lede">Manage every short link in this single-instance workspace.</p>
        </div>
        <Link className="secondary-action" to="/">
          New link
        </Link>
      </header>

      {error ? <p className="form-error" role="alert">{error}</p> : null}

      <section className="table-panel">
        {isLoading ? (
          <p className="empty-state">Loading links...</p>
        ) : links.length === 0 ? (
          <p className="empty-state">No links yet. Create the first one from the landing page.</p>
        ) : (
          <div className="table-scroll">
            <table className="links-table">
              <thead>
                <tr>
                  <th>Link</th>
                  <th>Status</th>
                  <th>Clicks</th>
                  <th>Created</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                {links.map((link) => (
                  <tr key={link.id}>
                    <td>
                      <Link className="link-title" to={`/dashboard/links/${link.id}`}>
                        {link.title || link.code}
                      </Link>
                      <a className="short-url" href={link.short_url} target="_blank" rel="noreferrer">
                        {link.short_url}
                      </a>
                      <span className="original-url">{link.original_url}</span>
                    </td>
                    <td>
                      <span className={link.is_active ? "status-pill active" : "status-pill inactive"}>
                        {link.is_active ? "Active" : "Inactive"}
                      </span>
                    </td>
                    <td>{link.clicks_count}</td>
                    <td>{formatDate(link.created_at)}</td>
                    <td>
                      <div className="table-actions">
                        <button type="button" onClick={() => toggleLink(link)} disabled={busyID === link.id}>
                          {link.is_active ? "Disable" : "Enable"}
                        </button>
                        <button type="button" className="danger-action" onClick={() => removeLink(link)} disabled={busyID === link.id}>
                          Delete
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>
    </main>
  );
}

function formatDate(value: string) {
  return new Intl.DateTimeFormat(undefined, {
    month: "short",
    day: "numeric",
    year: "numeric"
  }).format(new Date(value));
}
