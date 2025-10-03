import { useNavigate } from 'react-router-dom'
import './App.css'

export default function NotFound() {
    const navigate = useNavigate();

    const onBack = () => {
        navigate('/');
    };

    return (
        <>
            <div id="centered-container">
                <h1 className='text-6xl mt-8 mb-4'>Oops!</h1>
                <h6 className='text-1xl mb-4'>404 - Page Not Found</h6>

                <button className="nine-slice-button" onClick={onBack}>
                    Go Back
                </button>
            </div>
        </>
    );
}
