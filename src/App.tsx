import './App.css';
import CreditsButton from './Credits';
import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { FaDiscord } from 'react-icons/fa';
import './Log.mjs';

export async function copyText(text: string | undefined, setState: React.Dispatch<React.SetStateAction<any | null>>) {
  try {
    if (text) {
      await navigator.clipboard.writeText(text);
      setState(true);
      setTimeout(() => setState(false), 2000); // Reset after 2s
    } else {
      console.error("No text provided to copy");
    };
  } catch (err) {
    console.error("Copy failed:", err);
  };
};

export default function App() {
  const navigate = useNavigate();

  useEffect(() => {
    fetch("/session", { credentials: "include" })
      .then((res) => res.ok ? res.json() : null)
      .then((data) => {
        if (data?.username && data?.id) navigate("/dashboard");
      })
      .catch(() => {
        console.error("User unauthorized");
      });
  }, [navigate]);

  const handleLogin = () => {
    window.location.href = '/login';
  };

  return (
    <>
      <div id="background-scroll"></div>
      <div id="centered-container">
        <h1 className="text-3xl font-bold padding-4 mt-4 mb-8">
          Advertisement Manager
        </h1>
        <button
          className="nine-slice-button"
          onClick={handleLogin}
          style={{ display: "inline-flex", alignItems: "center", gap: "0.5rem", padding: "0.5rem 1rem" }}
          aria-label="Login with Discord"
        >
          <FaDiscord size={25} aria-hidden="true" />
          <span>Login with Discord</span>
        </button>
      </div>
      <CreditsButton />
    </>
  );
}
