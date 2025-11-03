import { useState } from "react";
import { createPortal } from "react-dom";
import "../page/Login.css";
import "./Rules.css";
import WarningIcon from "@mui/icons-material/WarningOutlined";

export default function RulesButton() {
  const [open, setOpen] = useState(false);
  return (
    <>
      <button
        className="nine-slice-button rules-button"
        onClick={() => setOpen(true)}
      >
        Rules
      </button>
      {open &&
        createPortal(
          <>
            <div
              className="rules-popup-overlay"
              onClick={() => setOpen(false)}
            />
            <div className="rules-popup-bg rules-popup-elastic">
              <h2 className="rules-title">Advertisement Rules</h2>
              <div className="rules-content">
                <div className="rules-list">
                  <div className="rule-item">
                    <WarningIcon className="rule-icon" />
                    <p>
                      Do not upload inappropriate or controversial
                      advertisements.
                    </p>
                  </div>
                  <div className="rule-item">
                    <WarningIcon className="rule-icon" />
                    <p>
                      Do not self-promote anything non-Geometry Dash related.
                      Memes or well-known creators are allowed.
                    </p>
                  </div>
                  <div className="rule-item">
                    <WarningIcon className="rule-icon" />
                    <p>No profanity or offensive text in the advertisement.</p>
                  </div>
                  <div className="rule-item">
                    <WarningIcon className="rule-icon" />
                    <p>
                      Do not promote any harmful, illegal, or offensive material
                      including both your Geometry Dash level and your
                      advertisement!
                    </p>
                  </div>
                  <div className="rule-item">
                    <WarningIcon className="rule-icon" />
                    <p>AI Generated advertisements are hard rejection.</p>
                  </div>
                  <div className="rule-item warning">
                    <WarningIcon className="rule-icon" />
                    <p>
                      Violating this may result in a ban. <b>No appeals!</b>
                    </p>
                  </div>
                </div>
              </div>
              <button
                className="nine-slice-button large rules-close-button"
                onClick={() => setOpen(false)}
              >
                Close
              </button>
            </div>
          </>,
          document.body
        )}
    </>
  );
}
