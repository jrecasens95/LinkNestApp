import { clearAuthToken, getAuthToken, setAuthToken } from "./links";

const API_URL = import.meta.env.VITE_API_URL ?? "http://localhost:4000";
const API_KEY = import.meta.env.VITE_API_KEY;

export type User = {
  id: number;
  name: string;
  email: string;
  created_at: string;
};

export type AuthResponse = {
  token: string;
  user: User;
};

type AuthPayload = {
  name?: string;
  email: string;
  password: string;
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
    // Keep fallback.
  }
  return new Error(message);
}

export async function register(payload: AuthPayload) {
  return authRequest("/api/auth/register", payload);
}

export async function login(payload: AuthPayload) {
  return authRequest("/api/auth/login", payload);
}

export async function me(): Promise<User> {
  const token = getAuthToken();
  const response = await fetch(`${API_URL}/api/auth/me`, {
    headers: token ? { Authorization: `Bearer ${token}` } : {}
  });

  if (!response.ok) {
    clearAuthToken();
    throw await parseAPIError(response, "Please sign in again.");
  }

  return response.json() as Promise<User>;
}

async function authRequest(path: string, payload: AuthPayload) {
  const headers: Record<string, string> = {
    "Content-Type": "application/json"
  };

  if (API_KEY && path === "/api/auth/register") {
    headers["X-API-Key"] = API_KEY;
  }

  const response = await fetch(`${API_URL}${path}`, {
    method: "POST",
    headers,
    body: JSON.stringify(payload)
  });

  if (!response.ok) {
    throw await parseAPIError(response, "Authentication failed.");
  }

  const data = (await response.json()) as AuthResponse;
  setAuthToken(data.token);
  return data;
}
