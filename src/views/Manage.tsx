import "../App.css";
import square02 from '../assets/square02.png';

function Manage() {
  // placeholder manage advertisements data
  const adverts = [
    {
      id: 'ad-banner', // likely going to be just a number
      type: 'Banner', // type of banner, square, vertical
      levelId: '123456', // level id associated with the ad
      expiration: '3 days left', // time left until expiration (7-6 days is green, 5-3 yellow, 2-0 red)
      image: 'https://via.placeholder.com/728x90?text=Banner+Ad', // endpoint of the ad image
    },
    {
      id: 'ad-square',
      type: 'Square',
      levelId: '654321',
      expiration: '5 days left',
      image: 'https://via.placeholder.com/180x180?text=Square+Ad',
    },
    {
      id: 'ad-vertical',
      type: 'Vertical',
      levelId: '112233',
      expiration: '1 day left',
      image: 'https://via.placeholder.com/90x728?text=Vertical+Ad',
    },
  ];

  return (
    <>
      <h1 className="text-2xl font-bold mb-6">Manage Advertisements</h1>
      <p className="text-lg">
        Manage and configure your advertisement.
      </p>
      <div style={{ display: 'flex', flexDirection: 'column', gap: '2em', marginTop: '1em' }}>
        {adverts.map(advert => (
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
              style={{ width: '160px', height: 'auto', borderRadius: '8px', objectFit: 'cover', background: '#444' }}
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
              onClick={() => {}}
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
