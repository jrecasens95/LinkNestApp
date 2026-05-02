import { FormEvent, useEffect, useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";
import { deleteLink, getLink, getLinkStats, updateLink, type LinkStats, type ShortLink } from "../api/links";

export function LinkDetailPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [link, setLink] = useState<ShortLink | null>(null);
  const [stats, setStats] = useState<LinkStats | null>(null);
  const [title, setTitle] = useState("");
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(true);
  const [isLoadingStats, setIsLoadingStats] = useState(false);
  const [isSaving, setIsSaving] = useState(false);
  const [copied, setCopied] = useState(false);

  useEffect(() => {
    async function loadLink() {
      if (!id) {
        return;
      }

      setIsLoading(true);
      setError("");

      try {
        const [data, statsData] = await Promise.all([getLink(id), getLinkStats(id)]);
        setLink(data);
        setStats(statsData);
        setTitle(data.title ?? "");
      } catch (err) {
        setError(err instanceof Error ? err.message : "Could not load this link.");
      } finally {
        setIsLoading(false);
      }
    }

    loadLink();
  }, [id]);

  async function refreshStats() {
    if (!id) {
      return;
    }

    setIsLoadingStats(true);
    setError("");

    try {
      const statsData = await getLinkStats(id);
      setStats(statsData);
      setLink((current) => current ? { ...current, clicks_count: statsData.total_clicks } : current);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Could not refresh stats.");
    } finally {
      setIsLoadingStats(false);
    }
  }

  async function saveLink(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!link) {
      return;
    }

    setIsSaving(true);
    setError("");

    try {
      const updated = await updateLink(link.id, { title, is_active: link.is_active });
      setLink(updated);
      setTitle(updated.title ?? "");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Could not save this link.");
    } finally {
      setIsSaving(false);
    }
  }

  async function toggleActive() {
    if (!link) {
      return;
    }

    setIsSaving(true);
    setError("");

    try {
      const updated = await updateLink(link.id, { is_active: !link.is_active });
      setLink(updated);
      setTitle(updated.title ?? "");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Could not update this link.");
    } finally {
      setIsSaving(false);
    }
  }

  async function copyShortURL() {
    if (!link) {
      return;
    }

    await navigator.clipboard.writeText(link.short_url);
    setCopied(true);
  }

  async function removeLink() {
    if (!link || !window.confirm(`Delete ${link.short_url}?`)) {
      return;
    }

    setIsSaving(true);
    setError("");

    try {
      await deleteLink(link.id);
      navigate("/dashboard");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Could not delete this link.");
      setIsSaving(false);
    }
  }

  return (
    <main className="dashboard-shell">
      <header className="dashboard-header">
        <div>
          <p className="eyebrow">Link detail</p>
          <h1>{link?.title || link?.code || "Link"}</h1>
          <p className="dashboard-lede">Edit the title, pause traffic, or remove the link from this instance.</p>
        </div>
        <Link className="secondary-action" to="/dashboard">
          Back
        </Link>
      </header>

      {error ? <p className="form-error" role="alert">{error}</p> : null}

      {isLoading ? (
        <section className="detail-panel">
          <p className="empty-state">Loading link...</p>
        </section>
      ) : link ? (
        <section className="detail-stack">
          <section className="detail-hero-panel">
            <div className="detail-hero-main">
              <span className={link.is_active ? "status-pill active" : "status-pill inactive"}>
                {link.is_active ? "Active" : "Inactive"}
              </span>
              <a className="detail-short-url" href={link.short_url} target="_blank" rel="noreferrer">
                {link.short_url}
              </a>
              <p>{link.original_url}</p>
            </div>

            <div className="detail-hero-actions">
              <button type="button" className="secondary-action button-reset" onClick={copyShortURL}>
                {copied ? "Copied" : "Copy URL"}
              </button>
              <button type="button" className="secondary-action button-reset" onClick={refreshStats} disabled={isLoadingStats}>
                {isLoadingStats ? "Refreshing" : "Refresh stats"}
              </button>
            </div>

            <div className="metric-grid">
              <article>
                <span>Total clicks</span>
                <strong>{stats?.total_clicks ?? link.clicks_count}</strong>
              </article>
              <article>
                <span>Referers</span>
                <strong>{stats?.referers.length ?? 0}</strong>
              </article>
              <article>
                <span>Last click</span>
                <strong>{stats?.recent_clicks[0] ? formatDate(stats.recent_clicks[0].created_at) : "None"}</strong>
              </article>
            </div>
          </section>

          <section className="detail-grid">
            <form className="detail-panel" onSubmit={saveLink}>
              <div className="panel-heading">
                <div>
                  <span className="detail-label">Settings</span>
                  <h2>Link details</h2>
                </div>
              </div>

            <div className="field-group">
              <label htmlFor="detail-title">Title</label>
              <input
                id="detail-title"
                type="text"
                value={title}
                placeholder="Untitled link"
                onChange={(event) => setTitle(event.target.value)}
                disabled={isSaving}
              />
            </div>

            <div className="detail-actions">
              <button className="primary-action" type="submit" disabled={isSaving}>
                {isSaving ? "Saving..." : "Save changes"}
              </button>
              <button type="button" className="secondary-action button-reset" onClick={toggleActive} disabled={isSaving}>
                {link.is_active ? "Disable link" : "Enable link"}
              </button>
              <button type="button" className="danger-button" onClick={removeLink} disabled={isSaving}>
                Delete link
              </button>
            </div>
          </form>

          <aside className="detail-panel stats-panel">
            <div className="panel-heading">
              <div>
                <span className="detail-label">Summary</span>
                <h2>Metadata</h2>
              </div>
            </div>
            <div>
              <span className="detail-label">Short URL</span>
              <a href={link.short_url} target="_blank" rel="noreferrer">{link.short_url}</a>
            </div>
            <div>
              <span className="detail-label">Original URL</span>
              <p>{link.original_url}</p>
            </div>
            <div className="stats-row">
              <div>
                <span className="detail-label">Clicks</span>
                <strong>{stats?.total_clicks ?? link.clicks_count}</strong>
              </div>
              <div>
                <span className="detail-label">Status</span>
                <strong>{link.is_active ? "Active" : "Inactive"}</strong>
              </div>
            </div>
            <div>
              <span className="detail-label">Created</span>
              <p>{formatDateTime(link.created_at)}</p>
            </div>
            <div>
              <span className="detail-label">Updated</span>
              <p>{formatDateTime(link.updated_at)}</p>
            </div>
          </aside>

          <section className="detail-panel analytics-panel">
            <div className="panel-heading">
              <div>
                <span className="detail-label">Traffic sources</span>
                <h2>Referers</h2>
              </div>
            </div>

              {stats?.referers.length ? (
                <div className="referer-list">
                  {stats.referers.map((referer) => (
                    <div className="referer-row" key={referer.referer}>
                      <span>{referer.referer || "Direct"}</span>
                      <strong>{referer.count}</strong>
                    </div>
                  ))}
                </div>
              ) : (
                <p className="empty-inline">No referers yet.</p>
              )}
          </section>

          <section className="detail-panel analytics-panel wide-panel">
            <div className="panel-heading">
              <div>
                <span className="detail-label">Events</span>
                <h2>Latest clicks</h2>
              </div>
            </div>

              {stats?.recent_clicks.length ? (
                <div className="click-list">
                  {stats.recent_clicks.map((click) => (
                    <article className="click-row" key={click.id}>
                      <div>
                        <strong>{formatDateTime(click.created_at)}</strong>
                        <span>{click.referer || "Direct"}</span>
                      </div>
                      <p>{click.user_agent || "Unknown user agent"}</p>
                      <small>{click.ip_address || "Anonymous IP"}</small>
                    </article>
                  ))}
                </div>
              ) : (
                <p className="empty-inline">No click events yet.</p>
              )}
          </section>
          </section>
        </section>
      ) : null}
    </main>
  );
}

function formatDate(value: string) {
  return new Intl.DateTimeFormat(undefined, {
    month: "short",
    day: "numeric"
  }).format(new Date(value));
}

function formatDateTime(value: string) {
  return new Intl.DateTimeFormat(undefined, {
    dateStyle: "medium",
    timeStyle: "short"
  }).format(new Date(value));
}
