import { useCallback, useEffect, useState } from "react";
import apiFetch from "../api/client";

export default function useFetch(url, options) {
  const [data, setData] = useState(null);
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(true);
  const [trigger, setTrigger] = useState(0);

  const fetchData = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);

      const response = await apiFetch(url, options);
      const result = await response.json();

      if (!result.success) {
        throw new Error(result.message);
      }

      setData(result.data);
    } catch (err) {
      setError(err);
    } finally {
      setLoading(false);
    }
  }, [url, trigger]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  const refetch = useCallback(() => setTrigger((t) => t + 1), []);

  return { data, error, loading, refetch };
}
