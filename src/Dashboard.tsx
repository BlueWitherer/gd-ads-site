import "./App.css";
import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import Statistics from "./dashboard/Statistics";
import Create from "./dashboard/Create";
import Leaderboard from "./dashboard/Leaderboard";
import Manage from "./dashboard/Manage";
import Account from "./dashboard/Account";
import CreditsButton from "./Credits";
import "./Log.mjs";

export default function Dashboard() {
  const navigate = useNavigate();
  const [selectedView, setSelectedView] = useState<
    "statistics" | "create" | "leaderboard" | "manage" | "account"
  >("statistics");
  const [isBanned, setIsBanned] = useState<boolean>(false);

  const [user, setUser] = useState<{
    id: string;
    username: string;
    avatar?: string | null;
    discriminator?: string | null;
  } | null>(null);

  useEffect(() => {
    fetch("/session", { credentials: "include" })
      .then((res) => (res.ok ? res.json() : null))
      .then((data) => {
        if (data?.username && data?.id) {
          setUser({
            id: data.id,
            username: data.username,
            avatar: data.avatar ?? null,
            discriminator: data.discriminator ?? null,
          });

          // Check if user is banned
          fetch("/account/user", { credentials: "include" })
            .then((res) => {
              if (res.status === 403) {
                setIsBanned(true);
              } else if (res.ok) {
                res.json().then((userData) => {
                  if (userData.banned) {
                    setIsBanned(true);
                  }
                });
              }
            })
            .catch(() => console.error("Failed to fetch user status"));
        } else {
          navigate("/");
        }
      })
      .catch(() => navigate("/"))
      .finally(() => console.log("User authorized"));
  }, [navigate]);

  const logout = () => {
    fetch("/logout", {
      method: "POST",
      credentials: "include",
    })
      .then(() => {
        navigate("/");
      })
      .catch(() => {
        navigate("/");
      });
  };

  const renderContent = () => {
    if (isBanned) {
      return (
        <>
          <h1 className="text-2xl font-bold mb-6" style={{ color: "#e74c3c" }}>
            Account Banned
          </h1>
          <p className="text-lg" style={{ color: "#e74c3c" }}>
            Your account has been banned. You no longer have access to this service.
          </p>
        </>
      );
    }

    switch (selectedView) {
      case "statistics":
        return <Statistics />;
      case "create":
        return <Create />;
      case "leaderboard":
        return <Leaderboard />;
      case "manage":
        return <Manage />;
      case "account":
        return <Account />;

      default:
        return null;
    }
  };

  return (
    <>
      <div id="background-scroll"></div>
      <div id="centered-container">
        <div className="sidebar-container">
          <button
            className="nine-slice-button padding-4 mt-4 mb-4"
            onClick={() => setSelectedView("create")}
          >
            Create
          </button>
          <button
            className="nine-slice-button padding-4 mt-4 mb-4"
            onClick={() => setSelectedView("statistics")}
          >
            Statistics
          </button>
          <button
            className="nine-slice-button padding-4 mt-4 mb-4"
            onClick={() => setSelectedView("leaderboard")}
          >
            Leaderboard
          </button>
          <button
            className="nine-slice-button padding-4 mt-4 mb-4"
            onClick={() => setSelectedView("manage")}
          >
            Manage
          </button>
          <button
            className="nine-slice-button padding-4 mt-4 mb-4"
            onClick={() => setSelectedView("account")}
          >
            Account
          </button>
        </div>
        <div
          className="user-container"
          style={{
            display: "flex",
            alignItems: "center",
            justifyContent: "space-between",
          }}
        >
          <div style={{ display: "flex", alignItems: "center", gap: 24 }}>
            {(() => {
              if (!user) return null;
              const avatarUrl =
                user.avatar && user.id
                  ? `https://cdn.discordapp.com/avatars/${user.id}/${user.avatar}.png?size=64`
                  : user.discriminator
                    ? `https://cdn.discordapp.com/embed/avatars/${parseInt(user.discriminator || "0", 10) % 5
                    }.png`
                    : null;
              return avatarUrl ? (
                <img
                  src={avatarUrl}
                  alt="avatar"
                  style={{ width: 64, height: 64, borderRadius: 9999 }}
                />
              ) : null;
            })()}

            <p className="text-lg">{user !== null ? user.username : "Guest"}</p>
          </div>

          <button
            title="Log Out"
            className="nine-slice-button"
            onClick={logout}
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="32"
              height="32"
              fill="currentColor"
              className="bi bi-box-arrow-left"
              viewBox="0 0 16 16"
            >
              <path
                fillRule="evenodd"
                d="M6 12.5a.5.5 0 0 0 .5.5h8a.5.5 0 0 0 .5-.5v-9a.5.5 0 0 0-.5-.5h-8a.5.5 0 0 0-.5.5v2a.5.5 0 0 1-1 0v-2A1.5 1.5 0 0 1 6.5 2h8A1.5 1.5 0 0 1 16 3.5v9a1.5 1.5 0 0 1-1.5 1.5h-8A1.5 1.5 0 0 1 5 12.5v-2a.5.5 0 0 1 1 0z"
              />
              <path
                fillRule="evenodd"
                d="M.146 8.354a.5.5 0 0 1 0-.708l3-3a.5.5 0 1 1 .708.708L1.707 7.5H10.5a.5.5 0 0 1 0 1H1.707l2.147 2.146a.5.5 0 0 1-.708.708z"
              />
            </svg>
          </button>
        </div>
        <div className="dashboard-container">{renderContent()}</div>
      </div>
      <CreditsButton />
    </>
  );
}
