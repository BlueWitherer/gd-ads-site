import { useNavigate } from 'react-router-dom';
import './App.css';
import CreditsButton from './Credits';

export default function App() {
  const navigate = useNavigate();

  return (
    <>
      <div id="background-scroll"></div>
      <div id="centered-container">
        <h1 className="text-3xl font-bold padding-4 mt-4 mb-8">
          Advertisement Manager
        </h1>
        <button className="nine-slice-button" onClick={() => navigate('/dashboard')}>
          Login
        </button>
      </div>
      <CreditsButton />
    </>
  );
}
