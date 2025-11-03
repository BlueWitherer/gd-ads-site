import { useEffect, useCallback } from "react";
import { useNavigate } from "react-router-dom";

/**
 * Custom hook to validate user session and redirect to login if invalid
 * @param enabled - Whether session validation is enabled (default: true)
 */
export function useSessionValidation(enabled: boolean = true) {
  const navigate = useNavigate();

  const validateSession = useCallback(async (): Promise<boolean> => {
    if (!enabled) return true;

    try {
      const res = await fetch("/session", { credentials: "include" });
      if (!res.ok) {
        console.warn("Session invalid, redirecting to login...");
        navigate("/");
        return false;
      }
      return true;
    } catch (error) {
      console.error("Error validating session:", error);
      navigate("/");
      return false;
    }
  }, [navigate, enabled]);

  return { validateSession };
}

/**
 * Custom hook to automatically validate session on component mount and interactions
 * @param onInvalidSession - Optional callback when session becomes invalid
 */
export function useAutoSessionValidation(onInvalidSession?: () => void) {
  const navigate = useNavigate();

  const checkSession = useCallback(async () => {
    try {
      const res = await fetch("/session", { credentials: "include" });
      if (!res.ok) {
        console.warn("Session expired or invalid, redirecting to login...");
        if (onInvalidSession) {
          onInvalidSession();
        }
        navigate("/");
      }
    } catch (error) {
      console.error("Error checking session:", error);
      navigate("/");
    }
  }, [navigate, onInvalidSession]);

  useEffect(() => {
    checkSession();

    // every 60 seconds
    const interval = setInterval(checkSession, 60000);
    let debounceTimer: number | null = null;
    let lastCheckTime = Date.now();

    const handleInteraction = () => {
      const now = Date.now();
      if (now - lastCheckTime < 5000) {
        return;
      }

      if (debounceTimer) {
        clearTimeout(debounceTimer);
      }

      debounceTimer = setTimeout(() => {
        lastCheckTime = Date.now();
        checkSession();
      }, 1000); // Wait 1 second after last interaction before checking
    };

    window.addEventListener("click", handleInteraction);
    window.addEventListener("keydown", handleInteraction);

    return () => {
      clearInterval(interval);
      if (debounceTimer) {
        clearTimeout(debounceTimer);
      }
      window.removeEventListener("click", handleInteraction);
      window.removeEventListener("keydown", handleInteraction);
    };
  }, [checkSession]);

  return { checkSession };
}

/**
 * Wrapper for fetch that automatically validates session before making requests
 * @param url - The URL to fetch
 * @param options - Fetch options
 * @param navigate - React Router navigate function
 * @returns Promise with the fetch response
 */
export async function fetchWithSessionCheck(
  url: string,
  options: RequestInit,
  navigate: (path: string) => void
): Promise<Response> {
  try {
    const response = await fetch(url, options);
    if (response.status === 401 || response.status === 403) {
      const sessionCheck = await fetch("/session", { credentials: "include" });
      if (!sessionCheck.ok) {
        console.warn(
          "Session invalid during API call, redirecting to login..."
        );
        navigate("/");
        throw new Error("Session invalid");
      }
    }

    return response;
  } catch (error) {
    console.error("Error in fetchWithSessionCheck:", error);
    throw error;
  }
}
