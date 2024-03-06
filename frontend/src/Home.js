import React, { useState, useEffect } from "react";
import { Link, useNavigate } from 'react-router-dom';
import { getNumberOfPlayers, getPicksDeadline, getPicks } from './backend'; // Correct import statement

const Home = () => {
    const [error, setError] = useState('');
    const [venmoName, setVenmoName] = useState("");
    const [pin, setPin] = useState("");
    const navigate = useNavigate();
    const [picksDeadline, setPicksDeadline] = useState('<insert_date>')
    const [prizePool, setPrizePool] = useState('<insert_data>')

    useEffect(() => {
        const fetchPicksDeadline = async () => {
          try {
            const resp = await getPicksDeadline();
            console.log(resp)
            if (resp.error) {
                console.error('Error fetching picks deadline:', resp.error);
              return
            }
            setPicksDeadline(resp.picksDeadline)
          } catch (error) {
            console.error('Error fetching picks deadline:', error);
          }
        };
        const fetchNumberOfPlayers = async () => {
          try {
            const resp = await getNumberOfPlayers();
            if (resp.error) {
                console.error('Error fetching number of players:', resp.error);
              return
            }
            setPrizePool(resp.numberOfPlayers * 20)
          } catch (error) {
            console.error('Error fetching number of players:', error);
          }
        };    
        fetchPicksDeadline();
        fetchNumberOfPlayers();
      }, []); // Empty dependency array ensures the effect runs only once after the initial render
        

    const handleEditPicks = async () => {
      if (!venmoName || !pin) {
        setError('Please fill in Venmo Name and/or PIN');
        return
      }
      // Make GET request with payload
        const resp = await getPicks(venmoName, pin)

        // if bad response; do nothing
        console.log(resp)
        if (resp.error) {
          setError(resp.error)
          return
        }

        // else pass response to PicksForm
        const firstName = resp.firstName
        const lastName = resp.lastName
        const teamName = resp.teamName
        const selectedTeams = resp.picks

        // Use the navigate function to go to the /picks route
        navigate('/picks', { state: { venmoName, pin, firstName, lastName, teamName, selectedTeams } });
    };


    return (
        <div>
            <h1>Welcome to 'The March Madness 16'!</h1>
            <h2>Brought to you by Magnus Fyhr</h2>
            <h2>Go to 'Rules' to understand how to play</h2>
            <h2>Picks lock @ {picksDeadline}</h2>
            <h2>Current prize pool : ${prizePool}.00</h2>
            <Link to="/leaderboard">
                <button>Leaderboard</button>
            </Link>
            <br />
            <Link to="/picks">
                <button>Make Picks</button>
            </Link>
            <br />
            <br />
            <div className="input-container">
              <input
                  type="text"
                  placeholder="Venmo Name"
                  value={venmoName}
                  onChange={(e) => setVenmoName(e.target.value)}
              />
            </div>
            <div className="input-container">
              <input
                  type="text"
                  placeholder="PIN"
                  value={pin}
                  onChange={(e) => setPin(e.target.value)}
              />
            </div>
            <button onClick={handleEditPicks}>Edit Picks</button>
            {error && <p style={{ color: 'red' }}>{error}</p>}
        </div>
    );
};

export default Home;
