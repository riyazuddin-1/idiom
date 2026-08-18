import { useEffect, useState } from "react";
import { useSearchParams } from "react-router-dom";
import { authenticate } from "../api/auth";

const Callback = () => {
  const [searchParams] = useSearchParams();
  const [status, setStatus] = useState("pending");

  useEffect(() => {
    const code = searchParams.get("code");
    if (!code) {
      setStatus("failure");
      return;
    }

    authenticate(code)
      .then(() => setStatus("success"))
      .catch(() => setStatus("failure"));
  }, []);

  useEffect(() => {
    if (status === "failure") {
      window.location.href = "/";
    } else if (status === "success") {
      window.location.href = "/console";
    }
  }, [status]);

  return (
    <div style={{ textAlign: "center", padding: "4rem 1rem" }}>
      {status === "pending" && <p>Authenticating...</p>}
      {status === "failure" && <p>Failed to authenticate. Redirecting...</p>}
      {status === "success" && <p>Authenticated. Redirecting to console...</p>}
    </div>
  );
};

export default Callback;
