import "./App.css";
import CreditsButton from "./popup/Credits";
import { useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { FaDiscord } from "react-icons/fa";
import "./Log.mjs";

export async function copyText(
  text: string | undefined,
  setState: React.Dispatch<React.SetStateAction<any | null>>
) {
  try {
    if (text) {
      await navigator.clipboard.writeText(text);
      setState(true);
      setTimeout(() => setState(false), 2500);
    } else {
      console.error("No text provided to copy");
    }
  } catch (err) {
    console.error("Copy failed:", err);
  }
}

export default function App() {
  const navigate = useNavigate();

  useEffect(() => {
    fetch("/session", { credentials: "include" })
      .then((res) => (res.ok ? res.json() : null))
      .then((data) => {
        if (data?.username && data?.id) navigate("/dashboard");
      })
      .catch(() => {
        console.error("User unauthorized");
      });
  }, [navigate]);

  const handleLogin = () => {
    window.location.href = "/login";
  };

  return (
    <>
      <div id="background-scroll"></div>
      <div id="centered-container">
        <div id="login-section">
          <h1 style={{ marginBottom: "2rem", color: "white" }}>GD Advertisement Manager</h1>
          <h2>
            Welcome to the GD Advertisement Manager! Manage all your Geometry
            Dash Advertisements here!
          </h2>
          <h2 style={{ marginBottom: "2rem", color: "white" }}>Login using your Discord Account to get started!</h2>
          <button
            className="nine-slice-button login-button"
            onClick={handleLogin}
            aria-label="Login with Discord"
          >
            <FaDiscord size={25} aria-hidden="true" />
            <span>Login with Discord</span>
          </button>
          <button
            className="nine-slice-button login-button"
            onClick={() =>
              window.open(
                "https://geode-sdk.org/mods/arcticwoof.player_advertisements",
                "_blank"
              )
            }
            aria-label="Install Geode Mod"
            style={{ marginTop: "1rem" }}
          >
            <span>Install Geode Mod</span>
          </button>
        </div>
      </div>
      <CreditsButton />
    </>
  );
}
