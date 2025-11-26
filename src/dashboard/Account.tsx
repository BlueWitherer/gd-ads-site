import "../page/Login.css";
import "./Dashboard.css";
import { useEffect, useState } from "react";
import square02 from "../assets/square02.png";
import { copyText } from "../page/Login";

import ContentCopyIcon from "@mui/icons-material/ContentCopyOutlined";
import DoneIcon from "@mui/icons-material/DoneOutlined";

type User = {
  id: string;
  username: string;
  is_admin: boolean;
  is_staff: boolean;
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

type Report = {
  id: number;
  ad: Ad;
  account_id: number;
  description: string;
  created_at: string;
};

export default function Account() {
  const [user, setUser] = useState<User | null>(null);
  const [pendingAds, setPendingAds] = useState<Ad[] | null>(null);
  const [showingPending, setShowingPending] = useState(false);
  const [pendingCount, setPendingCount] = useState<number>(0);

  const [reportedAds, setReportedAds] = useState<Report[] | null>(null);
  const [showingReported, setShowingReported] = useState(false);
  const [reportedCount, setReportedCount] = useState<number>(0);

  const [copied, setCopied] = useState(false);

  useEffect(() => {
    async function fetchUser() {
      try {
        const res = await fetch("/account/me", { credentials: "include" });
        if (res.ok) {
          const data = await res.json();
          setUser(data);
        }
      } catch (err) {
        console.error("Failed to fetch user:", err);
      }
    }

    async function fetchPendingCount() {
      try {
        const res = await fetch("/ads/pending", {
          method: "GET",
          credentials: "include",
        });
        if (res.ok) {
          const data = await res.json();
          setPendingCount(Array.isArray(data) ? data.length : 0);
        }
      } catch (err) {
        console.error("Failed to fetch pending ads count:", err);
      }
    }

    async function fetchReportedCount() {
      try {
        const res = await fetch("/ads/reports", {
          method: "GET",
          credentials: "include",
        });
        if (res.ok) {
          const data = await res.json();
          setReportedCount(Array.isArray(data) ? data.length : 0);
        }
      } catch (err) {
        console.error("Failed to fetch reported ads count:", err);
      }
    }

    fetchUser();
    fetchPendingCount();
    fetchReportedCount();
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
        setPendingCount(Array.isArray(data) ? data.length : 0);
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

  const handleApproveAd = async (adId: number) => {
    try {
      const res = await fetch(`/ads/pending/accept?id=${adId}`, {
        method: "POST",
        credentials: "include",
      });
      if (res.ok) {
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
    if (
      !confirm("Are you sure you want to reject and delete this advertisement?")
    )
      return;

    try {
      const res = await fetch(`/ads/delete?id=${adId}`, {
        method: "DELETE",
        credentials: "include",
      });
      if (res.ok) {
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

  const handleReportedAds = async () => {
    try {
      const res = await fetch("/ads/reports", {
        method: "GET",
        credentials: "include",
      });
      if (res.ok) {
        const data = await res.json();
        setReportedAds(data);
        setReportedCount(Array.isArray(data) ? data.length : 0);
        setShowingReported(true);
      } else {
        const txt = await res.text();
        alert("Failed to fetch reported ads: " + txt);
      }
    } catch (err) {
      console.error(err);
      alert("Failed to fetch reported ads.");
    }
  };

  const handleDeleteReportedAd = async (reportId: number) => {
    if (
      !confirm(
        "Are you sure you want to delete this advertisement? This action cannot be undone."
      )
    )
      return;

    try {
      const res = await fetch(`/ads/reports/action?id=${reportId}&action=1`, {
        method: "POST",
        credentials: "include",
      });
      if (res.ok) {
        handleReportedAds();
      } else {
        alert("Failed to delete advertisement");
      }
    } catch (err) {
      console.error("Delete failed:", err);
      alert("Failed to delete advertisement");
    }
  };

  const handleRejectReport = async (reportId: number) => {
    try {
      const res = await fetch(`/ads/reports/reject?id=${reportId}`, {
        method: "POST",
        credentials: "include",
      });
      if (res.ok) {
        handleReportedAds();
      } else {
        alert("Failed to reject report");
      }
    } catch (err) {
      console.error("Reject failed:", err);
      alert("Failed to reject report");
    }
  };

  if (showingReported) {
    return (
      <>
        <div className="pending-ads-header">
          <h2>Reported Advertisements</h2>
        </div>
        <div className="pending-ads-actions">
          <button
            className="nine-slice-button small pending-ads-action-button"
            onClick={handleReportedAds}
          >
            Refresh
          </button>
          <button
            className="nine-slice-button small pending-ads-action-button"
            onClick={() => setShowingReported(false)}
          >
            Back
          </button>
        </div>
        {!reportedAds || reportedAds.length === 0 ? (
          <div className="pending-ads-empty">No reported advertisements</div>
        ) : (
          <div className="pending-ads-list">
            {reportedAds.map((report) => (
              <div
                key={report.id}
                className="pending-ad-card"
                style={{
                  borderImage: `url(${square02}) 24 fill stretch`,
                }}
              >
                <div className="pending-ad-image-wrapper">
                  <a
                    href={report.ad.image_url}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="pending-ad-image-link"
                  >
                    <img
                      src={report.ad.image_url}
                      alt={`Ad ${report.ad.ad_id}`}
                      className="pending-ad-image"
                    />
                  </a>
                </div>
                <div className="pending-ad-info">
                  <div>
                    <strong>Report ID:</strong> {report.id}
                  </div>
                  <div>
                    <strong>Reason:</strong> {report.description}
                  </div>
                  <div>
                    <strong>Reporter ID:</strong> {report.account_id}
                  </div>
                  <hr style={{ margin: "8px 0", opacity: 0.2 }} />
                  <div>
                    <strong>Ad ID:</strong> {report.ad.ad_id}
                  </div>
                  <div>
                    <strong>Owner ID:</strong> {report.ad.user_id}
                  </div>
                  <div>
                    <strong>Level ID:</strong>{" "}
                    <a
                      href={`https://gdbrowser.com/${report.ad.level_id}`}
                      target="_blank"
                      rel="noopener noreferrer"
                      style={{
                        color: "#60a5fa",
                        textDecoration: "none",
                        cursor: "pointer",
                        transition: "color 0.2s ease",
                      }}
                      onMouseEnter={(e) =>
                        (e.currentTarget.style.color = "#3b82f6")
                      }
                      onMouseLeave={(e) =>
                        (e.currentTarget.style.color = "#60a5fa")
                      }
                    >
                      {report.ad.level_id}
                    </a>
                  </div>
                  <div>
                    <strong>Created:</strong>{" "}
                    {new Date(report.created_at).toLocaleString()}
                  </div>
                </div>
                <div className="pending-ad-actions">
                  <button
                    onClick={() => handleDeleteReportedAd(report.id)}
                    className="pending-ad-reject-button"
                    title="Delete the advertisement"
                  >
                    Delete Ad
                  </button>
                  <button
                    onClick={() => handleRejectReport(report.id)}
                    className="pending-ad-approve-button"
                    title="Reject the report (keep ad)"
                  >
                    Reject Report
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </>
    );
  }

  if (showingPending) {
    return (
      <>
        <div className="pending-ads-header">
          <h2>Pending Advertisements</h2>
        </div>
        <div className="pending-ads-actions">
          <button
            className="nine-slice-button small pending-ads-action-button"
            onClick={handlePendingAds}
          >
            Refresh
          </button>
          <button
            className="nine-slice-button small pending-ads-action-button"
            onClick={() => setShowingPending(false)}
          >
            Back
          </button>
        </div>
        {!pendingAds || pendingAds.length === 0 ? (
          <div className="pending-ads-empty">No pending advertisements</div>
        ) : (
          <div className="pending-ads-list">
            {pendingAds.map((ad) => (
              <div
                key={ad.ad_id}
                className="pending-ad-card"
                style={{
                  borderImage: `url(${square02}) 24 fill stretch`,
                }}
              >
                <div className="pending-ad-image-wrapper">
                  <a
                    href={ad.image_url}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="pending-ad-image-link"
                  >
                    <img
                      src={ad.image_url}
                      alt={`Ad ${ad.ad_id}`}
                      className="pending-ad-image"
                    />
                  </a>
                </div>
                <div className="pending-ad-info">
                  <div>
                    <strong>Ad ID:</strong> {ad.ad_id}
                  </div>
                  <div>
                    <strong>User ID:</strong> {ad.user_id}
                  </div>
                  <div>
                    <strong>Level ID:</strong>{" "}
                    <a
                      href={`https://gdbrowser.com/${ad.level_id}`}
                      target="_blank"
                      rel="noopener noreferrer"
                      style={{
                        color: "#60a5fa",
                        textDecoration: "none",
                        cursor: "pointer",
                        transition: "color 0.2s ease",
                      }}
                      onMouseEnter={(e) =>
                        (e.currentTarget.style.color = "#3b82f6")
                      }
                      onMouseLeave={(e) =>
                        (e.currentTarget.style.color = "#60a5fa")
                      }
                    >
                      {ad.level_id}
                    </a>
                  </div>
                  <div>
                    <strong>Type:</strong> {ad.type}
                  </div>
                  <div>
                    <strong>Created:</strong>{" "}
                    {new Date(ad.created_at).toLocaleString()}
                  </div>
                </div>
                <div className="pending-ad-actions">
                  <button
                    onClick={() => handleApproveAd(ad.ad_id)}
                    className="pending-ad-approve-button"
                  >
                    Approve
                  </button>
                  <button
                    onClick={() => handleRejectAd(ad.ad_id)}
                    className="pending-ad-reject-button"
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
      <h1 className="account-title">My Account</h1>
      <p className="account-subtitle">
        View and interact with your account information.
      </p>
      <p className="account-description">
        If you want to view your advertisements' stats in-game, copy the User ID
        into the settings prompted by the in-game popup.
      </p>
      {user && (
        <div className="account-user-info">
          <div className="account-user-details">
            <div>
              <strong>User ID: </strong> {user.id}{" "}
              <button
                onClick={handleCopyUserId}
                className="account-copy-button"
              >
                {copied ? (
                  <DoneIcon />
                ) : (
                  <ContentCopyIcon />
                )}
              </button>
            </div>
            <div>
              <strong>Account Created: </strong>{" "}
              {new Date(user.created_at).toLocaleString()}
            </div>
          </div>
        </div>
      )}

      <div className="account-actions">
        <button
          className="nine-slice-button small account-delete-button"
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

        {(user?.is_admin || user?.is_staff) && (
          <>
            <button
              className="nine-slice-button small"
              onClick={handlePendingAds}
            >
              Pending Ads ({pendingCount})
            </button>
            <button
              className="nine-slice-button small"
              onClick={handleReportedAds}
            >
              Reported Ads ({reportedCount})
            </button>
          </>
        )}

        {user?.is_admin && (
          <button
            className="nine-slice-button small account-admin-button"
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
