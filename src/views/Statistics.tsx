import '../App.css'

async function Statistics() {
  // replace with real data later via backend
  let totalViews = 420;
  let totalClicks = 69;

  console.debug("Rendering Statistics component...");

  // get the views count
  const resViews = await fetch('/api/stats/views');
  totalViews = parseInt(await resViews.text());

  // get the clicks count
  const resClicks = await fetch('/api/stats/clicks');
  totalClicks = parseInt(await resClicks.text());

  console.info("Returning Statistics component");

  return (
    <>
      <h1 className="text-2xl font-bold mb-6">Your Statistics</h1>
      
      {/* Total Views */}
      <div className="stat-box mb-6">
        <h2 className="text-xl font-bold mb-2">Total Views</h2>
        <p className="text-4xl font-bold">{totalViews.toLocaleString()}</p>
      </div>

      {/* Total Clicks */}
      <div className="stat-box mb-6">
        <h2 className="text-xl font-bold mb-2">Total Clicks</h2>
        <p className="text-4xl font-bold">{totalClicks.toLocaleString()}</p>
      </div>
    </>
  )
}

export default Statistics
