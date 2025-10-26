import "../App.css";
import { useEffect, useState } from "react";
import square02 from "../assets/square02.png";
import { copyText } from "../App";

import ContentCopyIcon from "@mui/icons-material/ContentCopyOutlined";
import DoneIcon from "@mui/icons-material/DoneOutlined";

type User = {
  id: string;
  username: string;
  is_admin: boolean;
  created_at: string;
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
  const [pendingAds, setPendingAds] = useState<Ad[] | null>(null);
  const [showingPending, setShowingPending] = useState(false);

  const [copied, setCopied] = useState(false);

  useEffect(() => {
    async function fetchUser() {
      try {
        const res = await fetch("/account/user", { credentials: "include" });
        if (res.ok) {
          const data = await res.json();
          setUser(data);
        };
      } catch (err) {
        console.error("Failed to fetch user:", err);
      };
    };

    fetchUser();
  }, []);

  const handleCopyUserId = async () => {
    await copyText(user?.id, setCopied);
  };

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
      };
    } catch (err) {
      console.error(err);
      alert("Failed to fetch pending ads.");
    };
  };

  const handleApproveAd = async (adId: number) => {
    try {
      const res = await fetch(`/ads/pending/accept?id=${adId}`, {
        method: "POST",
        credentials: "include",
      });
      if (res.ok) {
        alert("Advertisement approved successfully");
        // Refresh the pending ads list
        handlePendingAds();
      } else {
        alert("Failed to approve advertisement");
      }
    } catch (err) {
      console.error("Approve failed:", err);
      alert("Failed to approve advertisement");
    }
  };

  const handleRejectAd = async (adId: number) => {
    if (!confirm("Are you sure you want to reject and delete this advertisement?")) return;

    try {
      const res = await fetch(`/ads/delete?id=${adId}`, {
        method: "DELETE",
        credentials: "include",
      });
      if (res.ok) {
        alert("Advertisement rejected and deleted successfully");
        // Refresh the pending ads list
        handlePendingAds();
      } else {
        alert("Failed to reject advertisement");
      }
    } catch (err) {
      console.error("Reject failed:", err);
      alert("Failed to reject advertisement");
    }
  };

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
                  borderStyle: "solid",
                  borderWidth: "12px",
                  borderImage: `url(${square02}) 24 fill stretch`,
                  background: "transparent",
                  borderRadius: "0px",
                  display: "flex",
                  flexDirection: "column",
                  gap: "1rem",
                  alignItems: "center",
                }}
              >
                <div style={{ display: "flex", justifyContent: "center", alignItems: "center", flexShrink: 0 }}>
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
                </div>
                <div style={{ display: "flex", flexDirection: "column", gap: "0.5rem", textAlign: "center", width: "100%" }}>
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
                <div style={{ display: "flex", flexDirection: "column", gap: "0.5rem", width: "100%" }}>
                  <button
                    onClick={() => handleApproveAd(ad.ad_id)}
                    style={{
                      fontSize: "0.85rem",
                      padding: "8px 20px",
                      backgroundColor: "#27ae60",
                      color: "#fff",
                      border: "none",
                      borderRadius: "6px",
                      cursor: "pointer",
                      fontWeight: "bold",
                      transition: "background 0.2s",
                      width: "100%",
                    }}
                    onMouseEnter={(e) => (e.currentTarget.style.backgroundColor = "#229954")}
                    onMouseLeave={(e) => (e.currentTarget.style.backgroundColor = "#27ae60")}
                  >
                    Approve
                  </button>
                  <button
                    onClick={() => handleRejectAd(ad.ad_id)}
                    style={{
                      fontSize: "0.85rem",
                      padding: "8px 20px",
                      backgroundColor: "#e74c3c",
                      color: "#fff",
                      border: "none",
                      borderRadius: "6px",
                      cursor: "pointer",
                      fontWeight: "bold",
                      transition: "background 0.2s",
                      width: "100%",
                    }}
                    onMouseEnter={(e) => (e.currentTarget.style.backgroundColor = "#c0392b")}
                    onMouseLeave={(e) => (e.currentTarget.style.backgroundColor = "#e74c3c")}
                  >
                    Reject
                  </button>
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
      <p className="text-lg">
        Manage your account and view ads pending approval.
      </p>
      <p className="text-sm mb-6 text-gray-500">
        If you want to view your stats in-game, copy the User ID into the settings prompted by the popup.
      </p>
      {user && (
        <div
          style={{
            marginBottom: "2rem",
            padding: "1rem",
            backgroundColor: "rgba(255, 255, 255, 0.05)",
            borderRadius: "8px",
            border: "1px solid rgba(255, 255, 255, 0.1)",
          }}
        >
          <div style={{ display: "flex", flexDirection: "column", alignItems: "center", gap: "0.5rem" }}>
            <div>
              <strong>User ID: </strong> {user.id} <button onClick={handleCopyUserId} style={{ background: 'none', border: 'none', cursor: 'pointer' }}>
                {copied ? (
                  <>
                    <DoneIcon />
                  </>
                ) : (
                  <ContentCopyIcon />
                )}
              </button>
            </div>
            <div>
              <strong>Account Created: </strong> {new Date(user.created_at).toLocaleString()}
            </div>
          </div>
        </div>
      )}

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
