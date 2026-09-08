import { useEffect } from "react";
import { NavLink, Outlet } from "react-router-dom";
import {
  clearAccessToken,
  isAuthenticated,
  isMockAuthEnabled,
} from "../api/auth";

const Console = ({ onLogin }) => {
  const isLoggedIn = isAuthenticated();

  useEffect(() => {
    if (!isLoggedIn && !isMockAuthEnabled() && typeof onLogin === "function") {
      onLogin();
    }
  }, [isLoggedIn, onLogin]);

  const handleLogout = () => {
    clearAccessToken();
    window.location.href = "/";
  };

  if (!isLoggedIn) {
    return (
      <div className="auth-redirect-state">
        <h2>Redirecting to sign in…</h2>
      </div>
    );
  }

  return (
    <div className="console-layout">
      <aside className="console-sidebar">
        <div className="sidebar-header">
          <h1>Idiom</h1>
        </div>
        <nav className="sidebar-nav">
          <NavLink to="/console" end>
            Dashboard
          </NavLink>
          <NavLink to="/console/organizations">Organizations</NavLink>
        </nav>
        <div className="sidebar-footer">
          <button onClick={handleLogout}>Log out</button>
        </div>
      </aside>
      <main className="console-main">
        <Outlet />
      </main>
    </div>
  );
};

export default Console;
