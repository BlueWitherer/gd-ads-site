import { useNavigate } from 'react-router-dom'
import './App.css'

function App() {
  console.log('App rendered');
  const navigate = useNavigate();

  const handleLogin = () => {
    // test login
    navigate('/dashboard');
  };

  return (
    <>
      <div id="background-scroll"></div>
      <div id="centered-container">
        <h1 className="text-3xl font-bold padding-4 mt-4 mb-8">
          Advertisement Manager
        </h1>
        <button className="nine-slice-button" onClick={handleLogin}>
          Login
        </button>
      </div>
    </>
  )
}

export default App
