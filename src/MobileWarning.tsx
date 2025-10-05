import { useEffect, useState } from "react";

export default function MobileWarning() {
  const [showMobileWarning, setShowMobileWarning] = useState(false);

  useEffect(() => {
    const check = () => {
      const isNarrow = window.matchMedia("(max-width: 900px)").matches;
      const isTouch = window.matchMedia("(pointer: coarse)").matches;
      if ((isNarrow || isTouch)) {
        setShowMobileWarning(true);
      } else {
        setShowMobileWarning(false);
      }
    };

    check();
    window.addEventListener("resize", check);
    return () => window.removeEventListener("resize", check);
  }, []);

  if (!showMobileWarning) return null;

  return (
    <div
      style={{
        position: "fixed",
        left: 0,
        top: 0,
        width: "100%",
        height: "100%",
        background: "rgba(0,0,0,1)",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        zIndex: 10000,
      }}
      aria-live="polite"
    >
      <div
        style={{
          background: "#0b0b0b",
          color: "#fff",
          padding: "1.25rem",
          borderRadius: 8,
          maxWidth: 640,
          width: "92%",
          boxShadow: "0 12px 30px rgba(0,0,0,0.6)",
        }}
      >
        <h1 style={{ margin: 0, marginBottom: "0.5rem" }}>Small screen detected</h1>
        <p style={{ margin: 0 }}>
          It looks like you're on a mobile device or your screen is too narrow to
          display the full site layout. For best results, view this site on a
          wider screen.
        </p>
        <div style={{ display: "flex", justifyContent: "flex-end", marginTop: "1rem" }}>
        </div>
      </div>
    </div>
  );
}
