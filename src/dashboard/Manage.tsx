import "../App.css";
import square02 from '../assets/square02.png';
import { useEffect, useState } from 'react';

type Ad = {
  id: number;
  type: string;
  level_id: string;
  image: string;
  expiration: number;
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
        // Expecting array of { id, type, level_id, image, expiration }
        setAdverts(data.map((a: any) => ({
          id: a.ad_id,
          type: a.type,
          level_id: a.level_id,
          image: a.image_url,
          expiration: a.expiry,
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
              <div><strong>Level ID:</strong> {advert.level_id}</div>
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
