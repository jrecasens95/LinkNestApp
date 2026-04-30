const API_URL = import.meta.env.VITE_API_URL ?? "http://localhost:4000";
const API_KEY = import.meta.env.VITE_API_KEY;

export type CreateLinkPayload = {
  original_url: string;
  title?: string;
};

export type CreateLinkResponse = {
  code: string;
  short_url: string;
};

export type ShortLink = {
  id: number;
  code: string;
  original_url: string;
  title?: string;
  clicks_count: number;
  is_active: boolean;
  short_url: string;
  created_at: string;
  updated_at: string;
};

export type ListLinksResponse = {
  links: ShortLink[];
};

export type UpdateLinkPayload = {
  title?: string;
  is_active?: boolean;
};

export type ClickEvent = {
  id: number;
  user_agent: string;
  referer: string;
  ip_address: string;
  created_at: string;
};

export type RefererStat = {
  referer: string;
  count: number;
};

export type LinkStats = {
  total_clicks: number;
  recent_clicks: ClickEvent[];
  referers: RefererStat[];
};

type APIErrorResponse = {
  error?: string;
};

async function parseAPIError(response: Response, fallback: string) {
  let message = fallback;

  try {
    const data = (await response.json()) as APIErrorResponse;
    if (data.error) {
      message = data.error;
    }
  } catch {
    // Keep the fallback message when the API does not return JSON.
  }

  return new Error(message);
}

export async function createShortLink(payload: CreateLinkPayload): Promise<CreateLinkResponse> {
  const headers: Record<string, string> = {
    "Content-Type": "application/json"
  };

  if (API_KEY) {
    headers["X-API-Key"] = API_KEY;
  }

  const response = await fetch(`${API_URL}/api/links`, {
    method: "POST",
    headers,
    body: JSON.stringify(payload)
  });

  if (!response.ok) {
    throw await parseAPIError(response, "Could not shorten the URL. Please try again.");
  }

  return response.json() as Promise<CreateLinkResponse>;
}

export async function listLinks(): Promise<ShortLink[]> {
  const response = await fetch(`${API_URL}/api/links`);

  if (!response.ok) {
    throw await parseAPIError(response, "Could not load links.");
  }

  const data = (await response.json()) as ListLinksResponse;
  return data.links;
}

export async function getLink(id: string): Promise<ShortLink> {
  const response = await fetch(`${API_URL}/api/links/${id}`);

  if (!response.ok) {
    throw await parseAPIError(response, "Could not load this link.");
  }

  return response.json() as Promise<ShortLink>;
}

export async function updateLink(id: number, payload: UpdateLinkPayload): Promise<ShortLink> {
  const response = await fetch(`${API_URL}/api/links/${id}`, {
    method: "PATCH",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify(payload)
  });

  if (!response.ok) {
    throw await parseAPIError(response, "Could not update this link.");
  }

  return response.json() as Promise<ShortLink>;
}

export async function deleteLink(id: number): Promise<void> {
  const response = await fetch(`${API_URL}/api/links/${id}`, {
    method: "DELETE"
  });

  if (!response.ok) {
    throw await parseAPIError(response, "Could not delete this link.");
  }
}

export async function getLinkStats(id: string): Promise<LinkStats> {
  const response = await fetch(`${API_URL}/api/links/${id}/stats`);

  if (!response.ok) {
    throw await parseAPIError(response, "Could not load link stats.");
  }

  return response.json() as Promise<LinkStats>;
}
