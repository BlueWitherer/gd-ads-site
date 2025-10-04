import "./App.css";
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import Statistics from "./views/Statistics";
import Create from "./views/Create";
import Leaderboard from "./views/Leaderboard";
import Manage from "./views/Manage";
import Account from "./views/Account";
import CreditsButton from "./Credits";

export default function Dashboard() {
  const navigate = useNavigate();
  const [selectedView, setSelectedView] = useState<
    "statistics" | "create" | "leaderboard" | "manage" | "account"
  >("statistics");

  const logout = () => {
    // handle logout logic here
    navigate("/"); // Redirect to login page
  };

  const renderContent = () => {
    switch (selectedView) {
      case "statistics":
        return <Statistics />;
      case "create":
        return <Create />;
      case "leaderboard":
        return <Leaderboard />;
      case "manage":
        return <Manage />;
      case "account":
        return <Account />;
      default:
        return null;
    }
  };

  return (
    <>
      <div id="background-scroll"></div>
      <div id="centered-container">
        <div className="sidebar-container">
          <button
            className="nine-slice-button padding-4 mt-4 mb-4"
            onClick={() => setSelectedView("create")}
          >
            Create
          </button>
          <button
            className="nine-slice-button padding-4 mt-4 mb-4"
            onClick={() => setSelectedView("statistics")}
          >
            Statistics
          </button>
          <button
            className="nine-slice-button padding-4 mt-4 mb-4"
            onClick={() => setSelectedView("leaderboard")}
          >
            Leaderboard
          </button>
          <button
            className="nine-slice-button padding-4 mt-4 mb-4"
            onClick={() => setSelectedView("manage")}
          >
            Manage
          </button>
          <button
            className="nine-slice-button padding-4 mt-4 mb-4"
            onClick={() => setSelectedView("account")}
          >
            Account
          </button>
        </div>
        <div className="user-container">
          <p className="text-lg">User Info</p>
          <button
            title="Log Out"
            className="nine-slice-button"
            onClick={logout}
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="32"
              height="32"
              fill="currentColor"
              className="bi bi-box-arrow-left"
              viewBox="0 0 16 16"
            >
              <path
                fill-rule="evenodd"
                d="M6 12.5a.5.5 0 0 0 .5.5h8a.5.5 0 0 0 .5-.5v-9a.5.5 0 0 0-.5-.5h-8a.5.5 0 0 0-.5.5v2a.5.5 0 0 1-1 0v-2A1.5 1.5 0 0 1 6.5 2h8A1.5 1.5 0 0 1 16 3.5v9a1.5 1.5 0 0 1-1.5 1.5h-8A1.5 1.5 0 0 1 5 12.5v-2a.5.5 0 0 1 1 0z"
              />
              <path
                fill-rule="evenodd"
                d="M.146 8.354a.5.5 0 0 1 0-.708l3-3a.5.5 0 1 1 .708.708L1.707 7.5H10.5a.5.5 0 0 1 0 1H1.707l2.147 2.146a.5.5 0 0 1-.708.708z"
              />
            </svg>
          </button>
        </div>
        <div className="dashboard-container">{renderContent()}</div>
      </div>
      <CreditsButton />
    </>
  );
}
