import './App.css';

function App() {
  const publicUrl = process.env.PUBLIC_URL || '';
  const backgroundStyle = {
    backgroundImage: `linear-gradient(rgba(15,22,61,0.65), rgba(15,22,61,0.65)), url('${publicUrl}/background.png')`
  };

  return (
    <div className="App" style={backgroundStyle}>
    </div>
  );
}

export default App;
