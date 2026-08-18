import { NavLink, Outlet } from "react-router-dom";
import { getAccessToken, clearAccessToken } from "../api/auth";

const Console = ({ onLogin }) => {
  const isLoggedIn = !!getAccessToken();

  const handleLogout = () => {
    clearAccessToken();
    window.location.href = "/";
  };

  if (!isLoggedIn) {
    return (
      <div style={{ textAlign: "center", padding: "4rem 1rem" }}>
        <h2>Welcome to Idiom</h2>
        <p style={{ color: "var(--text)", margin: "12px 0 24px" }}>
          Sign in to access your console.
        </p>
        <button className="btn btn-primary" onClick={onLogin}>
          Sign In
        </button>
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
