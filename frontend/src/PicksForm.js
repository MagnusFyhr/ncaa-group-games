import React, { useState, useEffect } from 'react';
import { useLocation } from 'react-router-dom';
import { getTeams, putPicks, getPicks } from './backend'; // Correct import statement


const PicksForm = () => {
  const location = useLocation();
  const { state } = location;

  const [error, setError] = useState('');
  const [submitted, setSubmitted] = useState('');

  const [venmoName, setVenmoName] = useState(state?.venmoName || '');
  const [PIN, setPIN] = useState(state?.pin || '');
  const [firstName, setFirstName] = useState(state?.firstName || '');
  const [lastName, setLastName] = useState(state?.lastName || '');
  const [teamName, setTeamName] = useState(state?.teamName || '');
  const [selectedTeams, setSelectedTeams] = useState(state?.selectedTeams || Array.from({ length: 16 }, (_, i) => ''));
  const [allTeams, setAllTeams] = useState([]);

  const seeds = Array.from({ length: 16 }, (_, i) => i + 1);

  useEffect(() => {
    const fetchTeams = async () => {
      try {
        const resp = await getTeams();
        console.log(resp);
        if (resp.error) {
          setError(resp.error)
          return
        }
        setAllTeams(resp.teams)
      } catch (error) {
        console.error('Error fetching teams:', error);
      }
    };

    fetchTeams();
  }, []); // Empty dependency array ensures the effect runs only once after the initial render

  const validateForm = () => {
    if (!venmoName || !PIN || !firstName || !lastName || !teamName || selectedTeams.some(team => !team)) {
      setError('Please fill in all fields');
      setSubmitted('')
      return false;
    }
    setError('');
    return true;
  };

  const handleSubmit = async () => {
    console.log("handleSubmit")
    if (!validateForm()) return;
    console.log("handleSubmit")

    const picks = seeds.map((seed, index) => ({ seed, team: selectedTeams[index] }));

    // Make PUT request with createPicksPayload
    const resp = await putPicks(venmoName, PIN, firstName, lastName, teamName, picks.map(pick => pick.team))
    console.log(resp)
    if(resp.error) {
      setError('Failed to submit picks! ' + resp.error)
      setSubmitted('')
    } else {
      setSubmitted('Successfully submitted picks!')
      setError('')
    }
  };

  const handleLoad = async () => {
    if (!venmoName || !PIN) {
      setError('Please fill in Venmo Name and/or PIN');
      return
    }

    console.log("hello")
    // Make GET request with payload
    const resp = await getPicks(venmoName, PIN)
    console.log("hello2")
    console.log(resp)

    if (resp.error) {
      setError(resp.error)
      return
    }
    
    // set fields
    setVenmoName(resp.venmoName)
    setFirstName(resp.firstName)
    setLastName(resp.lastName)
    setTeamName(resp.teamName)
    setSelectedTeams(resp.picks)
  };

  return (
    <div>
      <div className="input-container">
        <input type="text" placeholder="Venmo Name" value={venmoName} onChange={(e) => setVenmoName(e.target.value)} />
      </div>
      <div className="input-container">
        <input type="text" placeholder="PIN" value={PIN} onChange={(e) => setPIN(e.target.value)} />
      </div>
      <div className="input-container">
        <input type="text" placeholder="First Name" value={firstName} onChange={(e) => setFirstName(e.target.value)} />
      </div>
      <div className="input-container">
        <input type="text" placeholder="Last Name" value={lastName} onChange={(e) => setLastName(e.target.value)} />
      </div>
      <div className="input-container">
        <input type="text" placeholder="Team Name" value={teamName} onChange={(e) => setTeamName(e.target.value)} />
      </div>
      {seeds.map((seed) => (
        <div key={seed} className="input-container">
          <label htmlFor={`seed${seed}`}>{`${seed} Seed:`}</label>
          <select
            id={`seed${seed}`}
            value={selectedTeams[seed - 1]}
            onChange={(e) => {
              const updatedSelectedTeams = [...selectedTeams];
              updatedSelectedTeams[seed - 1] = e.target.value;
              setSelectedTeams(updatedSelectedTeams);
            }}
          >
            <option value="">Select a team</option>
            {allTeams
              .filter((team) => team.seed === seed)
              .map((team) => (
                <option key={team.id} value={team.teamName}>
                  {team.teamName}
                </option>
              ))}
          </select>
        </div>
      ))}
      {submitted && <p style={{ color: 'green' }}>{submitted}</p>}
      {error && <p style={{ color: 'red' }}>{error}</p>}
      <button onClick={handleSubmit}>Submit</button>
      <br />
      <button onClick={handleLoad}>Load Picks</button>
    </div>
  );
};

export default PicksForm;
