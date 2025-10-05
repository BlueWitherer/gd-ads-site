import { useNavigate } from "react-router-dom";
import CreditsButton from "./Credits";
import "./App.css";
import "./Log.mjs";

export default function NotFound() {
  console.error("404 - Page Not Found");
  const navigate = useNavigate();

  const onBack = () => {
    console.info("Navigating back to home page");
    navigate("/");
  };

  return (
    <>
      <div id="background-scroll"></div>
      <div id="centered-container">
        <h1 className="text-6xl mt-8 mb-4">Oops!</h1>
        <h6 className="text-1xl mb-4">404 - Page Not Found</h6>

        <button className="nine-slice-button" onClick={onBack}>
          Go Back
        </button>
      </div>
      <CreditsButton />
    </>
  );
}
