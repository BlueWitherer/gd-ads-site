import React from 'react';
import '../App.css';

function Home() {
  const handleDiscordLogin = () => {
    window.location.href = '/auth/discord';
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
        <button
          className="center-button"
          style={{
            borderImageSource: `url('${process.env.PUBLIC_URL}/resources/button05.png')`
          }}
          onClick={handleDiscordLogin}
          aria-label="Login via Discord"
        >
          Login via Discord
        </button>
      </div>
    </div>
  );
}

export default Home;
