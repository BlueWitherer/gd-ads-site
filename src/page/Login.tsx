import "./Login.css";
import CreditsButton from "../popup/Credits";
import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { FaDiscord } from "react-icons/fa";
import "../misc/Log.mjs";

export async function copyText(
    text: string | undefined,
    setState: React.Dispatch<React.SetStateAction<any | null>>
) {
    try {
        if (text) {
            await navigator.clipboard.writeText(text);
            setState(true);
            setTimeout(() => setState(false), 2500);
        } else {
            console.error("No text provided to copy");
        }
    } catch (err) {
        console.error("Copy failed:", err);
    }
}

export default function App() {
    const navigate = useNavigate();
    const [randomAds, setRandomAds] = useState<
        Array<{
            url: string;
            id: number;
            top: number;
            scale: number;
            delay: number;
            speed: number;
            fadeIn: boolean;
        }>
    >([]);
    const [allImagesRef, setAllImagesRef] = useState<string[]>([]);
    const [nextIdRef, setNextIdRef] = useState(0);

    useEffect(() => {
        fetch("/session", { credentials: "include" })
            .then((res) => (res.ok ? res.json() : null))
            .then((data) => {
                if (data?.username && data?.id) navigate("/dashboard");
            })
            .catch(() => {
                console.error("User unauthorized");
            });
    }, [navigate]);

    useEffect(() => {
        async function fetchAndInitializeAds() {
            try {
                const adTypes = ["banner", "square", "skyscraper"];
                const allImages: string[] = [];
                for (const adType of adTypes) {
                    try {
                        const res = await fetch(`/cdn/${adType}/`);
                        if (res.ok) {
                            const html = await res.text();
                            const imageRegex = /href="([^"]+\.webp)"/g;
                            const matches = html.matchAll(imageRegex);
                            const images = Array.from(matches)
                                .map((m) => m[1])
                                .filter((img) => img !== "../");

                            images.forEach((img) => {
                                allImages.push(`/cdn/${adType}/${img}`);
                            });
                        }
                    } catch (err) {
                        console.error(`Failed to fetch ${adType} ads:`, err);
                    }
                }
                if (allImages.length > 0) {
                    setAllImagesRef(allImages);
                    const initialAds = Array.from({ length: 15 }, (_, i) => {
                        const randomImage =
                            allImages[Math.floor(Math.random() * allImages.length)];
                        return {
                            url: randomImage,
                            id: i,
                            top: Math.random() * 80,
                            scale: 1,
                            delay: -(15 - i * 1),
                            speed: 10 + Math.random() * 15,
                            fadeIn: false,
                        };
                    });
                    setRandomAds(initialAds);
                    setNextIdRef(15);
                }
            } catch (err) {
                console.error("Failed to fetch random ads:", err);
            }
        }

        fetchAndInitializeAds();
    }, []);

    useEffect(() => {
        if (allImagesRef.length === 0) return;

        const MAX_ADS = 15;
        const interval = setInterval(() => {
            setRandomAds((prevAds) => {
                if (prevAds.length >= MAX_ADS) {
                    return prevAds;
                }

                const newAd = {
                    url: allImagesRef[Math.floor(Math.random() * allImagesRef.length)],
                    id: nextIdRef,
                    top: Math.random() * 80,
                    scale: 1,
                    delay: 0,
                    speed: 10 + Math.random() * 15,
                    fadeIn: true,
                };
                setNextIdRef((prev) => prev + 1);
                return [...prevAds, newAd];
            });
        }, 2000);

        return () => clearInterval(interval);
    }, [allImagesRef]);

    useEffect(() => {
        const timer = setTimeout(() => {
            setRandomAds((prevAds) =>
                prevAds.map((ad) => ({ ...ad, fadeIn: false }))
            );
        }, 800);

        return () => clearTimeout(timer);
    }, [randomAds.length]);

    useEffect(() => {
        const timers: ReturnType<typeof setTimeout>[] = [];

        randomAds.forEach((ad) => {
            const totalTime = ad.speed + 0.8 + ad.speed * 0.05;
            const removeTimer = setTimeout(() => {
                setRandomAds((prevAds) => prevAds.filter((a) => a.id !== ad.id));
            }, totalTime * 1000);

            timers.push(removeTimer);
        });

        return () => {
            timers.forEach((timer) => clearTimeout(timer));
        };
    }, [randomAds]);

    const handleLogin = () => {
        window.location.href = "/login";
    };

    return (
        <>
            <div id="background-scroll"></div>
            <div id="centered-container">
                <div className="ads-container">
                    {randomAds.map((ad) => (
                        <img
                            key={ad.id}
                            src={ad.url}
                            alt={`Advertisement ${ad.id}`}
                            className={`ad-slide ${ad.fadeIn ? "fade-in" : ""}`}
                            style={{
                                top: `${ad.top}%`,
                                transform: `scale(${ad.scale})`,
                                animationDelay: `${ad.delay}s`,
                                animationDuration: `${ad.speed}s`,
                            }}
                        />
                    ))}
                </div>
                {/* Login Section */}
                <div id="login-section" className="login-section-wrapper">
                    {/* <div className="maintenance-banner">
            <p>
              Views and Clicks does not currently work on v1.0.6 or older.
              Please update your Geode Mod to v1.0.7 or newer.
            </p>
          </div> */}
                    <h1 className="login-title">GD Ads Manager</h1>
                    <h2>
                        Welcome to the GD Ads Manager! Manage all your Geometry
                        Dash Advertisements here!
                    </h2>
                    <h2 className="login-subtitle">
                        Login using your Discord Account to get started!
                    </h2>
                    <button
                        className="nine-slice-button login-button"
                        onClick={handleLogin}
                        aria-label="Login with Discord"
                    >
                        <FaDiscord size={25} aria-hidden="true" />
                        <span>Login with Discord</span>
                    </button>
                    <button
                        className="nine-slice-button login-button install-mod-button"
                        onClick={() =>
                            window.open(
                                "https://geode-sdk.org/mods/arcticwoof.player_advertisements",
                                "_blank"
                            )
                        }
                        aria-label="Install Geode Mod"
                    >
                        <span>Install Geode Mod</span>
                    </button>
                    <div className="footer-text">
                        Made with 💝 by ArcticWoof & Cheeseworks
                    </div>
                </div>
            </div>
            <CreditsButton />
        </>
    );
}
