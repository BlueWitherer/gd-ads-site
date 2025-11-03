import "../page/Login.css";
import "./Dashboard.css";
import square02 from "../assets/square02.png";
import { useEffect, useState } from "react";

type Ad = {
  id: number;
  type: string;
  level_id: string;
  image: string;
  expiration: number;
  pending?: boolean;
};

function getDaysRemaining(expirationTimestamp: number): {
  days: number;
  color: string;
} {
  const now = Date.now();
  const expirationMs = expirationTimestamp * 1000; // Convert seconds to milliseconds
  const diffMs = expirationMs - now;
  const days = Math.ceil(diffMs / (1000 * 60 * 60 * 24));

  let color = "#e74c3c"; // Red (1 day or less)
  if (days >= 5) {
    color = "#27ae60"; // Green (7-5 days)
  } else if (days >= 2) {
    color = "#f39c12"; // Orange (4-2 days)
  }

  return { days: Math.max(0, days), color };
}

function Manage() {
  const [adverts, setAdverts] = useState<Ad[] | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function load() {
      try {
        const res = await fetch("/ads/get", { credentials: "include" });
        if (!res.ok) {
          setError(`Failed to fetch ads: ${res.status}`);
          return;
        }

        const data = await res.json();
        // Expecting array of { id, type, level_id, image, expiration }
        setAdverts(
          data.map((a: any) => ({
            id: a.ad_id,
            type: a.type,
            level_id: a.level_id,
            image: a.image_url,
            expiration: a.expiry,
            pending: a.pending,
          }))
        );
      } catch (err: any) {
        setError(err.message || String(err));
      }
    }

    load();
  }, []);

  return (
    <>
      <h1 className="manage-title">Manage Advertisements</h1>
      <p className="manage-subtitle">
        Manage and configure your active advertisements.
      </p>
      <p className="manage-description">
        You can manually delete your advertisement or wait until the expiration
        date if you want to make a new one.
      </p>

      {error && <div className="manage-error">{error}</div>}

      <div className="manage-ads-list">
        {adverts === null ? (
          <div>Loading advertisements...</div>
        ) : adverts.length === 0 ? (
          <div>No advertisements found.</div>
        ) : (
          adverts.map((advert) => (
            <div
              key={advert.id}
              className="ad-card manage-ad-card"
              style={{
                borderImage: `url(${square02}) 24 fill stretch`,
              }}
            >
              <div className="manage-ad-content">
                <a
                  href={advert.image}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="manage-ad-image-link"
                >
                  <img
                    src={advert.image}
                    alt="Advertisement"
                    className="ad-card-image manage-ad-image"
                  />
                </a>
                <div className="ad-card-info manage-ad-info">
                  <div>
                    <strong>Ad ID:</strong> {advert.id}
                  </div>
                  <div>
                    <strong>Type:</strong> {advert.type}
                  </div>
                  <div>
                    <strong>Level ID:</strong> {advert.level_id}
                  </div>
                  <div>
                    <strong>Expiration:</strong>{" "}
                    {(() => {
                      const { days, color } = getDaysRemaining(
                        advert.expiration
                      );
                      return (
                        <span style={{ color, fontWeight: "bold" }}>
                          {days} day{days !== 1 ? "s" : ""}
                        </span>
                      );
                    })()}
                  </div>
                </div>
              </div>
              <div className="ad-card-badge manage-ad-badge-wrapper">
                {advert.pending ? (
                  <div className="manage-ad-pending-badge">PENDING</div>
                ) : (
                  <div className="manage-ad-approved-badge">APPROVED</div>
                )}
              </div>
              <button
                className="ad-card-delete manage-ad-delete-button"
                onClick={() => {
                  if (
                    confirm(
                      "Are you sure you want to delete this advertisement? This action cannot be undone."
                    )
                  ) {
                    fetch(`/ads/delete?id=${advert.id}`, {
                      method: "DELETE",
                      credentials: "include",
                    }).then(() => {
                      adverts.splice(adverts.indexOf(advert), 1);
                      setAdverts([...adverts]);
                    });
                  }
                }}
              >
                Delete
              </button>
            </div>
          ))
        )}
      </div>
    </>
  );
}

export default Manage;
