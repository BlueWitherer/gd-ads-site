import React from 'react';
import { useNavigate } from 'react-router-dom';
import '../App.css';

function Dashboard({ user, onLogout }) {
  const navigate = useNavigate();

  const handleLogout = () => {
    onLogout();
    navigate('/');
  };

  return (
    <div className="App" style={{
      backgroundImage: `url('${process.env.PUBLIC_URL}/resources/background.png')`,
      backgroundRepeat: 'repeat-x',
      backgroundSize: 'auto 100%'
    }}>
      <div className="center-square" style={{
        borderImageSource: `url('${process.env.PUBLIC_URL}/resources/squareBg.png')`
      }}>
        <div className="dashboard-content">
          <p className="dashboard-welcome">Welcome to the dashboard! yay</p>
        </div>

        {user && (
          <div className="discord-user-container">
            <div className="discord-user-info">
              Logged in as: {user.username}
              {user.discriminator && user.discriminator !== '0' ? `#${user.discriminator}` : ''}
            </div>
            <button
              className="center-button logout-button"
              style={{
                borderImageSource: `url('${process.env.PUBLIC_URL}/resources/button05.png')`
              }}
              onClick={handleLogout}
              aria-label="Logout"
            >
              Logout
            </button>
          </div>
        )}
      </div>
    </div>
  );
}

export default Dashboard;
