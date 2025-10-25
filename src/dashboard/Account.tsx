import "../App.css";
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
  pending?: boolean;
};

export default function Account() {
  const [user, setUser] = useState<User | null>(null);
  const [isBanned, setIsBanned] = useState<boolean>(false);
  const [pendingAds, setPendingAds] = useState<Ad[] | null>(null);
  const [showingPending, setShowingPending] = useState(false);

  useEffect(() => {
    async function fetchUser() {
      try {
        const res = await fetch("/account/user", { credentials: "include" });
        if (res.ok) {
          const data = await res.json();
          setUser(data);
        } else if (res.status === 403) {
          setIsBanned(true);
        }
      } catch (err) {
        console.error("Failed to fetch user:", err);
      }
    }

    fetchUser();
  }, []);

  const handlePendingAds = async () => {
    try {
      const res = await fetch("/ads/pending", {
        method: "GET",
        credentials: "include",
      });
      if (res.ok) {
        const data = await res.json();
        console.log("Pending ads response:", data);
        console.log("Is array?", Array.isArray(data));
        console.log("Length:", data?.length);
        setPendingAds(data);
        setShowingPending(true);
      } else {
        const txt = await res.text();
        console.warn("Pending ads endpoint:", res.status, txt);
        alert("Failed to fetch pending ads: " + txt);
      }
    } catch (err) {
      console.error(err);
      alert("Failed to fetch pending ads.");
    }
  };

  if (isBanned) {
    return (
      <>
        <h1 className="text-2xl font-bold mb-6" style={{ color: "#e74c3c" }}>
          Account Banned
        </h1>
        <p className="text-lg mb-6" style={{ color: "#e74c3c" }}>
          Your account has been banned. You no longer have access to this service.
        </p>
      </>
    );
  }

  if (showingPending) {
    return (
      <>
        <div style={{ display: "flex", alignItems: "center", justifyContent: "center", marginBottom: "1.5rem" }}>
          <h1 className="text-2xl font-bold">Pending Advertisements</h1>
        </div>
        {!pendingAds || pendingAds.length === 0 ? (
          <div style={{ color: "rgba(255, 255, 255, 0.7)", fontSize: "1.1rem", textAlign: "center" }}>
            No pending advertisements
          </div>
        ) : (
          <div
            style={{
              display: "flex",
              flexDirection: "column",
              gap: "1rem",
              maxHeight: "600px",
              overflowY: "auto",
            }}
          >
            {pendingAds.map((ad) => (
              <div
                key={ad.ad_id}
                style={{
                  padding: "1rem",
                  backgroundColor: "rgba(0, 0, 0, 0.3)",
                  borderRadius: "8px",
                  display: "flex",
                  gap: "1rem",
                  alignItems: "center",
                }}
              >
                <a
                  href={ad.image_url}
                  target="_blank"
                  rel="noopener noreferrer"
                  style={{
                    display: "block",
                    flexShrink: 0,
                    textDecoration: "none",
                  }}
                >
                  <img
                    src={ad.image_url}
                    alt={`Ad ${ad.ad_id}`}
                    style={{
                      width: "200px",
                      height: "auto",
                      aspectRatio: "16 / 9",
                      objectFit: "contain",
                      backgroundColor: "rgba(0, 0, 0, 0.5)",
                      cursor: "pointer",
                      transition: "opacity 0.2s ease",
                      borderRadius: "4px",
                    }}
                    onMouseEnter={(e) => (e.currentTarget.style.opacity = "0.8")}
                    onMouseLeave={(e) => (e.currentTarget.style.opacity = "1")}
                  />
                </a>
                <div style={{ flex: 1, display: "flex", flexDirection: "column", gap: "0.5rem" }}>
                  <div>
                    <strong>Ad ID:</strong> {ad.ad_id}
                  </div>
                  <div>
                    <strong>User ID:</strong> {ad.user_id}
                  </div>
                  <div>
                    <strong>Level ID:</strong> {ad.level_id}
                  </div>
                  <div>
                    <strong>Type:</strong> {ad.type}
                  </div>
                  <div>
                    <strong>Created:</strong> {new Date(ad.created_at).toLocaleString()}
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </>
    );
  }

  return (
    <>
      <h1 className="text-2xl font-bold mb-6">My Account</h1>
      <p className="text-lg mb-6">
        Manage your account and view ads pending approval.
      </p>
      <div
        style={{
          display: "flex",
          flexDirection: "column",
          gap: "1em",
          justifyContent: "center",
          alignItems: "center",
        }}
      >
        <div
          style={{
            display: "flex",
            gap: "1em",
            justifyContent: "center",
            alignItems: "center",
          }}
        >
          <button
            className="nine-slice-button"
            style={{ color: "#fff" }}
            onClick={async () => {
              if (
                !confirm(
                  "Are you sure you want to delete your account? This cannot be undone."
                )
              )
                return;
              try {
                const res = await fetch("/account/delete", {
                  method: "POST",
                  credentials: "include",
                });
                if (res.ok) {
                  alert("Account deleted. You will be logged out.");
                  window.location.href = "/";
                } else {
                  const txt = await res.text();
                  console.error("Delete failed:", txt);
                  alert("Failed to delete account.");
                }
              } catch (err) {
                console.error(err);
                alert("Failed to delete account.");
              }
            }}
          >
            Delete Account
          </button>

          {user?.is_admin && (
            <button
              className="nine-slice-button"
              onClick={handlePendingAds}
            >
              Pending Ads
            </button>
          )}
        </div>

        {user?.is_admin && (
          <button
            className="nine-slice-button"
            style={{ color: "#fff" }}
            onClick={() => {
              window.location.href = "/admin";
            }}
          >
            Admin Panel
          </button>
        )}
      </div>
    </>
  );
}
