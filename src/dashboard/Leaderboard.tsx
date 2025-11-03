import "../page/Login.css";
import "./Dashboard.css";
import { useState, useEffect } from "react";

import AdminPanelSettingsIcon from "@mui/icons-material/AdminPanelSettings";
import BuildIcon from "@mui/icons-material/Build";
import VerifiedIcon from "@mui/icons-material/Verified";
import Avatar from "@mui/material/Avatar";

interface User {
  id: string;
  username: string;
  avatar_url: string;
  total_views: number;
  total_clicks: number;
  is_admin: boolean;
  is_staff: boolean;
  verified: boolean;
  banned: boolean;
  created_at: string;
  updated_at: string;
}

export default function Leaderboard() {
  const [leaderboardData, setLeaderboardData] = useState<User[]>([]);
  const [page, setPage] = useState(0);
  const [loading, setLoading] = useState(false);
  const [hasNext, setHasNext] = useState(true);
  const [sortBy, setSortBy] = useState<"views" | "clicks">("views");

  const MAX_USERS = 10;

  useEffect(() => {
    fetchLeaderboard(page);
  }, [page, sortBy]);

  const fetchLeaderboard = async (pageNum: number) => {
    setLoading(true);
    try {
      const endpoint = sortBy === "views" ? "views" : "clicks";
      const response = await fetch(
        `/ads/leaderboard/${endpoint}?page=${pageNum}&max=${MAX_USERS}`
      );
      if (response.ok) {
        const data = await response.json();
        setLeaderboardData(data || []);
        // Check if we got a full page of results to determine if there's a next page
        setHasNext((data && data.length === MAX_USERS) || false);
      } else {
        setLeaderboardData([]);
        setHasNext(false);
      }
    } catch (error) {
      console.error("Failed to fetch leaderboard:", error);
      setLeaderboardData([]);
      setHasNext(false);
    } finally {
      setLoading(false);
    }
  };

  const handleSort = (newSort: "views" | "clicks") => {
    setSortBy(newSort);
    setPage(0); // Reset to first page when changing sort
  };

  return (
    <div className="leaderboard-container">
      <h1 className="leaderboard-title">Leaderboard</h1>

      {/* Tabs */}
      <div className="leaderboard-tabs">
        <button
          onClick={() => handleSort("views")}
          className={`nine-slice-button small leaderboard-tab-button ${
            sortBy === "views" ? "active" : ""
          }`}
        >
          Total Views
        </button>
        <button
          onClick={() => handleSort("clicks")}
          className={`nine-slice-button small leaderboard-tab-button ${
            sortBy === "clicks" ? "active" : ""
          }`}
        >
          Total Clicks
        </button>
      </div>

      {loading ? (
        <p className="leaderboard-loading">Loading...</p>
      ) : leaderboardData.length === 0 ? (
        <p className="leaderboard-empty">No users found</p>
      ) : (
        <>
          <div className="rounded-lg shadow-lg">
            <table className="w-full border-collapse">
              <thead className="bg-gray-800 text-white">
                <tr>
                  <th className="px-4 py-3 text-left font-semibold">Rank</th>
                  <th className="px-4 py-3 text-left font-semibold">
                    Username
                  </th>
                  <th className="px-4 py-3 text-right font-semibold">
                    {sortBy === "views" ? "Total Views" : "Total Clicks"}
                  </th>
                </tr>
              </thead>
              <tbody>
                {leaderboardData.map((user, index) => (
                  <tr
                    key={user.id}
                    className="border-b border-gray-200 transition-colors"
                  >
                    <td className="px-4 py-3">
                      {page * MAX_USERS + index + 1}
                    </td>
                    <td className="px-4 py-3 flex items-center gap-2">
                      {user.avatar_url ? (
                        <Avatar alt={user.username} src={user.avatar_url} />
                      ) : (
                        <Avatar alt={user.username}>
                          {user.username.charAt(0).toUpperCase()}
                        </Avatar>
                      )}
                      {user.username}
                      {(user.is_admin && (
                        <AdminPanelSettingsIcon titleAccess="Administrator" />
                      )) ||
                        (user.is_staff && <BuildIcon titleAccess="Staff" />) ||
                        (user.verified && (
                          <VerifiedIcon titleAccess="Verified" />
                        ))}
                    </td>
                    <td className="px-4 py-3 text-right">
                      {sortBy === "views"
                        ? user.total_views
                        : user.total_clicks}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          <div className="leaderboard-pagination">
            <button
              onClick={() => setPage(Math.max(0, page - 1))}
              disabled={page === 0}
              className="nine-slice-button leaderboard-pagination-button"
            >
              Previous
            </button>
            <span>Page {page + 1}</span>
            <button
              onClick={() => setPage(page + 1)}
              disabled={!hasNext}
              className="nine-slice-button leaderboard-pagination-button"
            >
              Next
            </button>
          </div>
        </>
      )}
    </div>
  );
}
