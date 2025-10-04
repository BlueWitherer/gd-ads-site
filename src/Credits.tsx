import { useState } from 'react';
import './App.css';
import newsIcon from './assets/newsIcon.png';

export default function CreditsButton() {
  const [open, setOpen] = useState(false);
  return (
    <>
      {/* Sprite Button */}
      <button
        className="sprite-button"
        style={{
          position: 'fixed',
          right: '32px',
          bottom: '32px',
          zIndex: 1000,
          width: '64px',
          height: '64px',
          background: 'transparent',
          border: 'none',
          padding: 0,
          cursor: 'pointer',
          transition: 'transform 0.1s ease',
        }}
        onClick={() => setOpen(true)}
        onMouseDown={e => e.currentTarget.style.transform = 'scale(0.95)'}
        onMouseUp={e => e.currentTarget.style.transform = 'scale(1)'}
        onMouseLeave={e => e.currentTarget.style.transform = 'scale(1)'}
      >
        <img
          src={newsIcon}
          alt="Credits"
          style={{ width: '100%', height: '100%', objectFit: 'contain' }}
        />
      </button>

      {/* Overlay and Popup */}
      {open && (
        <>
          {/* Dimmed overlay */}
          <div
            className="credit-popup-overlay"
            style={{
              position: 'fixed',
              top: 0,
              left: 0,
              width: '100vw',
              height: '100vh',
              background: 'rgba(0,0,0,0.7)',
              zIndex: 1999,
              animation: 'fadeIn 0.3s',
            }}
            onClick={() => setOpen(false)}
          />
          {/* Popup with elastic animation */}
          <div
            className="credit-popup-bg credit-popup-elastic"
            style={{
              position: 'fixed',
              left: '50%',
              top: '50%',
              transform: 'translate(-50%, -50%)',
              zIndex: 2000,
              width: '400px',
              minHeight: '220px',
              display: 'flex',
              flexDirection: 'column',
              alignItems: 'center',
              justifyContent: 'center',
              pointerEvents: 'auto',
              background: 'transparent',
              animation: 'popupElastic 0.7s cubic-bezier(.5,-0.5,.5,1.5)',
            }}
          >
            <h2 style={{ fontSize: '2rem', fontWeight: 'bold', marginBottom: '16px', color: '#fff' }}>Credits</h2>
            <div style={{ color: '#fff', fontSize: '1.2rem', textAlign: 'center', marginBottom: '24px' }}>
              <p>Frontend/UI/UX by: ArcticWoof</p>
              <p>Backend/API by: Cheeseworks</p>
              <br></br>
              <p>Assets by Geode Team and RobTop Games.</p>
            </div>
            <button
              className="nine-slice-button small"
              style={{ marginTop: '8px' }}
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