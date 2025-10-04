import { useEffect, useState } from "react";
import "../App.css";
import "../Log.mjs";
import { PieChart } from "@mui/x-charts/PieChart";

export default function Statistics() {
  console.debug("Rendering Statistics component...");

  const [views, setViews] = useState<number | null>(null);
  const [clicks, setClicks] = useState<number | null>(null);

  useEffect(() => {
    async function fetchStats(endpoint: string): Promise<number> {
      try {
        const res = await fetch(`/api/stats/${endpoint}`);
        if (res.ok) return parseInt(await res.text());
        return 1;
      } catch (err) {
        console.error("Error fetching stats:", err);
        return -1;
      }
    }

    fetchStats("views").then(setViews);
    fetchStats("clicks").then(setClicks);
  }, []);

  const settings = {
    width: 200,
    height: 200,
    hideLegend: true,
  };
  return (
    <>
      <h1 className="text-2xl font-bold mb-6">Statistics</h1>

      <div style={{ display: "flex", gap: "24px", alignItems: "flex-start" }}>
        <div className="stat-box" style={{ flex: "0 0 auto" }}>
          {views !== null && clicks !== null ? (
            <PieChart
              series={[
                {
                  data: [
                    { id: 0, value: views, label: "Views", color: "#2196f3" },
                    { id: 1, value: clicks, label: "Clicks", color: "#4caf50" },
                  ],
                  highlightScope: { fade: "global", highlight: "item" },
                  faded: {
                    innerRadius: 30,
                    outerRadius: 50,
                    color: "gray",
                  },
                },
              ]}
              {...settings}
              width={300}
              height={200}
            />
          ) : (
            <p>Loading chart...</p>
          )}
        </div>
        <div
          style={{
            display: "flex",
            flexDirection: "column",
            gap: "16px",
            flex: "1",
          }}
        >
          <div className="stat-box">
            <h2 className="text-xl font-bold mb-2">Total Views</h2>
            <p className="text-4xl font-bold">
              {views !== null ? views.toLocaleString() : "Loading..."}
            </p>
          </div>
          <div className="stat-box">
            <h2 className="text-xl font-bold mb-2">Total Clicks</h2>
            <p className="text-4xl font-bold">
              {clicks !== null ? clicks.toLocaleString() : "Loading..."}
            </p>
            <div
              style={{
                fontSize: "1.2rem",
                color: "#4caf50",
                marginTop: "12px",
              }}
            >
              {views !== null && clicks !== null && views > 0 ? (
                <span>
                  Click/View Ratio: {((clicks / views) * 100).toFixed(2)}%
                </span>
              ) : (
                <span>Click/View Ratio: N/A</span>
              )}
            </div>
          </div>
        </div>
      </div>
    </>
  );
}
