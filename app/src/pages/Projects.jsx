import { useState } from "react";
import { useParams, Link } from "react-router-dom";
import useFetch from "../hooks/useFetch";
import apiFetch from "../api/client";

const Projects = () => {
  const { oid } = useParams();
  const {
    data: projects,
    loading,
    error,
    refetch,
  } = useFetch(`/api/v1/organizations/${oid}/projects`);
  const [showForm, setShowForm] = useState(false);
  const [name, setName] = useState("");
  const [creating, setCreating] = useState(false);
  const [formError, setFormError] = useState(null);

  const handleCreate = async (e) => {
    e.preventDefault();
    setCreating(true);
    setFormError(null);

    try {
      const response = await apiFetch(`/api/v1/organizations/${oid}/projects`, {
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
      setFormError("Failed to create project.");
    } finally {
      setCreating(false);
    }
  };

  return (
    <div>
      <div style={{ marginBottom: 16 }}>
        <Link
          to="/console/organizations"
          style={{ fontSize: 14, color: "var(--text)" }}
        >
          &larr; Organizations
        </Link>
      </div>

      <h2>Projects</h2>

      <div className="toolbar">
        <span />
        <button className="btn btn-primary" onClick={() => setShowForm(!showForm)}>
          {showForm ? "Cancel" : "New Project"}
        </button>
      </div>

      {showForm && (
        <form onSubmit={handleCreate} style={{ marginBottom: 20 }}>
          <div className="form-row">
            <div className="form-field">
              <label htmlFor="project-name">Name</label>
              <input
                id="project-name"
                type="text"
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="Project name"
                required
              />
            </div>
            <button className="btn btn-primary" type="submit" disabled={creating}>
              {creating ? "Creating..." : "Create"}
            </button>
          </div>
          {formError && <p className="error" style={{ marginTop: 8 }}>{formError}</p>}
        </form>
      )}

      {loading && <div className="loading">Loading...</div>}
      {error && <div className="error">Failed to load projects.</div>}

      {!loading && !error && projects && projects.length === 0 && (
        <div className="empty">No projects in this organization yet.</div>
      )}

      {!loading && !error && projects && projects.length > 0 && (
        <div className="table-wrapper">
          <table>
            <thead>
              <tr>
                <th>Name</th>
                <th>Status</th>
                <th>Created</th>
              </tr>
            </thead>
            <tbody>
              {projects.map((project) => (
                <tr key={project.id}>
                  <td>{project.name}</td>
                  <td>
                    <span className="badge">{project.status || "active"}</span>
                  </td>
                  <td>{new Date(project.createdAt).toLocaleDateString()}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
};

export default Projects;
