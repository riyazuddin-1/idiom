import { useState } from "react";
import { useParams, Link } from "react-router-dom";
import useFetch from "../hooks/useFetch";
import apiFetch from "../api/client";
import "../styles/Project.css";
import FloatingInput from "../components/FloatingInput";

const ProjectSettings = ({ organization, project, onUpdated }) => {
  const [editing, setEditing] = useState(false);
  const [name, setName] = useState(project.name);
  const [description, setDescription] = useState(project.description || "");
  const [redirectUrls, setRedirectUrls] = useState(
    project.redirect_urls || [],
  );
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState(null);

  const handleSave = async (e) => {
    e.preventDefault();
    setSaving(true);
    setError(null);

    try {
      const response = await apiFetch(
        `/api/v1/organizations/${organization.id}/projects/${project.id}`,
        {
          method: "PUT",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            name,
            description,
            redirect_urls: redirectUrls,
          }),
        },
      );

      const result = await response.json();

      if (!result.success) {
        setError(result.message);
        return;
      }

      setEditing(false);
      onUpdated(result.data);
    } catch {
      setError("Failed to update project.");
    } finally {
      setSaving(false);
    }
  };

  const handleCancel = () => {
    setName(project.name);
    setDescription(project.description || "");
    setRedirectUrls(project.redirect_urls || []);
    setError(null);
    setEditing(false);
  };

  const handleRedirectUrlChange = (index, field, value) => {
    setRedirectUrls((prev) =>
      prev.map((item, i) =>
        i === index ? { ...item, [field]: value } : item,
      ),
    );
  };

  const handleAddRedirectUrl = () => {
    setRedirectUrls((prev) => [...prev, { device: "web", url: "" }]);
  };

  const handleRemoveRedirectUrl = (index) => {
    setRedirectUrls((prev) => prev.filter((_, i) => i !== index));
  };

  if (editing) {
    return (
      <div className="section">
        <div className="section-header">
          <h3>Settings</h3>
        </div>

        <form onSubmit={handleSave} className="section-form">
          <div className="form-row">
            <div className="form-field">
              <FloatingInput
                id="edit-project-name"
                label="Project name"
                required
                autoFocus
                value={name}
                onChange={(e) => setName(e.target.value)}
              />
            </div>
          </div>

          <div className="form-row" style={{ marginTop: 12 }}>
            <div className="form-field" style={{ width: "100%" }}>
              <FloatingInput
                id="edit-project-desc"
                label="Description"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
              />
            </div>
          </div>

          <div className="form-row" style={{ marginTop: 20 }}>
            <div className="form-field" style={{ width: "100%" }}>
              <div
                className="section-header"
                style={{ marginBottom: 12 }}
              >
                <div>
                  <h4>Redirect URLs</h4>
                  <p>Configure redirect URLs for web and mobile devices.</p>
                </div>

                <button
                  type="button"
                  className="btn btn-primary"
                  onClick={handleAddRedirectUrl}
                >
                  Add URL
                </button>
              </div>

              {redirectUrls.length === 0 && (
                <p>No redirect URLs configured.</p>
              )}

              {redirectUrls.map((redirect, index) => (
                <div
                  key={index}
                  className="form-row"
                  style={{
                    display: "flex",
                    alignItems: "center",
                    gap: 12,
                    marginBottom: 12,
                  }}
                >
                  <div className="form-field" style={{ minWidth: 120 }}>
                    <label htmlFor={`redirect-device-${index}`}>
                      Device
                    </label>
                    <select
                      id={`redirect-device-${index}`}
                      value={redirect.device}
                      onChange={(e) =>
                        handleRedirectUrlChange(
                          index,
                          "device",
                          e.target.value,
                        )
                      }
                    >
                      <option value="web">Web</option>
                      <option value="mobile">Mobile</option>
                    </select>
                  </div>

                  <div className="form-field" style={{ flex: 1 }}>
                    <FloatingInput
                      id={`redirect-url-${index}`}
                      label="Redirect URL"
                      type="url"
                      required
                      value={redirect.url}
                      onChange={(e) =>
                        handleRedirectUrlChange(
                          index,
                          "url",
                          e.target.value,
                        )
                      }
                    />
                  </div>

                  <button
                    type="button"
                    className="btn-cancel"
                    onClick={() => handleRemoveRedirectUrl(index)}
                    aria-label={`Remove redirect URL ${index + 1}`}
                  >
                    Remove
                  </button>
                </div>
              ))}
            </div>
          </div>

          {error && (
            <p className="error" style={{ marginTop: 8 }}>
              {error}
            </p>
          )}

          <div className="form-actions" style={{ marginTop: 12 }}>
            <button
              type="button"
              className="btn-cancel"
              onClick={handleCancel}
            >
              Cancel
            </button>

            <button
              type="submit"
              className="btn btn-primary"
              disabled={saving}
            >
              {saving ? "Saving..." : "Save"}
            </button>
          </div>
        </form>
      </div>
    );
  }

  return (
    <div className="section">
      <div className="section-header">
        <div>
          <h3>Settings</h3>
          <p>Project name, description, and redirect URLs.</p>
        </div>

        <button
          className="btn btn-primary"
          onClick={() => setEditing(true)}
        >
          Edit
        </button>
      </div>

      <div className="settings-display">
        <div className="settings-row">
          <span className="settings-label">Name</span>
          <span>{project.name}</span>
        </div>

        <div className="settings-row">
          <span className="settings-label">Description</span>
          <span>{project.description || "—"}</span>
        </div>

        <div className="settings-row">
          <span className="settings-label">Redirect URLs</span>
          <div>
            {project.redirect_urls?.length > 0 ? (
              project.redirect_urls.map((redirect, index) => (
                <div key={index} style={{ marginBottom: 12 }}>
                  <p>
                    <b>
                      {redirect.device === "web" ? "Web" : "Mobile"}:
                    </b>
                  </p>
                  <p>{redirect.url}</p>
                </div>
              ))
            ) : (
              <span>—</span>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

// const ApiKeysSection = ({ organization, project }) => {
//   const {
//     data: apiKeys,
//     loading,
//     error,
//     refetch,
//   } = useFetch(
//     `/api/v1/organizations/${organization.id}/projects/${project.id}/api-keys`,
//   );

//   const [showForm, setShowForm] = useState(false);
//   const [name, setName] = useState("");
//   const [scopes, setScopes] = useState("");
//   const [creating, setCreating] = useState(false);
//   const [formError, setFormError] = useState(null);
//   const [newKey, setNewKey] = useState(null);
//   const [copied, setCopied] = useState(false);
//   const [editingId, setEditingId] = useState(null);
//   const [editName, setEditName] = useState("");
//   const [editScopes, setEditScopes] = useState("");

//   const handleCreate = async (e) => {
//     e.preventDefault();
//     setCreating(true);
//     setFormError(null);
//     setNewKey(null);

//     try {
//       const scopeList = scopes
//         .split(",")
//         .map((s) => s.trim())
//         .filter(Boolean);

//       const response = await apiFetch(
//         `/api/v1/organizations/${organization.id}/projects/${project.id}/api-keys`,
//         {
//           method: "POST",
//           headers: { "Content-Type": "application/json" },
//           body: JSON.stringify({ name, scopes: scopeList }),
//         },
//       );

//       const result = await response.json();

//       if (!result.success) {
//         setFormError(result.message);
//         return;
//       }

//       setNewKey(result.data.key);
//       setName("");
//       setScopes("");
//       setShowForm(false);
//       refetch();
//     } catch {
//       setFormError("Failed to create API key.");
//     } finally {
//       setCreating(false);
//     }
//   };

//   const handleRotate = async (keyId) => {
//     if (!window.confirm("Rotate this key? The old key will stop working immediately.")) return;

//     try {
//       const response = await apiFetch(
//         `/api/v1/organizations/${organization.id}/projects/${project.id}/api-keys/${keyId}/rotate`,
//         { method: "POST" },
//       );

//       const result = await response.json();

//       if (!result.success) {
//         alert(result.message);
//         return;
//       }

//       setNewKey(result.data.key);
//       refetch();
//     } catch {
//       alert("Failed to rotate API key.");
//     }
//   };

//   const handleRevoke = async (keyId) => {
//     if (!window.confirm("Revoke this key? This cannot be undone.")) return;

//     try {
//       const response = await apiFetch(
//         `/api/v1/organizations/${organization.id}/projects/${project.id}/api-keys/${keyId}`,
//         { method: "DELETE" },
//       );

//       if (!response.ok && response.status !== 204) {
//         const result = await response.json();
//         alert(result.message || "Failed to revoke API key.");
//         return;
//       }

//       refetch();
//     } catch {
//       alert("Failed to revoke API key.");
//     }
//   };

//   const startEdit = (key) => {
//     setEditingId(key.id);
//     setEditName(key.name);
//     setEditScopes((key.scopes || []).join(", "));
//   };

//   const cancelEdit = () => {
//     setEditingId(null);
//     setEditName("");
//     setEditScopes("");
//   };

//   const handleEditSave = async (keyId) => {
//     const scopeList = editScopes
//       .split(",")
//       .map((s) => s.trim())
//       .filter(Boolean);

//     try {
//       const response = await apiFetch(
//         `/api/v1/organizations/${organization.id}/projects/${project.id}/api-keys/${keyId}`,
//         {
//           method: "PUT",
//           headers: { "Content-Type": "application/json" },
//           body: JSON.stringify({ name: editName, scopes: scopeList }),
//         },
//       );

//       const result = await response.json();

//       if (!result.success) {
//         alert(result.message);
//         return;
//       }

//       cancelEdit();
//       refetch();
//     } catch {
//       alert("Failed to update API key.");
//     }
//   };

//   const copyToClipboard = (text) => {
//     navigator.clipboard.writeText(text);
//     setCopied(true);
//     setTimeout(() => setCopied(false), 2000);
//   };

//   const formatDate = (dateStr) => {
//     if (!dateStr) return "—";
//     return new Date(dateStr).toLocaleDateString();
//   };

//   return (
//     <div className="section">
//       <div className="section-header">
//         <div>
//           <h3>API Keys</h3>
//           <p>Authenticate server-side requests to your project.</p>
//         </div>
//         <button
//           className="btn btn-primary"
//           onClick={() => {
//             setShowForm(!showForm);
//             setNewKey(null);
//           }}
//         >
//           {showForm ? "Cancel" : "+ New key"}
//         </button>
//       </div>

//       {newKey && (
//         <div className="new-key-banner">
//           <h4>Key created</h4>
//           <p>Copy this key now. You won't be able to see it again.</p>
//           <div className="new-key-row">
//             <code>{newKey}</code>
//             <button onClick={() => copyToClipboard(newKey)}>
//               {copied ? "Copied" : "Copy"}
//             </button>
//           </div>
//         </div>
//       )}

//       {showForm && (
//         <form onSubmit={handleCreate} className="section-form">
//           <div className="form-row">
//             <div className="form-field">
//               <FloatingInput
//                 id="key-name"
//                 label="Key name"
//                 required
//                 autoFocus
//                 value={name}
//                 onChange={(e) => setName(e.target.value)}
//               />
//             </div>
//             <div className="form-field">
//               <FloatingInput
//                 id="key-scopes"
//                 label="Scopes (comma-separated, optional)"
//                 value={scopes}
//                 onChange={(e) => setScopes(e.target.value)}
//               />
//             </div>
//           </div>
//           {formError && <p className="error">{formError}</p>}
//           <button
//             type="submit"
//             className="btn btn-primary"
//             disabled={creating}
//             style={{ marginTop: 12 }}
//           >
//             {creating ? "Creating..." : "Create key"}
//           </button>
//         </form>
//       )}

//       {loading && <div className="loading">Loading...</div>}
//       {error && <div className="error">Failed to load API keys.</div>}
//     </div>
//   );
// };

const Project = ({ organization }) => {
  const { pid } = useParams();
  const {
    data: project,
    loading,
    error,
    refetch,
  } = useFetch(`/api/v1/organizations/${organization.id}/projects/${pid}`);

  if (loading) {
    return (
      <div className="project-page">
        <div className="loading">Loading...</div>
      </div>
    );
  }

  if (error || !project) {
    return (
      <div className="project-page">
        <Link to="/console/projects" className="btn-back">
          &larr; Projects
        </Link>
        <div className="error">Failed to load project.</div>
      </div>
    );
  }

  return (
    <div className="project-page">
      <Link to="/console/projects" className="btn-back">
        &larr; Projects
      </Link>

      <div className="project-header">
        <div className="project-icon">
          {(project.name || "?").charAt(0).toUpperCase()}
        </div>
        <div>
          <h2>{project.name}</h2>
          <p>{project.slug || project.id}</p>
        </div>
      </div>

      <ProjectSettings
        organization={organization}
        project={project}
        onUpdated={() => refetch()}
      />

      {/* <ApiKeysSection organization={organization} project={project} /> */}
    </div>
  );
};

export default Project;
