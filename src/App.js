import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route, useNavigate, useLocation } from 'react-router-dom';
import './App.css';
import Home from './components/Home';
import Dashboard from './components/Dashboard';

function DashboardWrapper() {
  const [user, setUser] = useState(null);
  const navigate = useNavigate();
  const location = useLocation();

  useEffect(() => {
    const params = new URLSearchParams(location.search);
    const userParam = params.get('user');

    if (userParam) {
      try {
        const userData = JSON.parse(decodeURIComponent(userParam));
        setUser(userData);
        // Store in sessionStorage for persistence
        sessionStorage.setItem('discordUser', JSON.stringify(userData));
        // Clean up URL
        navigate('/dashboard', { replace: true });
      } catch (error) {
        console.error('Failed to parse user data:', error);
      }
    } else {
      // Try to load from sessionStorage
      const storedUser = sessionStorage.getItem('discordUser');
      if (storedUser) {
        setUser(JSON.parse(storedUser));
      } else {
        // No user data, redirect to home
        navigate('/');
      }
    }
  }, [location, navigate]);

  const handleLogout = () => {
    setUser(null);
    sessionStorage.removeItem('discordUser');
  };

  return <Dashboard user={user} onLogout={handleLogout} />;
}

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/dashboard" element={<DashboardWrapper />} />
      </Routes>
    </Router>
  );
}

export default App;
