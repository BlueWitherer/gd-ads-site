import './App.css';
import CreditsButton from './Credits';
import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { FaDiscord } from 'react-icons/fa';
import './Log.mjs';
import square03 from './assets/square03.png';

export async function copyText(text: string | undefined, setState: React.Dispatch<React.SetStateAction<any | null>>) {
  try {
    if (text) {
      await navigator.clipboard.writeText(text);
      setState(true);
      setTimeout(() => setState(false), 2500);
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
      <style>{`
        @media (max-width: 1024px) {
          * {
            box-sizing: border-box;
          }

          body {
            margin: 0;
            padding: 0;
            overflow: hidden;
          }

          #root {
            width: 100vw;
            height: 100vh;
            padding: 0;
            margin: 0;
          }

          #background-scroll {
            z-index: 0 !important;
          }

          #centered-container {
            position: fixed !important;
            top: 0 !important;
            left: 0 !important;
            width: 100vw !important;
            height: 100vh !important;
            transform: none !important;
            max-height: none !important;
            border-style: solid !important;
            border-width: 32px !important;
            border-image: url('${square03}') 32 fill stretch !important;
            padding: 1rem !important;
            margin: 0 !important;
            display: flex !important;
            flex-direction: column !important;
            justify-content: center !important;
            align-items: center !important;
            gap: 0 !important;
            z-index: 10 !important;
            pointer-events: auto !important;
            background: transparent !important;
            box-sizing: border-box;
          }

          #centered-container .nine-slice-button,
          #centered-container button,
          #centered-container a,
          #centered-container input,
          #centered-container select {
            pointer-events: auto !important;
          }

          .sprite-button {
            position: fixed !important;
            top: 2rem !important;
            right: 2rem !important;
            bottom: auto !important;
            z-index: 1001 !important;
          }
        }
      `}</style>
    </>
  );
}
