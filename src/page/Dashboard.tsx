import "./Login.css";
import "./Dashboard.css";
import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import Statistics from "../dashboard/Statistics";
import Create from "../dashboard/Create";
import Leaderboard from "../dashboard/Leaderboard";
import Manage from "../dashboard/Manage";
import Account from "../dashboard/Account";
import CreditsButton from "../popup/Credits";
import { useAutoSessionValidation } from "../utils/useSessionValidation";
import "../misc/Log.mjs";

import LogoutIcon from "@mui/icons-material/LogoutOutlined";
import AdminPanelSettingsIcon from "@mui/icons-material/AdminPanelSettings";
import BuildIcon from "@mui/icons-material/Build";
import VerifiedIcon from "@mui/icons-material/Verified";
import NoteAddIcon from "@mui/icons-material/NoteAddOutlined";
import QueryStatsIcon from "@mui/icons-material/QueryStatsOutlined";
import EmojiEventsIcon from "@mui/icons-material/EmojiEventsOutlined";
import AppSettingsAltIcon from "@mui/icons-material/AppSettingsAltOutlined";
import AccountCircleIcon from "@mui/icons-material/AccountCircleOutlined";
import MenuIcon from "@mui/icons-material/MenuOutlined";
import CloseIcon from "@mui/icons-material/CloseOutlined";
import Avatar from "@mui/material/Avatar";

export default function Dashboard() {
  const navigate = useNavigate();
  const [selectedView, setSelectedView] = useState<
    "statistics" | "create" | "leaderboard" | "manage" | "account"
  >("statistics");
  const [isBanned, setIsBanned] = useState<boolean>(false);
  const [isAdmin, setIsAdmin] = useState<boolean>(false);
  const [isStaff, setIsStaff] = useState<boolean>(false);
  const [verified, setVerified] = useState<boolean>(false);
  const [sidebarOpen, setSidebarOpen] = useState<boolean>(false);
  const [sidebarClosing, setSidebarClosing] = useState<boolean>(false);

  const [user, setUser] = useState<{
    id: string;
    username: string;
    avatar?: string | null;
    discriminator?: string | null;
  } | null>(null);

  useAutoSessionValidation();

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
          fetch("/account/me", { credentials: "include" })
            .then((res) => {
              if (res.status === 403) {
                setIsBanned(true);
              } else if (res.ok) {
                res.json().then((userData) => {
                  setIsAdmin(userData.is_admin);
                  setIsStaff(userData.is_staff);
                  setVerified(userData.verified);
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

  const closeSidebarWithAnimation = (callback?: () => void) => {
    setSidebarClosing(true);
    setTimeout(() => {
      setSidebarOpen(false);
      setSidebarClosing(false);
      callback?.();
    }, 300);
  };

  const renderContent = () => {
    if (isBanned) {
      return (
        <>
          <h1 className="banned-title">Account Banned</h1>
          <p className="banned-text">
            Your account has been banned. You no longer have access to this
            service.
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
        {/* Mobile Menu Toggle Button */}
        <button
          className="nine-slice-button mobile-menu-btn"
          onClick={() => {
            if (sidebarOpen) {
              setSidebarClosing(true);
              setTimeout(() => {
                setSidebarOpen(false);
                setSidebarClosing(false);
              }, 300);
            } else {
              setSidebarOpen(true);
            }
          }}
          id="mobile-menu-btn"
        >
          {sidebarOpen ? <CloseIcon /> : <MenuIcon />}
        </button>

        <div
          className={`sidebar-container sidebar-wrapper ${
            sidebarClosing ? "closing" : sidebarOpen ? "open" : ""
          }`}
          id="sidebar"
        >
          <button
            className="nine-slice-button padding-4 mt-4 mb-4"
            onClick={() => {
              setSelectedView("create");
              closeSidebarWithAnimation();
            }}
          >
            <NoteAddIcon /> Create
          </button>
          <button
            className="nine-slice-button padding-4 mt-4 mb-4"
            onClick={() => {
              setSelectedView("statistics");
              closeSidebarWithAnimation();
            }}
          >
            <QueryStatsIcon /> Statistics
          </button>
          <button
            className="nine-slice-button padding-4 mt-4 mb-4"
            onClick={() => {
              setSelectedView("leaderboard");
              closeSidebarWithAnimation();
            }}
          >
            <EmojiEventsIcon /> Leaderboard
          </button>
          <button
            className="nine-slice-button padding-4 mt-4 mb-4"
            onClick={() => {
              setSelectedView("manage");
              closeSidebarWithAnimation();
            }}
          >
            <AppSettingsAltIcon /> Manage
          </button>
          <button
            className="nine-slice-button padding-4 mt-4 mb-4"
            onClick={() => {
              setSelectedView("account");
              closeSidebarWithAnimation();
            }}
          >
            <AccountCircleIcon /> Account
          </button>
        </div>
        <div className="user-container user-container-wrapper">
          <div className="user-info-container">
            {(() => {
              if (!user)
                return (
                  <Avatar alt="Guest" sx={{ width: 48, height: 48 }}>
                    G
                  </Avatar>
                );
              const avatarUrl =
                user.avatar && user.id
                  ? `https://cdn.discordapp.com/avatars/${user.id}/${user.avatar}.png?size=64`
                  : user.discriminator
                  ? `https://cdn.discordapp.com/embed/avatars/${
                      parseInt(user.discriminator || "0", 10) % 5
                    }.png`
                  : null;
              return avatarUrl ? (
                <Avatar alt={user.username} src={avatarUrl} />
              ) : null;
            })()}

            <div style={{ display: "flex", alignItems: "center", gap: 4 }}>
              <span style={{ fontSize: "1rem" }}>
                {user !== null ? user.username : "Guest"}
              </span>
              {(isAdmin && (
                <AdminPanelSettingsIcon titleAccess="Administrator" />
              )) ||
                (isStaff && <BuildIcon titleAccess="Staff" />) ||
                (verified && <VerifiedIcon titleAccess="Verified" />)}
            </div>
          </div>

          <button
            title="Log Out"
            className="nine-slice-button logout-button"
            onClick={logout}
          >
            <LogoutIcon className="logout-icon" />
          </button>
        </div>
        <div className="dashboard-container">{renderContent()}</div>
      </div>
      <CreditsButton />
    </>
  );
}
