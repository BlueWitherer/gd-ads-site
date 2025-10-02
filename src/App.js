import './App.css';

function App() {
  const publicUrl = process.env.PUBLIC_URL || '';
  const appStyle = {
    backgroundImage: `url('${publicUrl}/resources/background.png')`,
    backgroundRepeat: 'repeat-x',
    backgroundSize: 'auto 100%'
  };

  const squareStyle = {
    borderImageSource: `url('${publicUrl}/resources/squareBg.png')`
  };

  const buttonStyle = {
    borderImageSource: `url('${publicUrl}/resources/button05.png')`
  };

  return (
    // shows the websites output i think lol
    <div className="App" style={appStyle}>
      <div className="center-square" style={squareStyle}>
        <button
          className="center-button"
          style={buttonStyle}
          onClick={() => console.log('Center button clicked')}
          aria-label="Login to Discord"
        >
          Login via Discord
        </button>
      </div>
    </div>
  );
}

export default App;
