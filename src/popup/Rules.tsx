import { useState } from "react";
import { createPortal } from "react-dom";
import "../page/Login.css";
import "./Rules.css";

import WarningIcon from "@mui/icons-material/WarningOutlined";
import InfoIcon from "@mui/icons-material/InfoOutlined";
import { FaDiscord } from "react-icons/fa";

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
                      Do not upload generally inappropriate or controversial
                      ads.
                    </p>
                  </div>
                  <div className="rule-item">
                    <WarningIcon className="rule-icon" />
                    <p>
                      Do not use ads to promote anything non-Geometry Dash
                      related. Memes or well-known creators are allowed.
                    </p>
                  </div>
                  <div className="rule-item">
                    <WarningIcon className="rule-icon" />
                    <p>No profanity or offensive text in the ad.</p>
                  </div>
                  <div className="rule-item">
                    <WarningIcon className="rule-icon" />
                    <p>
                      Do not promote any harmful, illegal, or offensive
                      material, including both within your Geometry Dash level
                      and ad!
                    </p>
                  </div>
                  <div className="rule-item">
                    <WarningIcon className="rule-icon" />
                    <p>AI-generated ads are not allowed.</p>
                  </div>
                  <div className="rule-item">
                    <WarningIcon className="rule-icon" />
                    <p>
                      Do not post the same ad image multiple times. Try add variety
                      on each advertisement. You can make different versions of
                      your ad in the same level.
                    </p>
                  </div>
                  <div className="rule-item warning">
                    <WarningIcon className="rule-icon" />
                    <p>
                      Do not attempt to inflate your views/clicks using bots or
                      other methods. Getting caught doing this will almost
                      immediately result in a ban.
                    </p>
                  </div>
                  <div className="rule-item warning">
                    <WarningIcon className="rule-icon" />
                    <p>
                      Violating the rules may result in a ban or removal of your
                      ad/s without warning. No appeals!
                    </p>
                  </div>
                  <div className="rule-item info">
                    <InfoIcon className="rule-icon" />
                    <p>
                      Enforcement of these rules is at the discretion of staff.
                      If you are unsure about something, feel free to ask in the{" "}
                      <a
                        href="https://discord.gg/gXcppxTNxC"
                        target="_blank"
                        rel="noreferrer"
                        style={{
                          display: "inline-flex",
                          alignItems: "center",
                          gap: "0.25rem",
                        }}
                      >
                        <FaDiscord size={15} aria-hidden="true" /> Discord
                        server
                      </a>
                      .
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
