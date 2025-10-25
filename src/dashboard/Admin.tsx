import "../App.css";
import { useNavigate } from "react-router-dom";
import { useEffect, useState } from "react";

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

  useEffect(() => {
    async function checkAdmin() {
      try {
        const res = await fetch("/account/user", { credentials: "include" });
        if (res.ok) {
          const user: User = await res.json();
          if (!user.is_admin) {
            navigate("/dashboard");
            return;
          }
          setIsAdmin(true);
        } else {
          navigate("/");
        }
      } catch (err) {
        console.error("Failed to fetch user:", err);
        navigate("/");
      } finally {
        setLoading(false);
      }
    }

    checkAdmin();
  }, [navigate]);

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
            ←
          </button>
          <h1 className="text-3xl font-bold" style={{ position: "absolute", top: "1rem" }}>
            Admin Panel
          </h1>

          <div
            style={{
              display: "flex",
              gap: "0.5rem",
              marginTop: "3rem",
              marginBottom: "1rem",
            }}
          >
            <input
              type="text"
              placeholder="Enter user ID"
              value={searchInput}
              onChange={(e) => setSearchInput(e.target.value)}
              onKeyDown={(e) => e.key === "Enter" && handleSearch()}
              style={{
                padding: "0.5rem 1rem",
                fontSize: "1rem",
                borderRadius: "4px",
                border: "1px solid rgba(255, 255, 255, 0.3)",
                backgroundColor: "rgba(0, 0, 0, 0.5)",
                color: "white",
              }}
            />
            <button
              className="nine-slice-button"
              onClick={handleSearch}
              disabled={searching}
              style={{ opacity: searching ? 0.5 : 1 }}
            >
              Search
            </button>
          </div>

          {error && <div style={{ color: "#e74c3c", marginBottom: "1rem" }}>{error}</div>}

          {searchResult && (
            <div
              style={{
                marginTop: "1rem",
                padding: "1rem",
                backgroundColor: "rgba(0, 0, 0, 0.3)",
                borderRadius: "8px",
                maxWidth: "600px",
                maxHeight: "500px",
                overflow: "auto",
              }}
            >
              <div style={{ marginBottom: "1rem", textAlign: "left" }}>
                <div>
                  <strong>Username:</strong> {searchResult.user.username}
                </div>
                <div>
                  <strong>User ID:</strong> {searchResult.user.id}
                </div>
                <div>
                  <strong>Total Views:</strong> {searchResult.user.total_views}
                </div>
                <div>
                  <strong>Total Clicks:</strong> {searchResult.user.total_clicks}
                </div>
                <div>
                  <strong>Admin:</strong> {searchResult.user.is_admin ? "Yes" : "No"}
                </div>
                <div>
                  <strong>Banned:</strong> {searchResult.user.banned ? "Yes" : "No"}
                </div>
              </div>

              <div style={{ marginBottom: "1rem" }}>
                <h3 style={{ marginBottom: "0.5rem" }}>Advertisements ({searchResult.ads.length})</h3>
                {searchResult.ads.length === 0 ? (
                  <div style={{ color: "rgba(255, 255, 255, 0.7)" }}>No advertisements</div>
                ) : (
                  <div
                    style={{
                      display: "flex",
                      flexDirection: "column",
                      gap: "0.5rem",
                    }}
                  >
                    {searchResult.ads.map((ad) => (
                      <div
                        key={ad.ad_id}
                        style={{
                          padding: "0.5rem",
                          backgroundColor: "rgba(255, 255, 255, 0.05)",
                          borderRadius: "4px",
                          fontSize: "0.9rem",
                        }}
                      >
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
                          <strong>Views:</strong> {ad.view_count || 0} |{" "}
                          <strong>Clicks:</strong> {ad.click_count || 0}
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </div>

              <div
                style={{
                  display: "flex",
                  gap: "1rem",
                  justifyContent: "center",
                }}
              >
                <button
                  className="nine-slice-button"
                  onClick={() => handleBanUser(searchResult.user.id)}
                  style={{ fontSize: "0.9rem", padding: "4px 12px" }}
                >
                  Ban User
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
          )}
        </div>
      ) : null}
    </>
  );
}
