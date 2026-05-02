import type { ReactNode } from "react";
import { Link, Navigate, Route, Routes, useLocation, useNavigate } from "react-router-dom";
import { clearAuthToken, getAuthToken } from "./api/links";
import { AuthPage } from "./pages/AuthPage";
import { DashboardPage } from "./pages/DashboardPage";
import { HomePage } from "./pages/HomePage";
import { LinkDetailPage } from "./pages/LinkDetailPage";

export function App() {
  const isSignedIn = Boolean(getAuthToken());
  const navigate = useNavigate();

  function logout() {
    clearAuthToken();
    navigate("/login", { replace: true });
  }

  return (
    <>
      <nav className="top-nav">
        <Link to="/">LinkNest</Link>
        <div>
          <Link to="/">Create</Link>
          <Link to="/dashboard">Dashboard</Link>
          {isSignedIn ? (
            <button type="button" onClick={logout}>
              Logout
            </button>
          ) : (
            <Link to="/login">Login</Link>
          )}
        </div>
      </nav>
      <Routes>
        <Route path="/" element={<ProtectedRoute><HomePage /></ProtectedRoute>} />
        <Route path="/login" element={<PublicOnlyRoute><AuthPage mode="login" /></PublicOnlyRoute>} />
        <Route path="/register" element={<PublicOnlyRoute><AuthPage mode="register" /></PublicOnlyRoute>} />
        <Route path="/dashboard" element={<ProtectedRoute><DashboardPage /></ProtectedRoute>} />
        <Route path="/dashboard/links/:id" element={<ProtectedRoute><LinkDetailPage /></ProtectedRoute>} />
      </Routes>
    </>
  );
}

function ProtectedRoute({ children }: { children: ReactNode }) {
  const location = useLocation();

  if (!getAuthToken()) {
    return <Navigate to="/login" replace state={{ from: location }} />;
  }

  return children;
}

function PublicOnlyRoute({ children }: { children: ReactNode }) {
  const location = useLocation();
  const from = location.state?.from?.pathname ?? "/dashboard";

  if (getAuthToken()) {
    return <Navigate to={from} replace />;
  }

  return children;
}
