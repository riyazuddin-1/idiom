import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";
import { isMockAuthEnabled } from "./api/auth";
import Callback from "./pages/Callback";
import Console from "./pages/Console";
import Dashboard from "./pages/Dashboard";
import HomePage from "./pages/Home";
import Organizations from "./pages/Organizations";
import Projects from "./pages/Projects";
import "./App.css";

const AUTH_LOGIN_URL =
  import.meta.env.VITE_AUTH_LOGIN_URL ||
  "https://idiom-identity.onrender.com/idiom/login";

function App() {
  const mockAuthEnabled = isMockAuthEnabled();

  return (
    <>
      {mockAuthEnabled && (
        <div className="mock-auth-banner">
          Development mode: mock auth enabled
        </div>
      )}

      <BrowserRouter>
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/callback" element={<Callback />} />
          <Route
            path="/console"
            element={
              <Console onLogin={() => (window.location.href = AUTH_LOGIN_URL)} />
            }
          >
            <Route index element={<Dashboard />} />
            <Route path="organizations" element={<Organizations />} />
            <Route path="organizations/:oid/projects" element={<Projects />} />
          </Route>
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </BrowserRouter>
    </>
  );
}

export default App;
