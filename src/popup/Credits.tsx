import { useState } from "react";
import "../page/Login.css";
import "./Credits.css";
import newsIcon from "../assets/newsIcon.png";

export default function CreditsButton() {
  const [open, setOpen] = useState(false);
  return (
    <>
      <button
        className="sprite-button credits-button"
        onClick={() => setOpen(true)}
      >
        <img src={newsIcon} alt="Credits" />
      </button>
      {open && (
        <>
          <div
            className="credit-popup-overlay"
            onClick={() => setOpen(false)}
          />
          <div className="credit-popup-bg credit-popup-elastic">
            <h2 className="credits-title">Credits</h2>
            <div className="credits-content">
              <div className="credits-grid">
                <div className="credits-person">
                  <a
                    href="https://github.com/DumbCaveSpider"
                    target="_blank"
                    rel="noopener noreferrer"
                  >
                    <img
                      src="https://avatars.githubusercontent.com/u/56347227"
                      alt="ArcticWoof avatar"
                      className="credits-avatar"
                    />
                  </a>
                  <span className="credits-name">ArcticWoof</span>
                  <span className="credits-role">Frontend/UI/UX</span>
                  <span className="credits-role">Geode Mod</span>
                </div>
                <div className="credits-person">
                  <a
                    href="https://github.com/BlueWitherer"
                    target="_blank"
                    rel="noopener noreferrer"
                  >
                    <img
                      src="https://avatars.githubusercontent.com/u/47698640?v=4"
                      alt="Cheeseworks avatar"
                      className="credits-avatar"
                    />
                  </a>
                  <span className="credits-name">Cheeseworks</span>
                  <span className="credits-role">Backend/API</span>
                </div>
              </div>
              <p className="credits-footer-text">
                Assets by Geode Team and RobTop Games.
              </p>
            </div>
            <a
              href="https://arcticwoof.com.au/privacy/"
              target="_blank"
              rel="noopener noreferrer"
              className="credits-privacy-link"
            >
              Privacy Policy
            </a>
            <button
              className="nine-slice-button small credits-close-button"
              onClick={() => setOpen(false)}
            >
              Close
            </button>
          </div>
        </>
      )}
    </>
  );
}
