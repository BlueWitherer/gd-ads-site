import '../App.css'

function Statistics() {
  // replace with real data later via backend
  const totalViews = 420;
  const totalClicks = 69;

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
