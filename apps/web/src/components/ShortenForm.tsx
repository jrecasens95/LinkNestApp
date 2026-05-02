import { FormEvent, useState } from "react";
import { createShortLink, type CreateLinkResponse } from "../api/links";

export function ShortenForm() {
  const [url, setUrl] = useState("");
  const [title, setTitle] = useState("");
  const [customAlias, setCustomAlias] = useState("");
  const [result, setResult] = useState<CreateLinkResponse | null>(null);
  const [error, setError] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [copied, setCopied] = useState(false);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setError("");
    setResult(null);
    setCopied(false);

    const trimmedURL = url.trim();
    if (!trimmedURL) {
      setError("Paste a URL to shorten.");
      return;
    }

    setIsSubmitting(true);

    try {
      const shortLink = await createShortLink({
        original_url: trimmedURL,
        title: title.trim() || undefined,
        custom_alias: customAlias.trim() || undefined
      });

      setResult(shortLink);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Could not shorten the URL.");
    } finally {
      setIsSubmitting(false);
    }
  }

  async function copyShortURL() {
    if (!result) {
      return;
    }

    await navigator.clipboard.writeText(result.short_url);
    setCopied(true);
  }

  return (
    <section className="shorten-panel" aria-label="Shorten URL">
      <form className="shorten-form" onSubmit={handleSubmit}>
        <div className="field-group">
          <label htmlFor="original-url">Long URL</label>
          <input
            id="original-url"
            name="original-url"
            type="url"
            inputMode="url"
            placeholder="https://example.com/newsletter/campaign"
            value={url}
            onChange={(event) => setUrl(event.target.value)}
            disabled={isSubmitting}
          />
        </div>

        <div className="field-group">
          <label htmlFor="link-title">Title</label>
          <input
            id="link-title"
            name="link-title"
            type="text"
            placeholder="Spring launch"
            value={title}
            onChange={(event) => setTitle(event.target.value)}
            disabled={isSubmitting}
          />
        </div>

        <div className="field-group">
          <label htmlFor="custom-alias">Custom alias</label>
          <input
            id="custom-alias"
            name="custom-alias"
            type="text"
            placeholder="spring_launch"
            minLength={3}
            maxLength={40}
            pattern="[A-Za-z0-9_-]{3,40}"
            value={customAlias}
            onChange={(event) => setCustomAlias(event.target.value)}
            disabled={isSubmitting}
          />
        </div>

        {error ? <p className="form-error" role="alert">{error}</p> : null}

        <button className="primary-action" type="submit" disabled={isSubmitting}>
          {isSubmitting ? "Shortening..." : "Shorten URL"}
        </button>
      </form>

      {result ? (
        <div className="result-panel" aria-live="polite">
          <span>Short URL</span>
          <div className="result-row">
            <a href={result.short_url} target="_blank" rel="noreferrer">
              {result.short_url}
            </a>
            <button type="button" onClick={copyShortURL}>
              {copied ? "Copied" : "Copy"}
            </button>
          </div>
        </div>
      ) : null}
    </section>
  );
}
