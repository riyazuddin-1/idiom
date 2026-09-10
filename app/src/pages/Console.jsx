import { useEffect, useState } from "react";
import { NavLink, Route, Routes } from "react-router-dom";

import {
  clearAccessToken,
  isAuthenticated,
  isMockAuthEnabled,
} from "../api/auth";

import "../styles/Console.css";
import useFetch from "../hooks/useFetch";
import Projects from "./Projects";
import FloatingInput from "../components/FloatingInput";

const Organizations = ({
  organization,
  setOrganization,
  setToggleOrganization,
}) => {
  const {
    data: organizations,
    loading,
    error,
    refetch,
  } = useFetch("/api/v1/organizations");
  const [showForm, setShowForm] = useState(false);
  const [name, setName] = useState("");
  const [creating, setCreating] = useState(false);
  const [formError, setFormError] = useState(null);

  const handleCreate = async (e) => {
    e.preventDefault();
    setCreating(true);
    setFormError(null);

    try {
      const response = await apiFetch("/api/v1/organizations", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name }),
      });

      const result = await response.json();

      if (!result.success) {
        setFormError(result.message);
        return;
      }

      setName("");
      setShowForm(false);
      refetch();
    } catch (err) {
      setFormError("Failed to create organization.");
    } finally {
      setCreating(false);
    }
  };

  useEffect(() => {
    if (!organizations || !organizations.length) return;
    const organizationId = window.sessionStorage.getItem("organizationId");
    if (
      organizationId &&
      organizations.filter((org) => org.id == organizationId).length > 0
    ) {
      setOrganization(
        organizations.filter((org) => org.id == organizationId)[0],
      );
    }
  }, [organizations]);

  return (
    <div
      className="organizations-layout"
      onClick={() => setToggleOrganization(false)}
    >
      <div
        className="organizations-container"
        onClick={(e) => e.stopPropagation()}
      >
        <h2>Organizations</h2>

        {showForm && (
          <form onSubmit={handleCreate} style={{ marginBottom: 20 }}>
            <div className="form-row">
              <div className="form-field">
                <FloatingInput
                  id="org-name"
                  label="Organization name"
                  required
                  autoFocus
                />
              </div>
            </div>
            {formError && (
              <p className="error" style={{ marginTop: 8 }}>
                {formError}
              </p>
            )}
          </form>
        )}

        {loading && <div className="loading">Loading...</div>}
        {error && <div className="error">Failed to load organizations.</div>}

        {!loading && !error && organizations && organizations.length === 0 && (
          <div className="empty">No organizations yet.</div>
        )}

        {!loading &&
          !error &&
          !showForm &&
          organizations &&
          organizations.length > 0 && (
            <ul className="organizations-list">
              {organizations.map((org) => (
                <li
                  key={org.id}
                  onClick={() => {
                    if (org.id != organization?.id) {
                      setOrganization(org);
                    }
                    setToggleOrganization(false);
                  }}
                >
                  {org.name}
                </li>
              ))}
            </ul>
          )}

        {!loading && (
          <>
            <hr />
            <button
              className="toggle-create"
              onClick={() => setShowForm(!showForm)}
            >
              {showForm ? "Cancel" : "Create"}
            </button>
          </>
        )}
      </div>
    </div>
  );
};

const Console = ({ onLogin }) => {
  const isLoggedIn = isAuthenticated();
  const [organization, setOrganization] = useState(null);
  const [toggleOrganization, setToggleOrganization] = useState(false);

  useEffect(() => {
    if (!isLoggedIn && !isMockAuthEnabled() && typeof onLogin === "function") {
      onLogin();
    }
  }, [isLoggedIn, onLogin]);

  useEffect(() => {
    if (organization) {
      window.sessionStorage.setItem("organizationId", organization.id);
    }
  }, [organization]);

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
    <div className="console">
      {!organization ? (
        <Organizations
          organization={organization}
          setOrganization={setOrganization}
          setToggleOrganization={setToggleOrganization}
        />
      ) : (
        <div className="console-layout">
          <aside className="console-sidebar">
            <div className="sidebar-header">
              <h1 onClick={() => setToggleOrganization(true)}>
                {organization.name}
              </h1>
            </div>

            <nav className="sidebar-nav">
              <NavLink to="/console/projects">Projects</NavLink>
            </nav>

            <div className="sidebar-footer">
              <button onClick={handleLogout}>Log out</button>
            </div>
          </aside>

          <main className="console-main">
            <Routes>
              <Route
                path="/projects"
                element={<Projects organization={organization} />}
              />
            </Routes>
          </main>
        </div>
      )}

      {organization && toggleOrganization && (
        <Organizations
          organization={organization}
          setOrganization={setOrganization}
          setToggleOrganization={setToggleOrganization}
        />
      )}
    </div>
  );
};

export default Console;
