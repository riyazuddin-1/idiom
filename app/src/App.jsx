import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";
import Callback from "./pages/Callback";
import Console from "./pages/Console";
import Dashboard from "./pages/Dashboard";
import Organizations from "./pages/Organizations";
import Projects from "./pages/Projects";
import "./App.css";

const AUTH_LOGIN_URL =
  import.meta.env.VITE_AUTH_LOGIN_URL ||
  "https://idiom-identity.onrender.com/idiom/login";

function App() {
  return (
    <BrowserRouter>
      <Routes>
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
        <Route path="*" element={<Navigate to="/console" replace />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
