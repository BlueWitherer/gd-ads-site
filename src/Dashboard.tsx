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
import MenuIcon from '@mui/icons-material/MenuOutlined';
import CloseIcon from '@mui/icons-material/CloseOutlined';
import Avatar from "@mui/material/Avatar";
import square03 from './assets/square03.png';
import blacksquare from './assets/blacksquare.png';

export default function Dashboard() {
  const navigate = useNavigate();
  const [selectedView, setSelectedView] = useState<
    "statistics" | "create" | "leaderboard" | "manage" | "account"
  >("statistics");
  const [isBanned, setIsBanned] = useState<boolean>(false);
  const [isAdmin, setIsAdmin] = useState<boolean>(false);
  const [sidebarOpen, setSidebarOpen] = useState<boolean>(false);

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
        {/* Mobile Menu Toggle Button */}
        <button
          className="nine-slice-button"
          onClick={() => setSidebarOpen(!sidebarOpen)}
          style={{
            position: "absolute",
            top: "1rem",
            left: "1rem",
            zIndex: 100,
            display: "none",
          }}
          id="mobile-menu-btn"
        >
          {sidebarOpen ? <CloseIcon /> : <MenuIcon />}
        </button>

        <div
          className="sidebar-container"
          id="sidebar"
          style={{
            position: "absolute",
            display: "flex",
            flexDirection: "column",
          }}
        >
          <button
            className="nine-slice-button padding-4 mt-4 mb-4"
            onClick={() => {
              setSelectedView("create");
              setSidebarOpen(false);
            }}
          >
            <NoteAddIcon /> Create
          </button>
          <button
            className="nine-slice-button padding-4 mt-4 mb-4"
            onClick={() => {
              setSelectedView("statistics");
              setSidebarOpen(false);
            }}
          >
            <QueryStatsIcon /> Statistics
          </button>
          <button
            className="nine-slice-button padding-4 mt-4 mb-4"
            onClick={() => {
              setSelectedView("leaderboard");
              setSidebarOpen(false);
            }}
          >
            <EmojiEventsIcon /> Leaderboard
          </button>
          <button
            className="nine-slice-button padding-4 mt-4 mb-4"
            onClick={() => {
              setSelectedView("manage");
              setSidebarOpen(false);
            }}
          >
            <AppSettingsAltIcon /> Manage
          </button>
          <button
            className="nine-slice-button padding-4 mt-4 mb-4"
            onClick={() => {
              setSelectedView("account");
              setSidebarOpen(false);
            }}
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
            gap: "1rem",
          }}
        >
          <div style={{ display: "flex", alignItems: "center", justifyContent: "center", gap: 15, flexWrap: "wrap", flex: 1 }}>
            {(() => {
              if (!user) return <Avatar alt="Guest" sx={{ "width": 48, "height": 48 }}>G</Avatar>;
              const avatarUrl =
                user.avatar && user.id
                  ? `https://cdn.discordapp.com/avatars/${user.id}/${user.avatar}.png?size=64`
                  : user.discriminator
                    ? `https://cdn.discordapp.com/embed/avatars/${parseInt(user.discriminator || "0", 10) % 5
                    }.png`
                    : null;
              return avatarUrl ? (
                <Avatar alt={user.username} src={avatarUrl} sx={{ "width": 48, "height": 48 }} />
              ) : null;
            })()}

            <div style={{ display: "flex", alignItems: "center", gap: 4 }}>
              {isAdmin && <AdminPanelSettingsIcon />}
              <span style={{ fontSize: "0.9rem" }}>{user !== null ? user.username : "Guest"}</span>
            </div>
          </div>

          <button
            title="Log Out"
            className="nine-slice-button"
            onClick={logout}
            style={{ padding: "2px 4px", fontSize: "10px", flexShrink: 0 }}
          >
            <LogoutIcon style={{ "scale": 1 }} />
          </button>
        </div>
        <div className="dashboard-container">{renderContent()}</div>
      </div>
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
            justify-content: flex-start !important;
            align-items: stretch !important;
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

          #mobile-menu-btn {
            display: block !important;
            position: fixed !important;
            top: 2rem !important;
            left: 2rem !important;
            z-index: 1000 !important;
            padding: 4px 8px !important;
            font-size: 12px !important;
          }

          .sidebar-container {
            position: fixed !important;
            top: 8rem !important;
            left: 50% !important;
            transform: translateX(-50%) !important;
            width: 90% !important;
            min-height: auto !important;
            max-height: calc(100vh - 380px) !important;
            margin: 0 !important;
            border-style: solid !important;
            border-width: 24px !important;
            border-image: url('${blacksquare}') 24 fill stretch !important;
            padding: 2rem !important;
            display: ${sidebarOpen ? "flex !important" : "none !important"};
            flex-direction: column;
            gap: 0.5rem;
            background: transparent;
            z-index: 999;
            overflow-y: auto;
            overflow-x: hidden;
            pointer-events: auto;
            box-sizing: border-box;
          }

          .sidebar-container .nine-slice-button {
            width: 100%;
            margin: 0 !important;
            padding: 6px 12px !important;
            font-size: 14px !important;
            pointer-events: auto;
            white-space: nowrap;
          }

          .user-container {
            position: fixed !important;
            top: 0 !important;
            left: 50% !important;
            transform: translateX(-50%) translateY(25%) !important;
            width: calc(100% - 18rem) !important;
            min-height: auto !important;
            margin: 0 !important;
            padding: 1rem 1.5rem !important;
            border-style: solid !important;
            z-index: 100;
            pointer-events: auto;
            display: flex !important;
            flex-direction: row !important;
            justify-content: space-between !important;
            align-items: center !important;
            gap: 1.5rem;
            background: transparent;
            box-sizing: border-box;
          }

          .user-container > div:first-child {
            width: auto;
            display: flex !important;
            justify-content: flex-start !important;
            align-items: center !important;
            flex-wrap: nowrap;
            gap: 0.75rem;
            flex-direction: row;
            flex: 0;
          }

          .user-container .nine-slice-button {
            transform: none !important;
            padding: 4px 8px !important;
            font-size: 10px !important;
            margin: 0 !important;
            pointer-events: auto;
            flex-shrink: 0;
          }

          .dashboard-container {
            position: fixed !important;
            left: 50% !important;
            top: 10rem !important;
            transform: translateX(-50%) !important;
            width: calc(100% - 8rem) !important;
            height: calc(100vh - 9rem) !important;
            margin: 0 !important;
            border-style: solid !important;
            border-width: 24px !important;
            border-image: url('${blacksquare}') 24 fill stretch !important;
            padding: 1rem !important;
            background: transparent;
            border-radius: 0;
            overflow-y: auto;
            overflow-x: hidden;
            flex-grow: 0;
            pointer-events: auto;
            box-sizing: border-box;
          }
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
      <CreditsButton />
    </>
  );
}
