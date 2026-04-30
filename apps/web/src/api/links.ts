const API_URL = import.meta.env.VITE_API_URL ?? "http://localhost:4000";

export type CreateLinkPayload = {
  original_url: string;
  title?: string;
};

export type CreateLinkResponse = {
  code: string;
  short_url: string;
};

type APIErrorResponse = {
  error?: string;
};

export async function createShortLink(payload: CreateLinkPayload): Promise<CreateLinkResponse> {
  const response = await fetch(`${API_URL}/api/links`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify(payload)
  });

  if (!response.ok) {
    let message = "Could not shorten the URL. Please try again.";

    try {
      const data = (await response.json()) as APIErrorResponse;
      if (data.error) {
        message = data.error;
      }
    } catch {
      // Keep the generic message when the API does not return JSON.
    }

    throw new Error(message);
  }

  return response.json() as Promise<CreateLinkResponse>;
}
