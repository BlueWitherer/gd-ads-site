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
        className="nine-slice-button small rules-button"
        onClick={() => setOpen(true)}
      >
        Submission Rules
      </button>
      {open &&
        createPortal(
          <>
            <div
              className="rules-popup-overlay"
              onClick={() => setOpen(false)}
            />
            <div className="rules-popup-bg rules-popup-elastic">
              <h2 className="rules-title">Advertisement Submission Rules</h2>
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
                    <p>AI generated advertisements are not allowed.</p>
                  </div>
                  <div className="rule-item">
                    <WarningIcon className="rule-icon" />
                    <p>
                      Do not post the same advertisement multiple times. Most
                      likely your duplicated advertisements will be rejected or
                      deleted.
                    </p>
                  </div>
                  <div className="rule-item warning">
                    <WarningIcon className="rule-icon" />
                    <p>
                      Do not attempt to inflate your views/clicks using bots or
                      other methods. Not only this negatively impacts the server
                      performance, but ruins the experiences for other users.
                      Caught doing this will result in an immediate ban.
                    </p>
                  </div>
                  <div className="rule-item warning">
                    <WarningIcon className="rule-icon" />
                    <p>
                      Violating the rules may result in a ban or your
                      advertisement/s removal without warning. No appeals!
                    </p>
                  </div>
                </div>
              </div>
              <button
                className="nine-slice-button small rules-close-button"
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
