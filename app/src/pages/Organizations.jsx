import { useState } from "react";
import { Link } from "react-router-dom";
import useFetch from "../hooks/useFetch";
import apiFetch from "../api/client";

const Organizations = () => {
  const { data: organizations, loading, error, refetch } = useFetch("/api/v1/organizations");
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

  return (
    <div>
      <h2>Organizations</h2>

      <div className="toolbar">
        <span />
        <button className="btn btn-primary" onClick={() => setShowForm(!showForm)}>
          {showForm ? "Cancel" : "New Organization"}
        </button>
      </div>

      {showForm && (
        <form onSubmit={handleCreate} style={{ marginBottom: 20 }}>
          <div className="form-row">
            <div className="form-field">
              <label htmlFor="org-name">Name</label>
              <input
                id="org-name"
                type="text"
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="Organization name"
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
      {error && <div className="error">Failed to load organizations.</div>}

      {!loading && !error && organizations && organizations.length === 0 && (
        <div className="empty">No organizations yet.</div>
      )}

      {!loading && !error && organizations && organizations.length > 0 && (
        <div className="table-wrapper">
          <table>
            <thead>
              <tr>
                <th>Name</th>
                <th>Status</th>
                <th>Created</th>
                <th />
              </tr>
            </thead>
            <tbody>
              {organizations.map((org) => (
                <tr key={org.id}>
                  <td>{org.name}</td>
                  <td>
                    <span className="badge">{org.status}</span>
                  </td>
                  <td>{new Date(org.createdAt).toLocaleDateString()}</td>
                  <td>
                    <Link to={`/console/organizations/${org.id}/projects`}>
                      Projects
                    </Link>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
};

export default Organizations;
