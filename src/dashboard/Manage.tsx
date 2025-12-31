import "../page/Login.css";
import "./Dashboard.css";
import square02 from "../assets/square02.png";
import { useEffect, useState } from "react";
import { SiKofi } from "react-icons/si";

type Ad = {
    id: number;
    type: string;
    level_id: string;
    image: string;
    expiration: number;
    pending?: boolean;
    views?: number;
    clicks?: number;
    boost_count?: number;
};



function getDaysRemaining(expirationTimestamp: number): {
    days: number;
    color: string;
} {
    const now = Date.now();
    const expirationMs = expirationTimestamp * 1000; // Convert seconds to milliseconds
    const diffMs = expirationMs - now;
    const days = Math.ceil(diffMs / (1000 * 60 * 60 * 24));

    let color = "#e74c3c"; // Red (2 days or less)
    if (days >= 7) {
        color = "#27ae60"; // Green (7-14 days)
    } else if (days >= 3) {
        color = "#f39c12"; // Orange (3-6 days)
    };

    return { days: Math.max(0, days), color };
};

function Manage() {
    const [adverts, setAdverts] = useState<Ad[] | null>(null);
    const [boostAmounts, setBoostAmounts] = useState<{ [key: number]: string }>({});
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        async function load() {
            try {
                const res = await fetch("/ads/get", { credentials: "include" });
                if (!res.ok) {
                    setError(`Failed to fetch ads: ${res.status}`);
                    return;
                };

                const data = await res.json();

                if (!data) {
                    setAdverts([]);
                    return;
                }

                // Expecting array of { id, type, level_id, image, expiration }
                setAdverts(
                    data.map((a: any) => ({
                        id: a.ad_id,
                        type: a.type,
                        level_id: a.level_id,
                        image: a.image_url,
                        expiration: a.expiry,
                        pending: a.pending,
                        views: a.views || 0,
                        clicks: a.clicks || 0,
                        boost_count: a.boost_count || 0,
                    }))
                );
            } catch (err: any) {
                setError(err.message || String(err));
            };
        };

        load();
    }, []);

    const handleBoost = async (id: number) => {
        const amount = parseInt(boostAmounts[id] || "0", 10);
        if (isNaN(amount) || amount < 1) {
            alert("Please enter a valid boost amount (minimum 1).");
            return;
        }

        try {
            const res = await fetch(`/ads/boost?id=${id}&boosts=${amount}`, {
                method: "POST",
                credentials: "include",
            });

            if (!res.ok) {
                const text = await res.text();
                alert(`Failed to boost ad: ${text}`);
                return;
            }

            alert("Ad boosted successfully!");
            setBoostAmounts((prev) => ({ ...prev, [id]: "" }));
            const refreshRes = await fetch("/ads/get", { credentials: "include" });
            if (refreshRes.ok) {
                const data = await refreshRes.json();
                if (data) {
                    setAdverts(
                        data.map((a: any) => ({
                            id: a.ad_id,
                            type: a.type,
                            level_id: a.level_id,
                            image: a.image_url,
                            expiration: a.expiry,
                            pending: a.pending,
                            views: a.views || 0,
                            clicks: a.clicks || 0,
                            boost_count: a.boost_count || 0,
                        }))
                    );
                }
            }
        } catch (err: any) {
            alert(`Error: ${err.message || String(err)}`);
        }
    };

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

            <a
                href="https://ko-fi.com/playerads"
                target="_blank"
                rel="noopener noreferrer"
                className="kofi-button"
            >
                <SiKofi size={24} />
                <span>Get More Reach!</span>
            </a>

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
                                <div className="ad-content-wrapper">
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
                                    <div className="ad-card-stats manage-ad-stats">
                                        <div className="stat-item">
                                            <strong>Views: {advert.views || 0}</strong>
                                        </div>
                                        <div className="stat-item">
                                            <strong>Clicks: {advert.clicks || 0}</strong>
                                        </div>
                                        <div className="stat-item">
                                            <strong>Boosts: {advert.boost_count || 0}</strong>
                                        </div>
                                        <div className="stat-item">
                                            <strong>
                                                Click/View Ratio: {advert.views && advert.views > 0
                                                    ? (
                                                        ((advert.clicks || 0) / advert.views) *
                                                        100
                                                    ).toFixed(2)
                                                    : "0.00"}
                                                %
                                            </strong>
                                        </div>
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
                            <div className="manage-ad-boost-wrapper">
                                <input
                                    type="number"
                                    min="1"
                                    placeholder="Amount"
                                    className="manage-ad-boost-input"
                                    value={boostAmounts[advert.id] || ""}
                                    onChange={(e) =>
                                        setBoostAmounts((prev) => ({
                                            ...prev,
                                            [advert.id]: e.target.value,
                                        }))
                                    }
                                />
                                <button
                                    className="manage-ad-boost-button"
                                    onClick={() => handleBoost(advert.id)}
                                >
                                    Boost Ad
                                </button>
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
};

export default Manage;
