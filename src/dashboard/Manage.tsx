import "../App.css";
import square02 from '../assets/square02.png';
import { useEffect, useState } from 'react';

type Ad = {
  id: number;
  type: string;
  levelId: string;
  image: string;
  expiration: string;
}

function Manage() {
  const [adverts, setAdverts] = useState<Ad[] | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function load() {
      try {
        const res = await fetch('/ads/get', { credentials: 'include' });
        if (!res.ok) {
          setError(`Failed to fetch ads: ${res.status}`);
          return;
        }

        const data = await res.json();
        // Expecting array of { id, type, levelId, image, expiration }
        setAdverts(data.map((a: any) => ({
          id: a.id,
          type: a.type,
          levelId: a.levelId,
          image: a.image,
          expiration: a.expiration || '',
        })));
      } catch (err: any) {
        setError(err.message || String(err));
      }
    }

    load();
  }, []);

  return (
    <>
      <h1 className="text-2xl font-bold mb-6">Manage Advertisements</h1>
      <p className="text-lg">Manage and configure your active advertisements.</p>

      {error && <div className="text-red-400">{error}</div>}

      <div style={{ display: 'flex', flexDirection: 'column', gap: '2em', marginTop: '1em' }}>
        {adverts === null ? (
          <div>Loading advertisements...</div>
        ) : adverts.length === 0 ? (
          <div>No advertisements found.</div>
        ) : adverts.map(advert => (
          <div
            key={advert.id}
            style={{
              display: 'flex',
              alignItems: 'center',
              gap: '1.5em',
              color: '#fff',
              borderStyle: 'solid',
              borderWidth: '12px',
              borderImage: `url(${square02}) 24 fill stretch`,
              background: 'transparent',
              borderRadius: '0px',
              padding: '1em',
              boxShadow: '0 2px 8px rgba(0,0,0,0.15)',
              maxWidth: '800px',
            }}
          >
            <img
              src={advert.image}
              alt="Advertisement"
              style={{ width: '160px', height: '160px', overflow: 'hidden', borderRadius: '10px', objectFit: 'contain', background: '#333333' }}
            />
            <div style={{ flex: 1 }}>
              <div><strong>ID:</strong> {advert.id}</div>
              <div><strong>Type:</strong> {advert.type}</div>
              <div><strong>Level ID:</strong> {advert.levelId}</div>
              <div><strong>Expiration:</strong> {advert.expiration}</div>
            </div>
            <button
              style={{
                background: '#e74c3c',
                color: '#fff',
                border: 'none',
                borderRadius: '6px',
                padding: '0.5em 1em',
                cursor: 'pointer',
                fontWeight: 'bold',
                fontSize: '1em',
                transition: 'background 0.2s',
              }}
              onClick={() => {
                fetch(`/ads/delete?id=${advert.id}`, { method: 'DELETE', credentials: 'include' }).then(() => {
                  adverts.splice(adverts.indexOf(advert), 1);
                  setAdverts([...adverts]);
                });
              }}
            >
              Delete
            </button>
          </div>
        ))}
      </div>
    </>
  );
}

export default Manage;
