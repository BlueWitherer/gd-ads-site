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

import LogoutIcon from "@mui/icons-material/LogoutOutlined";
import AdminPanelSettingsIcon from "@mui/icons-material/AdminPanelSettings";
import NoteAddIcon from '@mui/icons-material/NoteAddOutlined';
import QueryStatsIcon from '@mui/icons-material/QueryStatsOutlined';
import EmojiEventsIcon from '@mui/icons-material/EmojiEventsOutlined';
import AppSettingsAltIcon from '@mui/icons-material/AppSettingsAltOutlined';
import AccountCircleIcon from '@mui/icons-material/AccountCircleOutlined';

export default function Dashboard() {
  const navigate = useNavigate();
  const [selectedView, setSelectedView] = useState<
    "statistics" | "create" | "leaderboard" | "manage" | "account"
  >("statistics");
  const [isBanned, setIsBanned] = useState<boolean>(false);
  const [isAdmin, setIsAdmin] = useState<boolean>(false);

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
                  setIsBanned(userData.banned);
                  setIsAdmin(userData.is_admin);
                });
              };
            })
            .catch(() => console.error("Failed to fetch user status"));
        } else {
          navigate("/");
        };
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
    };

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
            <NoteAddIcon /> Create
          </button>
          <button
            className="nine-slice-button padding-4 mt-4 mb-4"
            onClick={() => setSelectedView("statistics")}
          >
            <QueryStatsIcon /> Statistics
          </button>
          <button
            className="nine-slice-button padding-4 mt-4 mb-4"
            onClick={() => setSelectedView("leaderboard")}
          >
            <EmojiEventsIcon /> Leaderboard
          </button>
          <button
            className="nine-slice-button padding-4 mt-4 mb-4"
            onClick={() => setSelectedView("manage")}
          >
            <AppSettingsAltIcon /> Manage
          </button>
          <button
            className="nine-slice-button padding-4 mt-4 mb-4"
            onClick={() => setSelectedView("account")}
          >
            <AccountCircleIcon /> Account
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

            <div style={{ display: "flex", alignItems: "center", gap: 6 }}>
              {isAdmin && <AdminPanelSettingsIcon />}
              <span className="text-lg">{user !== null ? user.username : "Guest"}</span>
            </div>
          </div>

          <button
            title="Log Out"
            className="nine-slice-button"
            onClick={logout}
          >
            <LogoutIcon style={{ "scale": 1.25 }} />
          </button>
        </div>
        <div className="dashboard-container">{renderContent()}</div>
      </div>
      <CreditsButton />
    </>
  );
}
