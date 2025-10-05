import './App.css';
import CreditsButton from './Credits';
import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import './Log.mjs';

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
        <button className="nine-slice-button" onClick={handleLogin}>
          Login
        </button>
      </div>
      <CreditsButton />
    </>
  );
}
