import "../App.css";
import { useNavigate } from "react-router-dom";
import { useEffect, useState } from "react";
import { copyText } from "../App";
import square03 from '../assets/square03.png';
import blacksquare from '../assets/blacksquare.png';

import ReplyIcon from '@mui/icons-material/ReplyOutlined';
import SearchIcon from '@mui/icons-material/SearchOutlined';
import ContentCopyIcon from "@mui/icons-material/ContentCopyOutlined";
import DoneIcon from "@mui/icons-material/DoneOutlined";
import PersonPinIcon from '@mui/icons-material/PersonPinOutlined';
import BadgeIcon from '@mui/icons-material/BadgeOutlined';
import VisibilityIcon from '@mui/icons-material/VisibilityOutlined';
import MouseIcon from '@mui/icons-material/MouseOutlined';
import AdminPanelSettingsIcon from "@mui/icons-material/AdminPanelSettings";
import GavelIcon from '@mui/icons-material/GavelOutlined';
import AccessTimeIcon from '@mui/icons-material/AccessTimeOutlined';
import DeleteForeverIcon from '@mui/icons-material/DeleteForeverOutlined';

type User = {
  id: string;
  username: string;
  is_admin: boolean;
};

type Ad = {
  ad_id: number;
  user_id: string;
  level_id: string;
  type: number;
  image_url: string;
  created_at: string;
  view_count?: number;
  click_count?: number;
  pending?: boolean;
};

type SearchResult = {
  user: {
    id: string;
    username: string;
    total_views: number;
    total_clicks: number;
    is_admin: boolean;
    banned: boolean;
  };
  ads: Ad[];
};

export default function Admin() {
  const navigate = useNavigate();
  const [isAdmin, setIsAdmin] = useState<boolean | null>(null);
  const [loading, setLoading] = useState(true);
  const [searchInput, setSearchInput] = useState("");
  const [searchResult, setSearchResult] = useState<SearchResult | null>(null);
  const [searching, setSearching] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const [copied, setCopied] = useState(false);

  useEffect(() => {
    async function checkAdmin() {
      try {
        const res = await fetch("/account/user", { credentials: "include" });
        if (res.ok) {
          const user: User = await res.json();
          if (!user.is_admin) {
            navigate("/dashboard");
            return;
          };
          setIsAdmin(true);
        } else {
          navigate("/");
        };
      } catch (err) {
        console.error("Failed to fetch user:", err);
        navigate("/");
      } finally {
        setLoading(false);
      };
    };

    checkAdmin();
  }, [navigate]);

  const handleCopyUserId = async () => {
    await copyText(searchResult?.user.id, setCopied);
  };

  const handleSearch = async () => {
    if (!searchInput.trim()) {
      setError("Please enter a user ID");
      return;
    }

    setSearching(true);
    setError(null);
    setSearchResult(null);

    try {
      const res = await fetch(`/users/${searchInput}`, {
        credentials: "include",
      });

      if (!res.ok) {
        setError("User not found");
        setSearching(false);
        return;
      }

      const data: SearchResult = await res.json();
      setSearchResult(data);
    } catch (err) {
      console.error("Search failed:", err);
      setError("Failed to search user");
    } finally {
      setSearching(false);
    }
  };

  const handleBanUser = async (userId: string) => {
    if (!confirm("Are you sure you want to ban this user?")) return;

    try {
      const res = await fetch(`/ban?id=${userId}`, {
        method: "POST",
        credentials: "include",
      });

      if (res.ok) {
        alert("User banned successfully");
        setSearchResult(null);
        setSearchInput("");
      } else {
        alert("Failed to ban user");
      }
    } catch (err) {
      console.error("Ban failed:", err);
      alert("Failed to ban user");
    }
  };

  const handleUnbanUser = async (userId: string) => {
    if (!confirm("Are you sure you want to unban this user?")) return;

    try {
      const res = await fetch(`/unban?id=${userId}`, {
        method: "POST",
        credentials: "include",
      });

      if (res.ok) {
        alert("User unbanned successfully");
        setSearchResult(null);
        setSearchInput("");
      } else {
        alert("Failed to unban user");
      }
    } catch (err) {
      console.error("Unban failed:", err);
      alert("Failed to unban user");
    }
  };

  const handleDeleteUser = async (userId: string) => {
    if (!confirm("Are you sure you want to delete this user? This cannot be undone."))
      return;

    try {
      const res = await fetch(`/users/${userId}`, {
        method: "DELETE",
        credentials: "include",
      });

      if (res.ok) {
        alert("User deleted successfully");
        setSearchResult(null);
        setSearchInput("");
      } else {
        alert("Failed to delete user");
      }
    } catch (err) {
      console.error("Delete failed:", err);
      alert("Failed to delete user");
    }
  };

  const handleDeleteAd = async (adId: number) => {
    if (!confirm("Are you sure you want to delete this advertisement?")) return;

    try {
      const res = await fetch(`/ads/delete?id=${adId}`, {
        method: "DELETE",
        credentials: "include",
      });

      if (res.ok) {
        alert("Advertisement deleted successfully");
        // Refresh the search to update the ad list
        if (searchInput) {
          handleSearch();
        }
      } else {
        alert("Failed to delete advertisement");
      }
    } catch (err) {
      console.error("Delete ad failed:", err);
      alert("Failed to delete advertisement");
    }
  };

  return (
    <>
      <div id="background-scroll"></div>
      {loading ? (
        <div id="centered-container">
          <div>Loading...</div>
        </div>
      ) : isAdmin ? (
        <div id="centered-container">
          <button
            className="nine-slice-button"
            onClick={() => navigate("/dashboard")}
            style={{
              position: "absolute",
              top: "1rem",
              left: "1rem",
            }}
          >
            <ReplyIcon style={{ "scale": 2.5 }} />
          </button>

          <h1 className="text-3xl font-bold" style={{ marginTop: "1rem", marginBottom: "2rem" }}>
            Admin Panel
          </h1>

          <div
            style={{
              display: "flex",
              gap: "0.5rem",
              marginBottom: "1.5rem",
              width: "auto",
              justifyContent: "center",
              zIndex: 10,
            }}
          >
            <input
              type="text"
              placeholder="Enter user ID or username"
              value={searchInput}
              onChange={(e) => setSearchInput(e.target.value)}
              onKeyDown={(e) => e.key === "Enter" && handleSearch()}
              className="custom-select"
              style={{
                position: "relative",
                zIndex: 5,
                maxWidth: "800px",
                minWidth: "400px",
              }}
            />
            <button
              className="nine-slice-button"
              onClick={handleSearch}
              disabled={searching}
              style={{ opacity: searching ? 0.5 : 1, position: "relative", zIndex: 5, fontSize: "0.9rem", padding: "4px 12px" }}
            >
              <SearchIcon />
            </button>
          </div>

          <div style={{ marginTop: "1rem" }}>
            {error && <div style={{ color: "#e74c3c", marginBottom: "1rem" }}>{error}</div>}

            {searchResult && (
              <div
                style={{
                  marginTop: "1rem",
                  padding: "1.5rem",
                  borderStyle: "solid",
                  borderWidth: "24px",
                  borderImage: `url('${blacksquare}') 24 fill stretch`,
                  maxWidth: "100%",
                  maxHeight: "600px",
                  overflow: "auto",
                  display: "flex",
                  gap: "2rem",
                  backgroundColor: "transparent",
                  boxSizing: "border-box",
                }}
              >
                <div style={{ flex: "0 0 350px" }}>
                  <div
                    style={{
                      marginBottom: "1.5rem",
                      textAlign: "left",
                      padding: "1rem",
                      backgroundColor: "rgba(0, 0, 0, 0.3)",
                      borderRadius: "8px",
                    }}
                  >
                    <div style={{ marginBottom: "0.5rem" }}>
                      <PersonPinIcon /><strong>Username:</strong> {searchResult.user.username}
                    </div>
                    <div style={{ marginBottom: "0.5rem" }}>
                      <BadgeIcon /><strong>User ID:</strong> {searchResult.user.id} <button
                        onClick={handleCopyUserId}
                        style={{ background: "none", border: "none", cursor: "pointer", display: "inline-flex", alignItems: "center", gap: "0.25rem" }}
                        aria-label="Copy user id"
                      >
                        {copied ? (
                          <>
                            <DoneIcon />
                          </>
                        ) : (
                          <ContentCopyIcon />
                        )}
                      </button>
                    </div>
                    <div style={{ marginBottom: "0.5rem" }}>
                      <VisibilityIcon /><strong>Total Views:</strong> {searchResult.user.total_views}
                    </div>
                    <div style={{ marginBottom: "0.5rem" }}>
                      <MouseIcon /><strong>Total Clicks:</strong> {searchResult.user.total_clicks}
                    </div>
                    <div style={{ marginBottom: "0.5rem" }}>
                      <AdminPanelSettingsIcon /><strong>Admin:</strong> {searchResult.user.is_admin ? "Yes" : "No"}
                    </div>
                    <div style={{ marginBottom: "0.5rem" }}>
                      <GavelIcon /><strong>Banned:</strong> {searchResult.user.banned ? "Yes" : "No"}
                    </div>
                  </div>

                  <div
                    style={{
                      display: "flex",
                      gap: "1rem",
                      justifyContent: "center",
                      marginTop: "1.5rem",
                      flexWrap: "wrap",
                    }}
                  >
                    <button
                      className="nine-slice-button"
                      onClick={() =>
                        searchResult.user.banned
                          ? handleUnbanUser(searchResult.user.id)
                          : handleBanUser(searchResult.user.id)
                      }
                      style={{ fontSize: "0.9rem", padding: "4px 12px" }}
                    >
                      <GavelIcon /> {searchResult.user.banned ? "Unban" : "Ban"}
                    </button>
                    <button
                      className="nine-slice-button"
                      onClick={() => handleDeleteUser(searchResult.user.id)}
                      style={{ fontSize: "0.9rem", padding: "4px 12px" }}
                    >
                      Delete User
                    </button>
                  </div>
                </div>

                <div style={{ flex: "1", overflowY: "auto" }}>
                  <h3 style={{ marginBottom: "1rem", marginTop: "0" }}>Advertisements ({searchResult.ads?.length || 0})</h3>
                  {!searchResult.ads || searchResult.ads.length === 0 ? (
                    <div style={{ color: "rgba(255, 255, 255, 0.7)" }}>No advertisements</div>
                  ) : (
                    <div
                      style={{
                        display: "grid",
                        gridTemplateColumns: "repeat(auto-fill, minmax(220px, 1fr))",
                        gap: "1rem",
                      }}
                    >
                      {searchResult.ads.map((ad) => (
                        <div
                          key={ad.ad_id}
                          style={{
                            padding: "0.75rem",
                            backgroundColor: "rgba(0, 0, 0, 0.3)",
                            borderRadius: "8px",
                            fontSize: "0.85rem",
                            display: "flex",
                            flexDirection: "column",
                            gap: "0.5rem",
                            position: "relative",
                            pointerEvents: "auto",
                          }}
                        >
                          {ad.pending && (
                            <div
                              style={{
                                position: "absolute",
                                top: "0.5rem",
                                right: "0.5rem",
                                backgroundColor: "#f39c12",
                                color: "black",
                                padding: "0.25rem 0.5rem",
                                borderRadius: "3px",
                                fontSize: "0.75rem",
                                fontWeight: "bold",
                                zIndex: 10,
                              }}
                            >
                              <AccessTimeIcon /> PENDING
                            </div>
                          )}
                          <a
                            href={ad.image_url}
                            target="_blank"
                            rel="noopener noreferrer"
                            style={{
                              display: "block",
                              width: "100%",
                              textDecoration: "none",
                              pointerEvents: "auto",
                            }}
                          >
                            <img
                              src={ad.image_url}
                              alt={`Ad ${ad.ad_id}`}
                              style={{
                                width: "100%",
                                height: "auto",
                                aspectRatio: "16 / 9",
                                objectFit: "contain",
                                backgroundColor: "rgba(0, 0, 0, 0.5)",
                                cursor: "pointer",
                                transition: "opacity 0.2s ease",
                                display: "block",
                              }}
                              onMouseEnter={(e) => (e.currentTarget.style.opacity = "0.8")}
                              onMouseLeave={(e) => (e.currentTarget.style.opacity = "1")}
                            />
                          </a>
                          <div>
                            <strong>Ad ID:</strong> {ad.ad_id}
                          </div>
                          <div>
                            <strong>Level ID:</strong> {ad.level_id}
                          </div>
                          <div>
                            <strong>Type:</strong> {ad.type}
                          </div>
                          <div>
                            <strong>Views:</strong> {ad.view_count || 0} | <strong>Clicks:</strong> {ad.click_count || 0}
                          </div>
                          <button
                            className="nine-slice-button"
                            onClick={(e) => {
                              e.preventDefault();
                              handleDeleteAd(ad.ad_id);
                            }}
                            style={{
                              fontSize: "0.75rem",
                              padding: "2px 8px",
                              marginTop: "0.25rem",
                              width: "100%",
                            }}
                          >
                            <DeleteForeverIcon /> Delete
                          </button>
                        </div>
                      ))}
                    </div>
                  )}
                </div>
              </div>
            )}
          </div>
        </div >
      ) : null
      }
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
            align-items: center !important;
            gap: 0 !important;
            z-index: 10 !important;
            pointer-events: auto !important;
            background: transparent !important;
            box-sizing: border-box;
            overflow-y: auto;
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
