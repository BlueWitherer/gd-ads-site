import "../page/Login.css";
import "./Admin.css";
import { useNavigate } from "react-router-dom";
import { useEffect, useState } from "react";
import { copyText } from "./Login";

import ReplyIcon from "@mui/icons-material/ReplyOutlined";
import SearchIcon from "@mui/icons-material/SearchOutlined";
import ContentCopyIcon from "@mui/icons-material/ContentCopyOutlined";
import DoneIcon from "@mui/icons-material/DoneOutlined";
// import PersonPinIcon from "@mui/icons-material/PersonPinOutlined";
import BadgeIcon from "@mui/icons-material/BadgeOutlined";
import VisibilityIcon from "@mui/icons-material/VisibilityOutlined";
import MouseIcon from "@mui/icons-material/MouseOutlined";
import AdminPanelSettingsIcon from "@mui/icons-material/AdminPanelSettings";
import BuildIcon from "@mui/icons-material/Build";
import VerifiedIcon from "@mui/icons-material/Verified";
import GavelIcon from "@mui/icons-material/GavelOutlined";
import AccessTimeIcon from "@mui/icons-material/AccessTimeOutlined";
import DeleteForeverIcon from "@mui/icons-material/DeleteForeverOutlined";
import Avatar from "@mui/material/Avatar";

type User = {
  id: string;
  username: string;
  avatar_url: string;
  total_views: number;
  total_clicks: number;
  is_admin: boolean;
  is_staff: boolean;
  verified: boolean;
  banned: boolean;
};

type Ad = {
  ad_id: number;
  user_id: string;
  level_id: string;
  type: number;
  image_url: string;
  created_at: string;
  views?: number;
  clicks?: number;
  pending?: boolean;
};

type SearchResult = {
  user: User;
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
  const [activeTab, setActiveTab] = useState<"info" | "ads">("info");

  const [copied, setCopied] = useState(false);

  useEffect(() => {
    async function checkAdmin() {
      try {
        const res = await fetch("/account/me", { credentials: "include" });
        if (res.ok) {
          const user: User = await res.json();
          if (!user.is_admin) {
            navigate("/dashboard");
            return;
          };

          setIsAdmin(true);
        } else {
          navigate("/");
        }
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
    if (
      !confirm(
        "Are you sure you want to delete this user? This cannot be undone."
      )
    )
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
    if (!confirm("Are you sure you want to delete this advertisement?")) {
      return;
    }

    try {
      const res = await fetch(`/ads/delete?id=${adId}`, {
        method: "DELETE",
        credentials: "include",
      });

      if (res.ok) {
        alert("Advertisement deleted successfully");
        if (searchResult) {
          setSearchResult({
            ...searchResult,
            ads: searchResult.ads.filter((ad) => ad.ad_id !== adId),
          });
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
            className="nine-slice-button back-button"
            onClick={() => navigate("/dashboard")}
          >
            <ReplyIcon className="back-icon" />
          </button>

          <div className="search-container">
            <input
              type="text"
              placeholder="Enter user ID or username"
              value={searchInput}
              onChange={(e) => setSearchInput(e.target.value)}
              onKeyDown={(e) => e.key === "Enter" && handleSearch()}
              className="custom-select"
            />
            <button
              className="nine-slice-button"
              onClick={handleSearch}
              disabled={searching}
            >
              <SearchIcon />
            </button>
          </div>

          {searchResult && (
            <div className="mobile-tabs">
              <button
                className={`mobile-tab ${activeTab === "info" ? "active" : ""}`}
                onClick={() => setActiveTab("info")}
              >
                User Info
              </button>
              <button
                className={`mobile-tab ${activeTab === "ads" ? "active" : ""}`}
                onClick={() => setActiveTab("ads")}
              >
                Advertisements
              </button>
            </div>
          )}

          <div className="results-container" data-active-tab={activeTab}>
            {error && <div className="error-message">{error}</div>}

            {!error && !searchResult && (
              <div className="empty-state">
                Search for a user to get started
              </div>
            )}

            {searchResult && (
              <>
                <div
                  style={{ flex: "0 0 350px", width: "100%" }}
                  className="user-info-section"
                >
                  <div className="user-info-box">
                    <div className="user-info-item user-header">
                      <Avatar alt={searchResult.user.username} src={searchResult.user.avatar_url} />
                      <span className="username">{searchResult.user.username}</span>
                    </div>
                    
                    <div className="user-info-item roles-section">
                      <strong>Roles:</strong>
                      <div className="roles-list">
                        {searchResult.user.is_admin && (
                          <div className="role-item">
                            <AdminPanelSettingsIcon titleAccess="Administrator" />
                            <span>Administrator</span>
                          </div>
                        )}
                        {searchResult.user.is_staff && (
                          <div className="role-item">
                            <BuildIcon titleAccess="Staff" />
                            <span>Staff</span>
                          </div>
                        )}
                        {searchResult.user.verified && (
                          <div className="role-item">
                            <VerifiedIcon titleAccess="Verified" />
                            <span>Verified</span>
                          </div>
                        )}
                        {!searchResult.user.is_admin && !searchResult.user.is_staff && !searchResult.user.verified && (
                          <span className="no-roles">None</span>
                        )}
                      </div>
                    </div>
                    
                    <div className="user-info-item">
                      <BadgeIcon />
                      <strong>UID:</strong> {searchResult.user.id}
                      <button
                        onClick={handleCopyUserId}
                        className="copy-button"
                        aria-label="Copy user ID"
                      >
                        {copied ? <DoneIcon /> : <ContentCopyIcon />}
                      </button>
                    </div>
                    
                    <div className="user-info-item">
                      <VisibilityIcon />
                      <span className="info-label"><strong>Total Views:</strong></span>
                      <span className="info-value">{searchResult.user.total_views}</span>
                    </div>
                    
                    <div className="user-info-item">
                      <MouseIcon />
                      <span className="info-label"><strong>Total Clicks:</strong></span>
                      <span className="info-value">{searchResult.user.total_clicks}</span>
                    </div>
                    
                    <div className="user-info-item">
                      <GavelIcon />
                      <span className="info-label"><strong>Banned:</strong></span>
                      <span className="info-value">{searchResult.user.banned ? "Yes" : "No"}</span>
                    </div>
                  </div>

                  <div className="user-actions">
                    <button
                      className="nine-slice-button action-button"
                      onClick={() =>
                        searchResult.user.banned
                          ? handleUnbanUser(searchResult.user.id)
                          : handleBanUser(searchResult.user.id)
                      }
                    >
                      <GavelIcon /> {searchResult.user.banned ? "Unban" : "Ban"}
                    </button>
                    <button
                      className="nine-slice-button action-button"
                      onClick={() => handleDeleteUser(searchResult.user.id)}
                    >
                      Delete User
                    </button>
                  </div>
                </div>

                <div
                  style={{ flex: "1", width: "100%" }}
                  className="advertisements-section"
                >
                  <h3 className="ads-section-title">
                    Advertisements ({searchResult.ads?.length || 0})
                  </h3>
                  {!searchResult.ads || searchResult.ads.length === 0 ? (
                    <div className="no-ads-message">No advertisements</div>
                  ) : (
                    <div className="ads-grid">
                      {searchResult.ads.map((ad) => (
                        <div key={ad.ad_id} className="ad-card">
                          <a
                            href={ad.image_url}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="ad-image-link"
                          >
                            <img
                              src={ad.image_url}
                              alt={`Ad ${ad.ad_id}`}
                              className="ad-image"
                            />
                          </a>
                          {ad.pending && (
                            <div className="pending-badge">
                              <AccessTimeIcon /> PENDING
                            </div>
                          )}
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
                            <strong>Views:</strong> {ad.views || 0} |{" "}
                            <strong>Clicks:</strong> {ad.clicks || 0}
                          </div>
                          <button
                            className="nine-slice-button delete-ad-button"
                            onClick={(e) => {
                              e.preventDefault();
                              handleDeleteAd(ad.ad_id);
                            }}
                          >
                            <DeleteForeverIcon /> Delete
                          </button>
                        </div>
                      ))}
                    </div>
                  )}
                </div>
              </>
            )}
          </div>
        </div>
      ) : null}
    </>
  );
}
