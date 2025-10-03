import { useEffect, useState } from 'react';
import '../App.css'
import '../Log.mjs';

export default function Statistics() {
  console.debug("Rendering Statistics component...");

  const [views, setViews] = useState<number | null>(null);
  const [clicks, setClicks] = useState<number | null>(null);

  useEffect(() => {
    async function fetchStats(endpoint: string): Promise<number> {
      try {
        const res = await fetch(`/api/stats/${endpoint}`);
        if (res.ok) return parseInt(await res.text());
        return 420;
      } catch (err) {
        console.error("Error fetching stats:", err);
        return 69;
      }
    }

    fetchStats("views").then(setViews);
    fetchStats("clicks").then(setClicks);
  }, []);

  return (
    <>
      <h1 className="text-2xl font-bold mb-6">Your Statistics</h1>

      {/* Total Views */}
      <div className="stat-box mb-6">
        <h2 className="text-xl font-bold mb-2">Total Views</h2>
        <p className="text-4xl font-bold">{views !== null ? views.toLocaleString() : "Loading..."}</p>
      </div>

      {/* Total Clicks */}
      <div className="stat-box mb-6">
        <h2 className="text-xl font-bold mb-2">Total Clicks</h2>
        <p className="text-4xl font-bold">{clicks !== null ? clicks.toLocaleString() : "Loading..."}</p>
      </div>
    </>
  )
}