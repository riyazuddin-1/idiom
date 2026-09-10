import { useState } from "react";
import { useParams } from "react-router-dom";
import useFetch from "../hooks/useFetch";
import apiFetch from "../api/client";
import "../styles/Projects.css";
import FloatingInput from "../components/FloatingInput";

const Projects = ({organization}) => {
  const {
    data: projects,
    loading,
    error,
    refetch,
  } = useFetch(`/api/v1/organizations/${organization.id}/projects`);
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
    } catch {
      setFormError("Failed to create project.");
    } finally {
      setCreating(false);
    }
  };

  const statusLabel = (status) => {
    if (!status) return "Active";
    return status.charAt(0).toUpperCase() + status.slice(1);
  };

  const projectSub = (project) =>
    project.slug || new Date(project.createdAt).toLocaleDateString();

  return (
    <div className="projects-page">
      <div className="projects-header">
        <div>
          <h2>Projects</h2>
          <p>Each product gets a dedicated auth workspace.</p>
        </div>
        <button
          className="btn btn-primary"
          onClick={() => setShowForm(!showForm)}
        >
          {showForm ? "Cancel" : "+ New project"}
        </button>
      </div>

      {showForm && (
        <form onSubmit={handleCreate} className="projects-form">
          <FloatingInput
            id="project-name"
            label="Project name"
            required
            autoFocus
          />
          {formError && <p className="error">{formError}</p>}
        </form>
      )}

      {loading && <div className="loading">Loading...</div>}
      {error && <div className="error">Failed to load projects.</div>}

      {!loading && !error && projects && projects.length === 0 && (
        <div className="empty">No projects in this organization yet.</div>
      )}

      {!loading && !error && projects && projects.length > 0 && (
        <div className="projects-grid">
          {projects.map((project) => (
            <div className="project-tile" key={project.id}>
              <div className="project-tile-icon">
                {(project.name || "?").charAt(0).toUpperCase()}
              </div>
              <h4>{project.name}</h4>
              <p>{projectSub(project)}</p>
              <div className="project-tile-status">
                <span className="status-dot" />
                {statusLabel(project.status)}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default Projects;
