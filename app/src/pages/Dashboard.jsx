import useFetch from "../hooks/useFetch";

const Dashboard = () => {
  const {
    data: organizations,
    loading,
    error,
  } = useFetch("/api/v1/organizations");

  if (loading) return <div className="loading">Loading...</div>;
  if (error) return <div className="error">Failed to load dashboard.</div>;

  const orgCount = organizations ? organizations.length : 0;

  return (
    <div>
      <h2>Dashboard</h2>
      <div className="card-grid">
        <div className="card">
          <div className="card-label">Organizations</div>
          <div className="card-value">{orgCount}</div>
        </div>
      </div>
      {orgCount === 0 && (
        <p style={{ color: "var(--text)" }}>
          No organizations yet. Create one from the{" "}
          <a href="/console/organizations">Organizations</a> page.
        </p>
      )}
    </div>
  );
};

export default Dashboard;
