import { useEffect, useState } from "react";
import "../App.css";
import "../Log.mjs";
import { PieChart } from "@mui/x-charts/PieChart";

export default function Statistics() {
  console.debug("Rendering Statistics component...");

  const [views, setViews] = useState<number | null>(null);
  const [clicks, setClicks] = useState<number | null>(null);
  const [globalViews, setGlobalViews] = useState<number | null>(null);
  const [globalClicks, setGlobalClicks] = useState<number | null>(null);
  const [adCount, setAdCount] = useState<number | null>(null);

  useEffect(() => {
    async function fetchStats(): Promise<number> {
      try {
        const res = await fetch(`/stats/get`);
        if (res.ok) return await res.json().then((data) => {
          setViews(data.views);
          setClicks(data.clicks);

          return 1;
        });

        return 1;
      } catch (err) {
        console.error("Error fetching ad stats:", err);
        return -1;
      };
    };

    fetchStats();
  }, []);

  useEffect(() => {
    async function fetchGlobalStats(): Promise<void> {
      try {
        const res = await fetch(`/stats/global`);
        if (res.ok) {
          const data = await res.json();
          setGlobalViews(data.total_views);
          setGlobalClicks(data.total_clicks);
          setAdCount(data.ad_count);
        }
      } catch (err) {
        console.error("Error fetching global stats:", err);
      }
    }

    fetchGlobalStats();
  }, []);

  const settings = {
    width: 200,
    height: 200,
    hideLegend: true,
  };

  return (
    <>
      <h1 className="text-2xl font-bold mb-6">Statistics</h1>
        <div className="text-sm mb-6 p-3 rounded bg-yellow-900/30 border border-yellow-500">
          <p>Please update the Geode Mod to v1.0.6 or higher for clicks and views to be tracked properly.</p>
        </div>
      <div style={{ 
        display: "flex", 
        gap: "24px", 
        alignItems: "flex-start",
        flexWrap: "wrap",
        justifyContent: "center",
        marginBottom: "48px"
      }}>
        <div className="stat-box" style={{ flex: "0 0 auto", minWidth: 0 }}>
          {views !== null && clicks !== null && globalViews !== null && globalClicks !== null ? (
            <PieChart
              series={[
                {
                  data: [
                    { id: 0, value: views, label: "My Views", color: "#2196f3" },
                    { id: 1, value: clicks, label: "My Clicks", color: "#4caf50" },
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
              height={250}
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
            flex: "1 1 auto",
            minWidth: "200px",
          }}
        >
          <div className="stat-box">
            <h2 className="text-xl font-bold mb-2">Your Total Views</h2>
            <p className="text-4xl font-bold">
              {views !== null ? views.toLocaleString() : "Loading..."}
            </p>
          </div>
          <div className="stat-box">
            <h2 className="text-xl font-bold mb-2">Your Total Clicks</h2>
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

      <h1 className="text-2xl font-bold mb-6">Global Statistics</h1>
      <div style={{ 
        display: "flex", 
        gap: "24px", 
        flexWrap: "wrap",
        justifyContent: "center"
      }}>
        <div className="stat-box">
          <h3 className="text-xl font-bold mb-2">Total Views</h3>
          <p className="text-4xl font-bold">
            {globalViews !== null ? globalViews.toLocaleString() : "Loading..."}
          </p>
        </div>
        <div className="stat-box">
          <h3 className="text-xl font-bold mb-2">Total Clicks</h3>
          <p className="text-4xl font-bold">
            {globalClicks !== null ? globalClicks.toLocaleString() : "Loading..."}
          </p>
        </div>
        <div className="stat-box">
          <h3 className="text-xl font-bold mb-2">Active Advertisements</h3>
          <p className="text-4xl font-bold">
            {adCount !== null ? adCount.toLocaleString() : "Loading..."}
          </p>
        </div>
      </div>
    </>
  );
}
