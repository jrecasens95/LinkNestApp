import { FormEvent, useEffect, useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";
import { deleteLink, getLink, updateLink, type ShortLink } from "../api/links";

export function LinkDetailPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [link, setLink] = useState<ShortLink | null>(null);
  const [title, setTitle] = useState("");
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);

  useEffect(() => {
    async function loadLink() {
      if (!id) {
        return;
      }

      setIsLoading(true);
      setError("");

      try {
        const data = await getLink(id);
        setLink(data);
        setTitle(data.title ?? "");
      } catch (err) {
        setError(err instanceof Error ? err.message : "Could not load this link.");
      } finally {
        setIsLoading(false);
      }
    }

    loadLink();
  }, [id]);

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
        <section className="detail-grid">
          <form className="detail-panel" onSubmit={saveLink}>
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
                Save changes
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
                <strong>{link.clicks_count}</strong>
              </div>
              <div>
                <span className="detail-label">Status</span>
                <strong>{link.is_active ? "Active" : "Inactive"}</strong>
              </div>
            </div>
            <div>
              <span className="detail-label">Created</span>
              <p>{new Date(link.created_at).toLocaleString()}</p>
            </div>
          </aside>
        </section>
      ) : null}
    </main>
  );
}
